// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const (
	// MachineFinalizer allows ReconcileMetalMachine to clean up resources associated with MetalMachine before
	// removing it from the apiserver.
	MachineFinalizer = "metalmachine.infrastructure.cluster.x-k8s.io"

	// DefaultReconcilerRequeue is the default value for the reconcile retry.
	DefaultReconcilerRequeue = 10 * time.Second
)

// MetalMachineSpec defines the desired state of MetalMachine
type MetalMachineSpec struct {
	// ProviderID is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderID *string `json:"providerID,omitempty"`

	// ServerSelector specifies matching criteria for labels on Servers.
	// This is used to claim specific Server types for a MetalMachine.
	// +optional
	ServerSelector ServerSelector `json:"serverSelector,omitempty"`
}

// MetalMachineStatus defines the observed state of MetalMachine
type MetalMachineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MetalMachine is the Schema for the metalmachines API
type MetalMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MetalMachineSpec   `json:"spec,omitempty"`
	Status MetalMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MetalMachineList contains a list of MetalMachine
type MetalMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MetalMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MetalMachine{}, &MetalMachineList{})
}
