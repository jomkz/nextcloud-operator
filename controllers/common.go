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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// NextCloudAppName is the key for the "app name" label
	NextCloudAppName = "nextcloud"

	// NextCloudDefaultApacheConfig is the default apache vhosts configuration
	NextCloudDefaultApacheConfig = `
<VirtualHost *:8080>
    DocumentRoot /usr/src/nextcloud
    ErrorLog ${APACHE_LOG_DIR}/error.log
    CustomLog ${APACHE_LOG_DIR}/access.log combined
</VirtualHost>`

	// NextCloudDefaultHTTPPort is the default HTTP listen port for NextCloud server
	NextCloudDefaultHTTPPort = 8080

	// NextCloudDefaultPortsConfig is the default apache ports configuration
	NextCloudDefaultPortsConfig = "Listen 8080"

	// NextCloudDefaultServerImage is the default container image to use for NextCloud server
	NextCloudDefaultServerImage = "nextcloud"

	// NextCloudDefaultServerVersion is the default version tag to use for the NextCloud server container image
	NextCloudDefaultServerVersion = "20.0.4-apache"

	// NextCloudDefaultServerResourceLimitCPU is the default CPU limits for NextCloud server
	NextCloudDefaultServerResourceLimitCPU = "1000m"

	// NextCloudDefaultServerResourceLimitMemory is the default memory limits for NextCloud server
	NextCloudDefaultServerResourceLimitMemory = "1Gi"

	// NextCloudDefaultServerResourceRequestCPU is the default CPU requests for NextCloud server
	NextCloudDefaultServerResourceRequestCPU = "250m"

	// NextCloudDefaultServerResourceRequestMemory is the default memory requests for NextCloud server
	NextCloudDefaultServerResourceRequestMemory = "512Mi"

	// NextCloudKeyAppName is the key for the "app name" label
	NextCloudKeyAppName = "app.kubernetes.io/name"

	// NextCloudKeyPartOf is the key for the "part-of" label
	NextCloudKeyPartOf = "app.kubernetes.io/part-of"
)

// defaultLabels returns the default set of labels for controllers.
func defaultLabels(name string) map[string]string {
	return map[string]string{
		NextCloudKeyAppName: name,
		NextCloudKeyPartOf:  NextCloudAppName,
	}
}

// fetchObject will retrieve the object with the given namespace and name using the Kubernetes API.
// The result will be stored in the given object.
func fetchObject(client client.Client, namespace string, name string, obj client.Object) error {
	return client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, obj)
}

// IsObjectFound will perform a basic check that the given object exists via the Kubernetes API.
// If an error occurs as part of the check, the function will return false.
func isObjectFound(client client.Client, namespace string, name string, obj client.Object) bool {
	return !apierrors.IsNotFound(fetchObject(client, namespace, name, obj))
}

// nameWithSuffix will return a string using the Name from the given ObjectMeta with the provded suffix appended.
// Example: If ObjectMeta.Name is "test" and suffix is "object", the value of "test-object" will be returned.
func nameWithSuffix(meta metav1.ObjectMeta, suffix string) string {
	return fmt.Sprintf("%s-%s", meta.Name, suffix)
}
