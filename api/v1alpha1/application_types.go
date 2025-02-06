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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ApplicationSpec defines the desired state of Application.
type ApplicationSpec struct {
	// Source defines where to get the application manifests from
	Source ApplicationSource `json:"source"`

	// TargetClusters defines which clusters to deploy to
	TargetCluster string `json:"targetClusters"`

	// TargetNamespace defines the namespace to deploy the application to
	TargetNamespace string `json:"targetNamespace"`

	// Chart defines the name of the Helm chart if using Helm
	Chart ApplicationChartSource `json:"chart,omitempty"`

	// Values for helm charts (if using helm)
	Values map[string]string `json:"values,omitempty"`
}

// ApplicationSource defines where to get the application manifests from
type ApplicationSource struct {
	// Type is either "helm" or "kustomize"
	Type string `json:"type"`

	// Path to the chart/manifests
	Path string `json:"path"`

	// Repository URL
	RepoURL string `json:"repoUrl"`

	// Version/tag to use
	Version string `json:"version"`
}

// ApplicationChartSource defines the Helm chart details
type ApplicationChartSource struct {
	// Repository URL
	Repository string `json:"repository"`

	// Name of the repository
	RepoName string `json:"repoName"`
}

// ApplicationStatus defines the observed state of Application.
type ApplicationStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Application is the Schema for the applications API.
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ApplicationList contains a list of Application.
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
