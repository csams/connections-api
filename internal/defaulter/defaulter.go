/*
defaulter abstracts over CustomDefaulter so we can write simple, type specific mutators
*/
package defaulter

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Defaulter instances can be added to DefaulterHookRegistry
type Defaulter interface {
	admission.CustomDefaulter
	Object() runtime.Object
}

type defaulter[T any] struct {
	Mutator Mutator[*T]
	obj     runtime.Object
}

// Default handles type checking and casting before forwarding the type safe object to the Mutator
func (db *defaulter[T]) Default(ctx context.Context, obj runtime.Object) error {
	if o, ok := any(obj).(*T); ok {
		return db.Mutator.Mutate(ctx, o)
	}
	return fmt.Errorf(fmt.Sprintf("obj is of type %T. Expected %T", obj, new(T)))
}

func (db *defaulter[T]) Object() runtime.Object {
	return db.obj
}

type Mutator[T any] interface {
	Mutate(context.Context, T) error
}

func New[T any](m Mutator[*T]) Defaulter {
	t := new(T)
	if obj, ok := any(t).(runtime.Object); ok {
		return &defaulter[T]{
			Mutator: m,
			obj:     obj,
		}
	} else {
		panic(fmt.Errorf("type %T does not implement runtime.Object", t))
	}
}
