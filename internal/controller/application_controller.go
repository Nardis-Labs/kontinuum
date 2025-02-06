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

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kontinuumv1alpha1 "lab.nardis.io/kontinuum/api/v1alpha1"
	"lab.nardis.io/kontinuum/pkg/providers"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=kontinuum.nardis.io,resources=applications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kontinuum.nardis.io,resources=applications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kontinuum.nardis.io,resources=applications/finalizers,verbs=update
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var app kontinuumv1alpha1.Application
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Deploy the application to the target cluster
	var cluster kontinuumv1alpha1.MemberCluster
	if err := r.Get(ctx, client.ObjectKey{Name: app.Spec.TargetCluster}, &cluster); err != nil {
		logger.Error(err, "failed to get target cluster", "cluster", app.Spec.TargetCluster)
		return ctrl.Result{}, err
	}

	if err := r.deployToCluster(ctx, &app, &cluster); err != nil {
		logger.Error(err, "failed to deploy application", "cluster", app.Spec.TargetCluster)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kontinuumv1alpha1.Application{}).
		Named("application").
		Complete(r)
}

func (r *ApplicationReconciler) deployToCluster(ctx context.Context, app *kontinuumv1alpha1.Application, cluster *kontinuumv1alpha1.MemberCluster) error {
	// Get the appropriate cloud provider client
	providerClient, err := providers.GetCloudProviderClient(cluster.Spec.Provider)
	if err != nil {
		return err
	}

	// Deploy the application
	switch app.Spec.Source.Type {
	case "helm":
		return providerClient.DeployHelmChart(ctx, app, cluster)
	case "kustomize":
		return providerClient.DeployKustomizeManifests(ctx, app, cluster)
	default:
		return fmt.Errorf("unsupported application source type: %s", app.Spec.Source.Type)
	}
}
