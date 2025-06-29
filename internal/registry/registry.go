package registry

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/csams/connections-api/internal/defaulters/api"
)

type DefaulterHookRegistry struct {
	Scheme   *runtime.Scheme
	registry map[schema.GroupVersionKind]*admission.Webhook
}

// Add wraps the defaulter in an admission.Webhook and adds it to the registry
func (d *DefaulterHookRegistry) Add(defaulter api.Defaulter) {
	obj := defaulter.Object()
	if gvk, err := apiutil.GVKForObject(obj, d.Scheme); err == nil {
		if _, found := d.registry[gvk]; found {
			panic(fmt.Errorf("Duplicate registration: %v", gvk))
		}
		d.registry[gvk] = admission.WithCustomDefaulter(d.Scheme, obj, defaulter)
	} else {
		panic(err)
	}
}

func (d *DefaulterHookRegistry) Lookup(gvk schema.GroupVersionKind) (*admission.Webhook, bool) {
	hook, found := d.registry[gvk]
	return hook, found
}

func New(scheme *runtime.Scheme) *DefaulterHookRegistry {
	return &DefaulterHookRegistry{
		Scheme:   scheme,
		registry: map[schema.GroupVersionKind]*admission.Webhook{},
	}
}
