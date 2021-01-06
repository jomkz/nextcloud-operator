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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NextCloudSpec defines the desired state of NextCloud
type NextCloudSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of NextCloud. Edit NextCloud_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// NextCloudStatus defines the observed state of NextCloud
type NextCloudStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// NextCloud is the Schema for the nextclouds API
type NextCloud struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NextCloudSpec   `json:"spec,omitempty"`
	Status NextCloudStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NextCloudList contains a list of NextCloud
type NextCloudList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NextCloud `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NextCloud{}, &NextCloudList{})
}
