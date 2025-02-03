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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CloudProvider string

const (
	ProviderAKS CloudProvider = "AKS"
	ProviderEKS CloudProvider = "EKS"
	ProviderGKE CloudProvider = "GKE"
)

// HubClusterSpec defines the desired state of HubCluster.
type HubClusterSpec struct {
	// Provider is the cloud provider of the remote cluster
	Provider CloudProvider `json:"provider"`

	// Region is the cloud provider region where the cluster is located
	Region string `json:"region"`

	// Credentials references a secret containing cloud provider credentials
	Credentials SecretReference `json:"credentials"`

	// ClusterName is the name of the remote cluster in the cloud provider
	ClusterName string `json:"clusterName"`
}

// HubClusterStatus defines the observed state of HubCluster.
type HubClusterStatus struct {
	// Connected indicates if the hub can communicate with the remote cluster
	Connected bool `json:"connected"`

	// LastConnectionTime is the last time the hub successfully connected to the remote
	LastConnectionTime *metav1.Time `json:"lastConnectionTime,omitempty"`

	// Message provides additional status information
	Message string `json:"message,omitempty"`
}

// SecretReference references a secret containing cloud provider credentials
type SecretReference struct {
	Data map[string]string `json:"data"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HubCluster is the Schema for the hubclusters API.
type HubCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HubClusterSpec   `json:"spec,omitempty"`
	Status HubClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HubClusterList contains a list of HubCluster.
type HubClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HubCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HubCluster{}, &HubClusterList{})
}
