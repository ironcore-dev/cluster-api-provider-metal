// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"fmt"

	"github.com/ironcore-dev/cluster-api-provider-metal/internal/scope"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "github.com/ironcore-dev/cluster-api-provider-metal/api/v1alpha1"
)

// MetalClusterReconciler reconciles a MetalCluster object
type MetalClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

func (r *MetalClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	metalCluster := &infrav1.MetalCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, metalCluster); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Get owner cluster
	cluster, err := util.GetOwnerCluster(ctx, r.Client, metalCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cluster == nil {
		logger.Info("waiting for Cluster Controller to set OwnerRef on MetalCluster")
		return ctrl.Result{}, nil
	}

	logger = logger.WithValues("cluster", klog.KObj(cluster))
	ctx = ctrl.LoggerInto(ctx, logger)

	if annotations.IsPaused(cluster, metalCluster) {
		logger.Info("MetalCluster or owning Cluster is marked as paused, not reconciling")
		return ctrl.Result{}, nil
	}

	// Create the scope.
	clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client:         r.Client,
		Logger:         &logger,
		Cluster:        cluster,
		MetalCluster:   metalCluster,
		ControllerName: "metalcluster",
	})

	if err != nil {
		return reconcile.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}

	// Always close the scope when exiting this function, so we can persist any MetalCluster changes.
	defer func() {
		if err := clusterScope.Close(); err != nil && err == nil {
			logger.Error(err, "failed to close MetalCluster scope")
		}
	}()

	// Handle deleted clusters
	if !metalCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, clusterScope)
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, clusterScope)

}

func (r *MetalClusterReconciler) reconcileDelete(ctx context.Context, clusterScope *scope.ClusterScope) (reconcile.Result, error) {
	// We want to prevent deletion unless the owning cluster was flagged for deletion.
	if clusterScope.Cluster.DeletionTimestamp.IsZero() {
		clusterScope.Error(errors.New("deletion was requested but owning cluster wasn't deleted"), "Unable to delete MetalCluster")
		// We stop reconciling here. It will be triggered again once the owning cluster was deleted.
		return reconcile.Result{}, nil
	}

	clusterScope.Logger.V(4).Info("reconciling MetalCluster delete")
	// Deletion usually should be triggered through the deletion of the owning cluster.
	// If the MetalCluster was also flagged for deletion (e.g. deletion using the manifest file)
	// we should only allow to remove the finalizer when there are no MetalMachines left.
	machines, err := r.listMetalMachinesForCluster(ctx, clusterScope)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "could not retrieve metal machines for cluster %q", clusterScope.InfraClusterName())
	}

	// Requeue if there are one or more machines left.
	if len(machines) > 0 {
		clusterScope.Info("waiting for machines to be deleted", "remaining", len(machines))
		return ctrl.Result{RequeueAfter: infrav1.DefaultReconcilerRequeue}, nil
	}

	clusterScope.Info("cluster deleted successfully")
	ctrlutil.RemoveFinalizer(clusterScope.MetalCluster, infrav1.ClusterFinalizer)
	return ctrl.Result{}, nil
}

//nolint:unparam
func (r *MetalClusterReconciler) reconcileNormal(_ context.Context, clusterScope *scope.ClusterScope) (reconcile.Result, error) {
	clusterScope.Logger.Info("Reconciling MetalCluster")

	// If the MetalCluster doesn't have our finalizer, add it.
	ctrlutil.AddFinalizer(clusterScope.MetalCluster, infrav1.ClusterFinalizer)

	conditions.MarkTrue(clusterScope.MetalCluster, infrav1.MetalClusterReady)

	clusterScope.MetalCluster.Status.Ready = true

	return ctrl.Result{}, nil
}

func (r *MetalClusterReconciler) listMetalMachinesForCluster(ctx context.Context, clusterScope *scope.ClusterScope) ([]infrav1.MetalMachine, error) {
	var machineList infrav1.MetalMachineList
	err := r.List(ctx, &machineList, client.InNamespace(clusterScope.Namespace()), client.MatchingLabels{
		clusterv1.ClusterNameLabel: clusterScope.Name(),
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("listing machines", clusterScope.Name(), machineList.Items)
	return machineList.Items, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MetalClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.MetalCluster{}).
		WithEventFilter(predicates.ResourceNotPaused(mgr.GetScheme(), ctrl.LoggerFrom(ctx))).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(ctx, infrav1.GroupVersion.WithKind("MetalCluster"), mgr.GetClient(), &infrav1.MetalCluster{})),
		).
		Complete(r)
}
