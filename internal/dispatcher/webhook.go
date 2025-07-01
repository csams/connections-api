package dispatcher

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/csams/connections-api/internal/registry"
)

type Dispatcher struct {
	Scheme         *runtime.Scheme
	DefaulterHooks *registry.DefaulterHookRegistry
}

// New creates an admission.Webhook dispatches each request to an admission.Webhook wrapping a type-specific
// admission.CustomDefaulter
func New(scheme *runtime.Scheme, hooks *registry.DefaulterHookRegistry) *admission.Webhook {
	return &admission.Webhook{
		Handler: &Dispatcher{
			Scheme:         scheme,
			DefaulterHooks: hooks,
		},
	}
}

func (dispatcher *Dispatcher) Handle(ctx context.Context, req admission.Request) admission.Response {
	// Will workloads always have a name when they're created?
	// If not, how do we tell which connections should apply to it?

	// what do we have to work with that doesn't require modifying the object to give us a hint?
	// - userinfo (maybe?) - name and maybe groups
	// - namespace of the object, GVK of the object, any of the object's labels

	// if we need a hint on the object:
	// the workload either must have a label with an alternate identity we can use to match against a binding object
	// or it must have annotations that specify the connections that apply
	// the dashboard currently requires annotations.

	logger := log.FromContext(ctx)

	gvk := req.Kind
	if hook, ok := dispatcher.DefaulterHooks.Lookup(gvk); ok {
		logger.Info(fmt.Sprintf("Processing: %v", gvk))
		return hook.Handle(ctx, req)
	}

	logger.Info(fmt.Sprintf("Skipping: %v", gvk))

	return admission.Allowed("")
}
