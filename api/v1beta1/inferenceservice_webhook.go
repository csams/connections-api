/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var inferenceservicelog = logf.Log.WithName("inferenceservice-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *InferenceService) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-serving-kserve-io-opendatahub-io-v1beta1-inferenceservice,mutating=true,failurePolicy=fail,sideEffects=None,groups=serving.kserve.io.opendatahub.io,resources=inferenceservices,verbs=create;update,versions=v1beta1,name=minferenceservice.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &InferenceService{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *InferenceService) Default() {
	inferenceservicelog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}
