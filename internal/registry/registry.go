package registry

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/csams/connections-api/internal/defaulter"
)

type DefaulterHookRegistry struct {
	Scheme   *runtime.Scheme
	registry map[metav1.GroupVersionKind]*admission.Webhook
}

// Add wraps the defaulter in an admission.Webhook and adds it to the registry
// This associateѕ one admission.CustomDefaulter with each GVK, but we could use an
// admission.MultiMutatingHandler for each GVK and then dispatch to multiple CustomDefaulters
func (d *DefaulterHookRegistry) Add(defaulter defaulter.Defaulter) {
	obj := defaulter.Object()
	if gvk, err := apiutil.GVKForObject(obj, d.Scheme); err == nil {
		// metav1.GroupVersionKind is the type available on admission.Request
		g := metav1.GroupVersionKind(gvk)
		if _, found := d.registry[g]; found {
			panic(fmt.Errorf("Duplicate registration: %v", gvk))
		}
		d.registry[g] = admission.WithCustomDefaulter(d.Scheme, obj, defaulter)
	} else {
		panic(err)
	}
}

func (d *DefaulterHookRegistry) Lookup(gvk metav1.GroupVersionKind) (*admission.Webhook, bool) {
	hook, found := d.registry[gvk]
	return hook, found
}

func New(scheme *runtime.Scheme) *DefaulterHookRegistry {
	return &DefaulterHookRegistry{
		Scheme:   scheme,
		registry: map[metav1.GroupVersionKind]*admission.Webhook{},
	}
}
