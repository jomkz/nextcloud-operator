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
	"fmt"

	serverv1 "github.com/jmckind/nextcloud-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// newConfigMap returns a new ConfigMap instance for the given NextCloud.
func newConfigMap(meta metav1.ObjectMeta) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.Name,
			Namespace: meta.Namespace,
			Labels:    defaultLabels(meta.Name),
		},
	}
}

// newConfigMapWithName creates a new ConfigMap with the given name for the given NextCloud.
func newConfigMapWithName(name string, meta metav1.ObjectMeta) *corev1.ConfigMap {
	cm := newConfigMap(meta)
	cm.ObjectMeta.Name = name
	return cm
}

// newConfigMapWithName creates a new ConfigMap with the given suffix appended to the name.
// The name for the CongifMap is based on the name of the given NextCloud.
func newConfigMapWithSuffix(suffix string, meta metav1.ObjectMeta) *corev1.ConfigMap {
	return newConfigMapWithName(fmt.Sprintf("%s-%s", meta.Name, suffix), meta)
}

// reconcileRedisHAConfigMap will ensure that the main ConfigMap is present for the given NextCloud.
func (r *NextCloudReconciler) reconcileNextCloudConfigMap(cr *serverv1.NextCloud) error {
	cm := newConfigMapWithSuffix("config", cr.ObjectMeta)
	if isObjectFound(r.Client, cr.Namespace, cm.Name, cm) {
		return nil // ConfigMap found with nothing to do, move along...
	}

	cm.Data = map[string]string{
		"ports.conf":  NextCloudDefaultPortsConfig,
		"apache.conf": NextCloudDefaultApacheConfig,
	}

	if err := controllerutil.SetControllerReference(cr, cm, r.Scheme); err != nil {
		return err
	}
	return r.Create(context.TODO(), cm)
}
