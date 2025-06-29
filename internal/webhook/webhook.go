package webhook

import (
	"context"
	"fmt"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/csams/connections-api/internal/registry"
)

type Dispatcher struct {
	Scheme         *runtime.Scheme
	Decoder        admission.Decoder
	DefaulterHooks *registry.DefaulterHookRegistry
}

// NewDispatcher creates an admission.Webhook that decodes each request to a PartialObjectMetadata and then
// dispatches it to an admission.Webhook wrapping a type-specific admission.CustomDefaulter
func NewDispatcher(scheme *runtime.Scheme, hooks *registry.DefaulterHookRegistry) *admission.Webhook {
	dispatcher := &Dispatcher{
		Scheme:         scheme,
		Decoder:        admission.NewDecoder(scheme),
		DefaulterHooks: hooks,
	}

	return (&admission.Webhook{
		Handler: dispatcher,
	})
}

func (dispatcher *Dispatcher) Handle(ctx context.Context, req admission.Request) admission.Response {
	// the Object may not have a name on create
	// how do we tell which connection should apply to it?

	// what do we have to work with that doesn't require modifying the object to give us a hint?
	// - userinfo (maybe?) - name and maybe groups
	// - namespace of the object, GVK of the object, any of the object's labels

	// if we need a hint on the object:
	// the workload either must have a label with an alternate identity we can use to match against a binding object
	// or it must have annotations that specify the connections that apply

	logger := log.FromContext(ctx)

	// partially decode the object so we can get its GVK
	pm := &metav1.PartialObjectMetadata{}
	if err := dispatcher.Decoder.Decode(req, pm); err != nil {
		logger.Error(err, "Error decoding request")
		return admission.Errored(http.StatusBadRequest, err)
	}

	// get the object's gvk and dispatch to the right webhook
	if gvk, err := apiutil.GVKForObject(pm, dispatcher.Scheme); err == nil {
		if defaulter, ok := dispatcher.DefaulterHooks.Lookup(gvk); ok {
			logger.Info(fmt.Sprintf("Processing: %v", gvk))
			return defaulter.Handle(ctx, req)
		} else {
			logger.Info(fmt.Sprintf("Skipping: %v", gvk))
		}
	} else {
		logger.Error(err, "Error getting gvk")
		return admission.Errored(http.StatusBadRequest, err)
	}

	return admission.Allowed("")
}
