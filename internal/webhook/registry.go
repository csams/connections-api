package webhook

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	appsv1 "k8s.io/api/apps/v1"

	// https://pkg.go.dev/github.com/kserve/kserve@v0.15.2/pkg/apis/serving/v1beta1
	// for InferenceService
	serving "github.com/kserve/kserve/pkg/apis/serving/v1beta1"

	// kubeflow has a PodDefaults CR, but it doesn't look like it's been updated in 3 years:
	// https://github.com/kubeflow/kubeflow/blob/master/components/admission-webhook/README.md
	// for notebooks
	nbv1 "github.com/kubeflow/kubeflow/components/notebook-controller/api/v1"
	nbv1beta1 "github.com/kubeflow/kubeflow/components/notebook-controller/api/v1beta1"

	// TODO: add Notebooks v2

	// requires go 1.24 - version chasing may be an issue
	// modelservice "github.com/llm-d/llm-d-model-service/api/v1alpha1"

	// "Reuse storage-initializer for model loading and authentication from a specified URI (e.g., from Hugging Face).":
	// https://docs.google.com/document/d/11ZQJ2VhTc42S9K4yau2dMs3Q3f4jqWJL_7Sq14C3hzY/edit?tab=t.0

	// This will be part of KServe, but I'm not sure this is defined anywhere in an accessible git repo yet
	// LLMInferenceService

	"github.com/csams/connections-api/internal/webhook/defaulters"
)

type ConnectionBinderRegistry map[schema.GroupVersionKind]*admission.Webhook

var (
	defaulterRegistry = map[runtime.Object]admission.CustomDefaulter{
		&appsv1.Deployment{}:        &defaulters.DeploymentDefaulter{},
		&serving.InferenceService{}: &defaulters.InferenceServiceDefaulter{},
		&nbv1.Notebook{}:            &defaulters.NotebookV1Defaulter{},
		&nbv1beta1.Notebook{}:       &defaulters.NotebookV1Beta1Defaulter{},
	}

	// we'll always be chasing workload types and fighting import issues
	localSchemeBuilder = runtime.SchemeBuilder{
		clientgoscheme.AddToScheme,
		serving.AddToScheme,
		nbv1.AddToScheme,
		nbv1beta1.AddToScheme,
	}

	AddToScheme = localSchemeBuilder.AddToScheme
)

func NewConnectionBinderRegistry(scheme *runtime.Scheme) ConnectionBinderRegistry {
	binderRegistry := ConnectionBinderRegistry{}

	// associate the defaulters with the GVK's of the objects they handle
	for obj, defaulter := range defaulterRegistry {
		if gvk, err := apiutil.GVKForObject(obj, scheme); err == nil {
			binderRegistry[gvk] = admission.WithCustomDefaulter(scheme, obj, defaulter)
		} else {
			panic(err)
		}
	}

	return binderRegistry
}
