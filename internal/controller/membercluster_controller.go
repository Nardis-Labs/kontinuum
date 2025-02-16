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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kontinuumv1alpha1 "lab.nardis.io/kontinuum/api/v1alpha1"
	"lab.nardis.io/kontinuum/pkg/providers"
)

// MemberClusterReconciler reconciles a MemberCluster object
type MemberClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=kontinuum.nardis.io,resources=MemberClusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kontinuum.nardis.io,resources=MemberClusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kontinuum.nardis.io,resources=MemberClusters/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *MemberClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Get the RemoteCluster resource
	var remoteCluster kontinuumv1alpha1.MemberCluster
	if err := r.Get(ctx, req.NamespacedName, &remoteCluster); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get the appropriate cloud provider client
	providerClient, err := providers.GetCloudProviderClient(remoteCluster.Spec.Provider)
	if err != nil {
		logger.Error(err, "failed to get cloud provider client")
		return ctrl.Result{}, err
	}

	// Verify cluster connectivity
	connected, err := providerClient.VerifyClusterConnection(ctx, remoteCluster.Spec)
	if err != nil {
		logger.Error(err, "failed to verify cluster connection")
		remoteCluster.Status.Connected = false
		remoteCluster.Status.Message = fmt.Sprintf("Connection failed: %v", err)
	} else {
		remoteCluster.Status.Connected = connected
		if connected {
			remoteCluster.Status.Message = "Successfully connected to remote cluster"
			now := metav1.Now()
			remoteCluster.Status.LastConnectionTime = &now
		}
	}

	// Update the RemoteCluster status
	if err := r.Status().Update(ctx, &remoteCluster); err != nil {
		logger.Error(err, "failed to update RemoteCluster status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MemberClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kontinuumv1alpha1.MemberCluster{}).
		Named("MemberCluster").
		Complete(r)
}
