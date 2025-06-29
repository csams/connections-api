// We look at deployments for testing. We likely won't check them as a workload type.
package deployment

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"

	"sigs.k8s.io/controller-runtime/pkg/log"
	//"github.com/csams/connections-api/internal/defaulters/api"
)

type DeploymentBinder struct{}

func (d *DeploymentBinder) Bind(ctx context.Context, obj *appsv1.Deployment) error {
	logger := log.FromContext(ctx)

	gvk := obj.GetObjectKind().GroupVersionKind()
	logger.Info(fmt.Sprintf("Handling: %v\n", gvk))

	return nil
}

func New() *DeploymentBinder {
	return &DeploymentBinder{}
}
