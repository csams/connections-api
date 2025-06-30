package inferenceservice

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"

	// "Reuse storage-initializer for model loading and authentication from a specified URI (e.g., from Hugging Face).":
	// https://docs.google.com/document/d/11ZQJ2VhTc42S9K4yau2dMs3Q3f4jqWJL_7Sq14C3hzY/edit?tab=t.0

	// https://pkg.go.dev/github.com/kserve/kserve@v0.15.2/pkg/apis/serving/v1beta1
	// for InferenceService
	"github.com/csams/connections-api/internal/defaulter"
	serving "github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	// This will be part of KServe, but I'm not sure this is defined anywhere in an accessible git repo yet
	// LLMInferenceService
	// requires go 1.24 - version chasing may be an issue
	// ModelService "github.com/llm-d/llm-d-model-service/api/v1alpha1"
)

var (
	_ defaulter.Mutator[*serving.InferenceService] = &InferenceServiceBinder{}

	AddToScheme = serving.AddToScheme
)

type InferenceServiceBinder struct{}

func (dep *InferenceServiceBinder) Mutate(ctx context.Context, obj *serving.InferenceService) error {
	logger := log.FromContext(ctx)

	gvk := obj.GetObjectKind().GroupVersionKind()
	logger.Info(fmt.Sprintf("Handling %s\n", gvk.String()))

	return nil
}

func New() *InferenceServiceBinder {
	return &InferenceServiceBinder{}
}
