package deployment

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"

	"github.com/csams/connections-api/internal/defaulter"

	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ defaulter.Mutator[*appsv1.Deployment] = &DeploymentBinder{}

// DeploymentBinder exists so we have a standard object for testing. It's too generic of a workload type to monitor.
type DeploymentBinder struct{}

func (d *DeploymentBinder) Mutate(ctx context.Context, obj *appsv1.Deployment) error {
	logger := log.FromContext(ctx)

	if req, err := admission.RequestFromContext(ctx); err == nil {
		userName := req.UserInfo.Username
		logger.Info(fmt.Sprintf("Handling request of type %T from %s", obj, userName))
	}

	return nil
}

func New() *DeploymentBinder {
	return &DeploymentBinder{}
}
