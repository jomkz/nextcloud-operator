/*
A Kubernetes operator for managing Nextcloud clusters.
Copyright (C) 2021 NextCloud Operator Developers

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	serverv1 "github.com/jmckind/nextcloud-operator/api/v1"
)

// NextCloudReconciler reconciles a NextCloud object
type NextCloudReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=server.nextcloud.com,resources=nextclouds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=server.nextcloud.com,resources=nextclouds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=server.nextcloud.com,resources=nextclouds/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps;services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *NextCloudReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("nextcloud", req.NamespacedName)
	r.Log.Info("Reconcile starting...")

	nextcloud := &serverv1.NextCloud{}
	if err := r.Get(ctx, req.NamespacedName, nextcloud); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	if nextcloud.ObjectMeta.DeletionTimestamp.IsZero() {
		if err := r.reconcileNextCloudResources(nextcloud); err != nil {
			// Error reconciling NextCloud sub-resources - requeue the request.
			return ctrl.Result{}, err
		}
	}

	r.Log.Info("Reconcile complete")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NextCloudReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&serverv1.NextCloud{}).
		Complete(r)
}

// reconcileNextCloudResources will ensure that all resources are reconciled
// for the given NextCloud instance.
func (r *NextCloudReconciler) reconcileNextCloudResources(cr *serverv1.NextCloud) error {
	if err := r.reconcileNextCloudConfigMap(cr); err != nil {
		return err
	}
	if err := r.reconcileServerDeployment(cr); err != nil {
		return err
	}
	if err := r.reconcileNextCloudService(cr); err != nil {
		return err
	}
	return nil
}
