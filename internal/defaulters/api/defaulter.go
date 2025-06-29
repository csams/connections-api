package api

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Defaulter instances can be added to ConnectionBinderRegistry
type Defaulter interface {
	admission.CustomDefaulter
	Object() runtime.Object
}

type defaulterBase[T any] struct {
	Binder connectionBinder[*T]
	obj    runtime.Object
}

func (db *defaulterBase[T]) Default(ctx context.Context, obj runtime.Object) error {

	if o, ok := any(obj).(*T); ok {
		return db.Binder.Bind(ctx, o)
	}

	return fmt.Errorf(fmt.Sprintf("obj is of type %T. Expected %T", obj, new(T)))
}

func (db *defaulterBase[T]) Object() runtime.Object {
	return db.obj
}

type connectionBinder[T any] interface {
	Bind(context.Context, T) error
}

func New[T any](b connectionBinder[*T]) Defaulter {
	t := new(T)
	if obj, ok := any(t).(runtime.Object); ok {
		return &defaulterBase[T]{
			Binder: b,
			obj:    obj,
		}
	} else {
		panic(fmt.Errorf("type %T does not implement runtime.Object", t))
	}
}
