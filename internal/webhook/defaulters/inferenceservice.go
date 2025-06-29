package defaulters

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	serving "github.com/kserve/kserve/pkg/apis/serving/v1beta1"
)

type InferenceServiceDefaulter struct{}

func (dep *InferenceServiceDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	logger := log.FromContext(ctx)

	if _, ok := obj.(*serving.InferenceService); ok {
		gvk := obj.GetObjectKind().GroupVersionKind()
		logger.Info(fmt.Sprintf("Handling %s\n", gvk.String()))
	} else {
		logger.Info(fmt.Sprintf("Error Handling %T\n", obj))
	}

	return nil
}
