// We look at deployments for testing. We likely won't check them as a workload type.
package defaulters

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DeploymentDefaulter struct{}

func (d *DeploymentDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	logger := log.FromContext(ctx)

	if dep, ok := obj.(*appsv1.Deployment); ok {
		gvk := dep.GetObjectKind().GroupVersionKind()
		logger.Info(fmt.Sprintf("Handling: %s\n", gvk.String()))
	} else {
		logger.Info(fmt.Sprintf("Error Handling %T\n", obj))
	}

	return nil
}
