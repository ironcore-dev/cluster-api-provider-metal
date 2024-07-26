// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"

	"github.com/ironcore-dev/cluster-api-provider-metal/internal/scope"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	infrav1 "github.com/ironcore-dev/cluster-api-provider-metal/api/v1alpha1"
)

// MetalMachineReconciler reconciles a MetalMachine object
type MetalMachineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines/finalizers,verbs=update
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinedeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinesets,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=kubeadmcontrolplanes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *MetalMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the MetalMachine.
	metalMachine := &infrav1.MetalMachine{}
	err := r.Get(ctx, req.NamespacedName, metalMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Fetch the Machine.
	machine, err := util.GetOwnerMachine(ctx, r.Client, metalMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		logger.Info("Machine Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	logger = logger.WithValues("machine", klog.KObj(machine))

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		logger.Info("Machine is missing cluster label or cluster does not exist")
		return ctrl.Result{}, nil
	}

	if annotations.IsPaused(cluster, metalMachine) {
		logger.Info("MetalMachine or linked Cluster is marked as paused, not reconciling")
		return ctrl.Result{}, nil
	}

	logger = logger.WithValues("cluster", klog.KObj(cluster))

	metalClusterName := client.ObjectKey{
		Namespace: metalMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}

	metalCluster := &infrav1.MetalCluster{}
	if err := r.Client.Get(ctx, metalClusterName, metalCluster); err != nil {
		if apierrors.IsNotFound(err) || !metalCluster.Status.Ready {
			logger.Info("MetalCluster is not available yet")
			return ctrl.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Create the cluster scope.
	clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client:         r.Client,
		Logger:         &logger,
		Cluster:        cluster,
		MetalCluster:   metalCluster,
		ControllerName: "metalcluster",
	})

	if err != nil {
		return reconcile.Result{}, errors.Errorf("failed to create cluster scope: %+v", err)
	}

	// Create the machine scope
	machineScope, err := scope.NewMachineScope(scope.MachineScopeParams{
		Client:       r.Client,
		Cluster:      cluster,
		Machine:      machine,
		MetalCluster: metalCluster,
		MetalMachine: metalMachine,
	})

	if err != nil {
		return reconcile.Result{}, errors.Errorf("failed to create machine scope: %+v", err)
	}

	// Always close the scope when exiting this function, so we can persist any MetalMachine changes.
	// TODO: revisit side effects of closure errors
	defer func() {
		if err := machineScope.Close(); err != nil {
			logger.Error(err, "failed to close MetalMachine scope")
		}
	}()

	// Return early if the object or Cluster is paused.
	if annotations.IsPaused(cluster, metalMachine) {
		logger.Info("MetalMachine or linked Cluster is marked as paused. Won't reconcile normally")
		return reconcile.Result{}, nil
	}

	// Handle deleted machines
	if !metalMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, machineScope)
	}

	// Handle non-deleted machines
	return r.reconcileNormal(ctx, machineScope, clusterScope)
}

// SetupWithManager sets up the controller with the Manager.
func (r *MetalMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.MetalMachine{}).
		Watches(
			&clusterv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrav1.GroupVersion.WithKind("MetalMachine"))),
		).
		Complete(r)
}

// TODO: remove nolint tag
//
//nolint:unparam
func (r *MetalMachineReconciler) reconcileDelete(_ context.Context, machineScope *scope.MachineScope) (ctrl.Result, error) {
	machineScope.Logger.Info("Handling deleted MetalMachine")

	// insert ServerClaim deletion logic here

	// ServerClaim is being deleted
	return reconcile.Result{RequeueAfter: infrav1.DefaultReconcilerRequeue}, nil
}

// TODO: remove nolint tag
//
//nolint:unparam
func (r *MetalMachineReconciler) reconcileNormal(_ context.Context, machineScope *scope.MachineScope, clusterScope *scope.ClusterScope) (reconcile.Result, error) {
	clusterScope.Logger.V(4).Info("Reconciling MetalMachine")

	// If the MetalMachine is in an error state, return early.
	if machineScope.HasFailed() {
		machineScope.Info("Error state detected, skipping reconciliation")
		return ctrl.Result{}, nil
	}

	if !machineScope.Cluster.Status.InfrastructureReady {
		machineScope.Info("Cluster infrastructure is not ready yet")
		// TBD: update conditions
		return ctrl.Result{}, nil
	}

	// Make sure bootstrap data is available and populated.
	if machineScope.Machine.Spec.Bootstrap.DataSecretName == nil {
		machineScope.Info("Bootstrap data secret reference is not yet available")
		// TBD: update conditions
		return ctrl.Result{}, nil
	}

	// TBD add finalizer

	// Get or create the ServerClaim.
	// TBD
	machineScope.Info("Creating ServerClaim", "claim", machineScope.MetalMachine.Name)

	machineScope.SetReady()
	machineScope.Logger.Info("MetalMachine is ready")

	return reconcile.Result{}, nil
}
