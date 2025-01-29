// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scope

import (
	"context"

	"github.com/go-logr/logr"
	infrav1 "github.com/ironcore-dev/cluster-api-provider-metal/api/v1alpha1"
	"github.com/ironcore-dev/metal-operator/api/v1alpha1"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// MachineScopeParams defines the input parameters used to create a new Scope.
type MachineScopeParams struct {
	Client       client.Client
	Logger       *logr.Logger
	Cluster      *clusterv1.Cluster
	Machine      *clusterv1.Machine
	MetalCluster *infrav1.MetalCluster
	MetalMachine *infrav1.MetalMachine
}

// MachineScope defines the basic context for an actuator to operate upon.
type MachineScope struct {
	*logr.Logger
	client       client.Client
	patchHelper  *patch.Helper
	Cluster      *clusterv1.Cluster
	Machine      *clusterv1.Machine
	MetalCluster *infrav1.MetalCluster
	MetalMachine *infrav1.MetalMachine
	ServerClaim  *v1alpha1.ServerClaim
}

// NewMachineScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewMachineScope(params MachineScopeParams) (*MachineScope, error) {
	if params.Client == nil {
		return nil, errors.New("Client is required when creating a MachineScope")
	}
	if params.Cluster == nil {
		return nil, errors.New("Cluster is required when creating a MachineScope")
	}
	if params.Machine == nil {
		return nil, errors.New("Machine is required when creating a MachineScope")
	}
	if params.MetalCluster == nil {
		return nil, errors.New("MetalCluster is required when creating a MachineScope")
	}
	if params.MetalMachine == nil {
		return nil, errors.New("MetalMachine is required when creating a MachineScope")
	}
	if params.Logger == nil {
		logger := log.FromContext(context.Background())
		params.Logger = &logger
	}

	machineScope := &MachineScope{
		Logger:       params.Logger,
		client:       params.Client,
		Cluster:      params.Cluster,
		Machine:      params.Machine,
		MetalCluster: params.MetalCluster,
		MetalMachine: params.MetalMachine,
	}

	helper, err := patch.NewHelper(params.MetalMachine, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	machineScope.patchHelper = helper

	return machineScope, nil
}

// SetReady sets the MetalMachine Ready Status.
func (m *MachineScope) SetReady() {
	m.MetalMachine.Status.Ready = true
}

// SetNotReady sets the MetalMachine Ready Status to false.
func (m *MachineScope) SetNotReady() {
	m.MetalMachine.Status.Ready = false
}

// SetFailureMessage sets the MetalMachine status failure message.
func (m *MachineScope) SetFailureMessage(v error) {
	m.MetalMachine.Status.FailureMessage = ptr.To(v.Error())
}

// SetFailureReason sets the MetalMachine status failure reason.
func (m *MachineScope) SetFailureReason(v string) {
	m.MetalMachine.Status.FailureReason = v
}

// HasFailed returns the failure state of the machine scope.
func (m *MachineScope) HasFailed() bool {
	return m.MetalMachine.Status.FailureReason != "" || m.MetalMachine.Status.FailureMessage != nil
}

// PatchObject persists the Machine configuration and status.
func (s *MachineScope) PatchObject() error {
	// always update the readyCondition.
	// TBD readyCondition

	return s.patchHelper.Patch(context.TODO(), s.MetalMachine)
}

// Close closes the current scope persisting the Machine configuration and status.
func (s *MachineScope) Close() error {
	return s.PatchObject()
}
