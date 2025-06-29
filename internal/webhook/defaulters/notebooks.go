package defaulters

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	nbv1 "github.com/kubeflow/kubeflow/components/notebook-controller/api/v1"
	nbv1beta1 "github.com/kubeflow/kubeflow/components/notebook-controller/api/v1beta1"
)

type NotebookV1Defaulter struct{}

func (dep *NotebookV1Defaulter) Default(ctx context.Context, obj runtime.Object) error {
	logger := log.FromContext(ctx)

	if _, ok := obj.(*nbv1.Notebook); ok {
		gvk := obj.GetObjectKind().GroupVersionKind()
		logger.Info(fmt.Sprintf("Handling %s\n", gvk.String()))
	} else {
		logger.Info(fmt.Sprintf("Error Handling %T\n", obj))
	}

	return nil
}

type NotebookV1Beta1Defaulter struct{}

func (dep *NotebookV1Beta1Defaulter) Default(ctx context.Context, obj runtime.Object) error {
	logger := log.FromContext(ctx)

	if _, ok := obj.(*nbv1beta1.Notebook); ok {
		gvk := obj.GetObjectKind().GroupVersionKind()
		logger.Info(fmt.Sprintf("Handling %s\n", gvk.String()))
	} else {
		logger.Info(fmt.Sprintf("Error Handling %T\n", obj))
	}

	return nil
}
