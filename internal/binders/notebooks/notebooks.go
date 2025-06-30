package notebooks

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	// kubeflow has a PodDefaults CR, but it doesn't look like it's been updated in 3 years:
	// https://github.com/kubeflow/kubeflow/blob/master/components/admission-webhook/README.md
	// for notebooks
	"github.com/csams/connections-api/internal/defaulter"
	nbv1 "github.com/kubeflow/kubeflow/components/notebook-controller/api/v1"
	nbv1beta1 "github.com/kubeflow/kubeflow/components/notebook-controller/api/v1beta1"
	// TODO: add Notebooks v2
)

var (
	_ defaulter.Mutator[*nbv1.Notebook]      = &NotebookV1Binder{}
	_ defaulter.Mutator[*nbv1beta1.Notebook] = &NotebookV1Beta1Binder{}

	localSchemeBuilder = runtime.SchemeBuilder{
		nbv1.AddToScheme,
		nbv1beta1.AddToScheme,
	}

	AddToScheme = localSchemeBuilder.AddToScheme
)

type NotebookV1Binder struct{}

func (dep *NotebookV1Binder) Mutate(ctx context.Context, obj *nbv1.Notebook) error {
	logger := log.FromContext(ctx)

	gvk := obj.GetObjectKind().GroupVersionKind()
	logger.Info(fmt.Sprintf("Handling %s\n", gvk.String()))

	return nil
}

func NewV1Binder() *NotebookV1Binder {
	return &NotebookV1Binder{}
}

type NotebookV1Beta1Binder struct{}

func (dep *NotebookV1Beta1Binder) Mutate(ctx context.Context, obj *nbv1beta1.Notebook) error {
	logger := log.FromContext(ctx)

	gvk := obj.GetObjectKind().GroupVersionKind()
	logger.Info(fmt.Sprintf("Handling %s\n", gvk.String()))

	return nil
}

func NewV1Beta1Binder() *NotebookV1Beta1Binder {
	return &NotebookV1Beta1Binder{}
}
