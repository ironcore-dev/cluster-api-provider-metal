// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// ClusterFinalizer allows MetalClusterReconciler to clean up resources associated with MetalCluster before
	// removing it from the apiserver.
	ClusterFinalizer = "metalcluster.infrastructure.cluster.x-k8s.io"
)

// MetalClusterSpec defines the desired state of MetalCluster
type MetalClusterSpec struct {
	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint,omitempty"`
}

// MetalClusterStatus defines the observed state of MetalCluster
type MetalClusterStatus struct {
	// Ready denotes that the cluster (infrastructure) is ready.
	// +optional
	Ready bool `json:"ready"`

	// Conditions defines current service state of the MetalCluster.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MetalCluster is the Schema for the metalclusters API
type MetalCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MetalClusterSpec   `json:"spec,omitempty"`
	Status MetalClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MetalClusterList contains a list of MetalCluster
type MetalClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MetalCluster `json:"items"`
}

// GetConditions returns the observations of the operational state of the MetalCluster resource.
func (c *MetalCluster) GetConditions() clusterv1.Conditions {
	return c.Status.Conditions
}

// SetConditions sets the underlying service state of the MetalCluster to the predescribed clusterv1.Conditions.
func (c *MetalCluster) SetConditions(conditions clusterv1.Conditions) {
	c.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&MetalCluster{}, &MetalClusterList{})
}
