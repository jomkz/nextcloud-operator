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

	serverv1 "github.com/jmckind/nextcloud-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// newService returns a new Service instance for the given ObjectMeta.
func newService(meta metav1.ObjectMeta) *corev1.Service {
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.Name,
			Namespace: meta.Namespace,
			Labels:    defaultLabels(meta.Name),
		},
	}

	svc.Spec.Selector = defaultLabels(meta.Name)

	return &svc
}

// reconcileNextCloudService will ensure that the Service is present for the given NextCloud.
func (r *NextCloudReconciler) reconcileNextCloudService(cr *serverv1.NextCloud) error {
	svc := newService(cr.ObjectMeta)
	if isObjectFound(r.Client, cr.Namespace, svc.Name, svc) {
		return nil // Service found with nothing to do, move along...
	}

	svc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "http",
			Port:       NextCloudDefaultHTTPPort,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromInt(NextCloudDefaultHTTPPort),
		},
	}

	if err := controllerutil.SetControllerReference(cr, svc, r.Scheme); err != nil {
		return err
	}
	return r.Create(context.TODO(), svc)
}
