// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// MetalMachineTemplateSpec defines the desired state of MetalMachineTemplate
type MetalMachineTemplateSpec struct {
	Template MetalMachineTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true

// MetalMachineTemplate is the Schema for the metalmachinetemplates API
type MetalMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MetalMachineTemplateSpec `json:"spec,omitempty"`
}

// MetalMachineTemplateResource defines the spec and metadata for MetalMachineTemplate supported by capi.
type MetalMachineTemplateResource struct {
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta clusterv1.ObjectMeta `json:"metadata,omitempty"`
	Spec       MetalMachineSpec     `json:"spec"`
}

// +kubebuilder:object:root=true

// MetalMachineTemplateList contains a list of MetalMachineTemplate
type MetalMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MetalMachineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MetalMachineTemplate{}, &MetalMachineTemplateList{})
}
