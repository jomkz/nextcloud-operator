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
	"strings"

	serverv1 "github.com/jmckind/nextcloud-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// CombineImageTag will return the combined image and tag in the proper format for tags and digests.
func combineImageTag(img string, tag string) string {
	if strings.Contains(tag, ":") {
		return fmt.Sprintf("%s@%s", img, tag) // Digest
	} else if len(tag) > 0 {
		return fmt.Sprintf("%s:%s", img, tag) // Tag
	}
	return img // No tag, use default
}

// getNextCloudContainerImage returns the container image to use for NextCloud server.
func getNextCloudContainerImage() string {
	img := NextCloudDefaultServerImage
	tag := NextCloudDefaultServerVersion
	return combineImageTag(img, tag)
}

// newDeployment returns a new Deployment instance for the given ObjectMeta.
func newDeployment(meta metav1.ObjectMeta) *appsv1.Deployment {
	deploy := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.Name,
			Namespace: meta.Namespace,
			Labels:    defaultLabels(meta.Name),
		},
	}

	deploy.Spec = appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: defaultLabels(meta.Name),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: defaultLabels(meta.Name),
			},
		},
	}

	return &deploy
}

// reconcileServerDeployment will ensure the Deployment resource is present for
// the NextCloud Server component.
func (r *NextCloudReconciler) reconcileServerDeployment(cr *serverv1.NextCloud) error {
	deploy := newDeployment(cr.ObjectMeta)
	if isObjectFound(r.Client, cr.Namespace, deploy.Name, deploy) {
		return nil // Deployment found with nothing to do, move along...
	}

	deploy.Spec.Template.Spec.Containers = []corev1.Container{{
		Image: getNextCloudContainerImage(),
		Name:  "nextcloud",
		Ports: []corev1.ContainerPort{
			{
				ContainerPort: 8080,
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "nextcloud-config",
				MountPath: "/etc/apache2/ports.conf",
				SubPath:   "ports.conf",
			}, {
				Name:      "nextcloud-config",
				MountPath: "/etc/apache/sites-available/000-default.conf",
				SubPath:   "apache.conf",
			},
		},
	}}

	deploy.Spec.Template.Spec.Volumes = []corev1.Volume{
		{
			Name: "nextcloud-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: nameWithSuffix(cr.ObjectMeta, "config"),
					},
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(cr, deploy, r.Scheme); err != nil {
		return err
	}
	return r.Create(context.TODO(), deploy)
}
