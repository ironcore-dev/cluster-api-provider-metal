// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// MachineFinalizer allows ReconcileIroncoreMetalMachine to clean up resources associated with IroncoreMetalMachine before
	// removing it from the apiserver.
	MachineFinalizer = "ironcoremetalmachine.infrastructure.cluster.x-k8s.io"

	// DefaultReconcilerRequeue is the default value for the reconcile retry.
	DefaultReconcilerRequeue = 5 * time.Second
)

// IroncoreMetalMachineSpec defines the desired state of IroncoreMetalMachine
type IroncoreMetalMachineSpec struct {
	// ProviderID is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderID *string `json:"providerID,omitempty"`

	// Image specifies the boot image to be used for the server.
	Image string `json:"image"`

	// ServerSelector specifies matching criteria for labels on Servers.
	// This is used to claim specific Server types for a IroncoreMetalMachine.
	// +optional
	ServerSelector *metav1.LabelSelector `json:"serverSelector,omitempty"`
}

// IroncoreMetalMachineStatus defines the observed state of IroncoreMetalMachine
type IroncoreMetalMachineStatus struct {
	// Ready indicates the Machine infrastructure has been provisioned and is ready.
	// +optional
	Ready bool `json:"ready"`

	// FailureReason will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a succinct value suitable
	// for machine interpretation.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a more verbose string suitable
	// for logging and human consumption.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// IroncoreMetalMachine is the Schema for the ironcoremetalmachines API
type IroncoreMetalMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IroncoreMetalMachineSpec   `json:"spec,omitempty"`
	Status IroncoreMetalMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IroncoreMetalMachineList contains a list of IroncoreMetalMachine
type IroncoreMetalMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IroncoreMetalMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IroncoreMetalMachine{}, &IroncoreMetalMachineList{})
}
