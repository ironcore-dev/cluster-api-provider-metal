// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

// ServerSelector specifies matching criteria for labels on Server.
// This is used to claim specific Server types for a Machine
type ServerSelector struct {
	// Key/value pairs of labels that must exist on a chosen Server
	// +optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`
}
