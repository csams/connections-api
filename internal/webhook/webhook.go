package webhook

import (
	"context"
	"fmt"
	"net/http"

	// https://pkg.go.dev/github.com/kserve/kserve@v0.15.2/pkg/apis/serving/v1beta1
	// InferenceService
	serving "github.com/kserve/kserve/pkg/apis/serving/v1beta1"

	// https://deepwiki.com/llm-d/llm-d-deployer/3-modelservice-crd
	// requires go 1.24 and I don't feel like making a fork to try to downgrade it right now
	//modelservice "github.com/llm-d/llm-d-model-service/api/v1alpha1"

	// for notebooks https://github.com/kubeflow/kubeflow/blob/master/components/admission-webhook/README.md
	// kubeflow has a PodDefaults CR, but it doesn't look like it's been updated in 3 years.
	nbv1 "github.com/kubeflow/kubeflow/components/notebook-controller/api/v1"
	nbv1beta1 "github.com/kubeflow/kubeflow/components/notebook-controller/api/v1beta1"

	// https://docs.google.com/document/d/11ZQJ2VhTc42S9K4yau2dMs3Q3f4jqWJL_7Sq14C3hzY/edit?tab=t.0
	// "Reuse storage-initializer for model loading and authentication from a specified URI (e.g., from Hugging Face)."
	// LLMInferenceService

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	//"k8s.io/apimachinery/pkg/runtime/schema"

	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	InferenceServiceGVK = (&serving.InferenceService{}).GroupVersionKind()
	NBV1GVK             = (&nbv1.Notebook{}).GroupVersionKind()
	NBV1Beta1GVK        = (&nbv1beta1.Notebook{}).GroupVersionKind()
	// ModelServiceGVK     = (&modelservice.ModelService{}).GroupVersionKind().String()
)

type DeploymentDefaulter struct{}

func (dep *DeploymentDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	logger := log.FromContext(ctx)

	if _, ok := obj.(*appsv1.Deployment); ok {
		gvk := obj.GetObjectKind().GroupVersionKind()
		logger.Info(fmt.Sprintf("Handling %s\n", gvk.String()))
	} else {
		logger.Info(fmt.Sprintf("Error Handling %T\n", obj))
	}

	return nil
}

func NewConnectionBindingWebhook(scheme *runtime.Scheme) *admission.Webhook {
	binders := map[schema.GroupVersionKind]*admission.Webhook{}

	defaulters := map[runtime.Object]admission.CustomDefaulter{
		&appsv1.Deployment{}: &DeploymentDefaulter{},
	}

	for obj, defaulter := range defaulters {
		if gvk, err := apiutil.GVKForObject(obj, scheme); err == nil {
			binders[gvk] = admission.WithCustomDefaulter(scheme, obj, defaulter)
		} else {
			panic(err)
		}
	}

	decoder := admission.NewDecoder(scheme)
	dispatcher := makeDispatcher(scheme, decoder, binders)
	return (&admission.Webhook{
		Handler: admission.HandlerFunc(dispatcher),
	}).WithRecoverPanic(true)
}

func makeDispatcher(scheme *runtime.Scheme, decoder admission.Decoder, binders map[schema.GroupVersionKind]*admission.Webhook) func(context.Context, admission.Request) admission.Response {
	return func(ctx context.Context, req admission.Request) admission.Response {
		logger := log.FromContext(ctx)
		// the Object may not have a name on create
		// how do we tell which connection should apply to it?

		// what do we have to work with that doesn't require modifying the object to give us a hint?
		// - userinfo (maybe?) - name and maybe groups
		// - namespace of the object, GVK of the object, any of the object's labels

		// if we need a hint,
		// the workload either must have a label with an alternate identity we can use to match against a binding object
		// or it must have annotations that specify the connections that apply

		pm := &metav1.PartialObjectMetadata{}
		if err := decoder.Decode(req, pm); err != nil {
			logger.Error(err, "Error decoding request")
			return admission.Errored(http.StatusBadRequest, err)
		}

		if gvk, err := apiutil.GVKForObject(pm, scheme); err != nil {
			logger.Error(err, "Error getting gvk")
			return admission.Errored(http.StatusBadRequest, err)
		} else {
			for k := range binders {
				logger.Info(fmt.Sprintf("Binders: %v", k))
			}
			logger.Info(fmt.Sprintf("Processing: %v", gvk))
			if binder, ok := binders[gvk]; ok {
				return binder.Handle(ctx, req)
			}
		}

		return admission.Allowed("")
	}
}
