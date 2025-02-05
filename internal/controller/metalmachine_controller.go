// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/ironcore-dev/cluster-api-provider-ironcore-metal/internal/scope"
	"github.com/ironcore-dev/controller-utils/clientutils"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"

	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	infrav1alpha1 "github.com/ironcore-dev/cluster-api-provider-ironcore-metal/api/v1alpha1"
	metalv1alpha1 "github.com/ironcore-dev/metal-operator/api/v1alpha1"
)

// MetalMachineReconciler reconciles a MetalMachine object
type MetalMachineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	MetalMachineFinalizer        = "infrastructure.cluster.x-k8s.io/metalmachine"
	DefaultIgnitionSecretKeyName = "ignition"
)

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines/finalizers,verbs=update
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinedeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinesets,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=kubeadmcontrolplanes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=serverclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *MetalMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the MetalMachine.
	metalMachine := &infrav1alpha1.MetalMachine{}
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

	metalCluster := &infrav1alpha1.MetalCluster{}
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
		For(&infrav1alpha1.MetalMachine{}).
		Watches(
			&clusterapiv1beta1.Machine{},
			handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrav1alpha1.GroupVersion.WithKind("MetalMachine"))),
		).
		Complete(r)
}

func (r *MetalMachineReconciler) reconcileDelete(ctx context.Context, machineScope *scope.MachineScope) (ctrl.Result, error) {
	machineScope.Logger.Info("Deleting MetalMachine")

	// insert ServerClaim deletion logic here

	if modified, err := clientutils.PatchEnsureNoFinalizer(ctx, r.Client, machineScope.MetalMachine, MetalMachineFinalizer); !apierrors.IsNotFound(err) || modified {
		return ctrl.Result{}, err
	}
	machineScope.Logger.Info("Ensured that the finalizer has been removed")

	return reconcile.Result{RequeueAfter: infrav1alpha1.DefaultReconcilerRequeue}, nil
}

func (r *MetalMachineReconciler) reconcileNormal(ctx context.Context, machineScope *scope.MachineScope, clusterScope *scope.ClusterScope) (reconcile.Result, error) {
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

	if modified, err := clientutils.PatchEnsureFinalizer(ctx, r.Client, machineScope.MetalMachine, MetalMachineFinalizer); err != nil || modified {
		return ctrl.Result{}, err
	}
	machineScope.Logger.Info("Ensured finalizer has been added")

	// Fetch the bootstrap data secret.
	bootstrapSecret := &corev1.Secret{}
	secretName := types.NamespacedName{
		Namespace: machineScope.Machine.Namespace,
		Name:      *machineScope.Machine.Spec.Bootstrap.DataSecretName,
	}
	if err := r.Client.Get(ctx, secretName, bootstrapSecret); err != nil {
		machineScope.Error(err, "failed to get bootstrap data secret")
		return ctrl.Result{}, err
	}

	machineScope.Info("Creating IgnitionSecret", "Secret", machineScope.MetalMachine.Name)
	ignitionSecret, err := r.applyIgnitionSecret(ctx, machineScope.Logger, machineScope.MetalMachine, bootstrapSecret)
	if err != nil {
		machineScope.Error(err, "failed to create or patch ignition secret")
		return ctrl.Result{}, err
	}

	machineScope.Info("Creating ServerClaim", "ServerClaim", machineScope.MetalMachine.Name)
	serverClaim, err := r.applyServerClaim(ctx, machineScope.Logger, machineScope.MetalMachine, ignitionSecret)
	if err != nil {
		machineScope.Error(err, "failed to create or patch ServerClaim")
		return ctrl.Result{}, err
	}

	bound, _ := r.ensureServerClaimBound(ctx, serverClaim)
	if !bound {
		machineScope.Info("Waiting for ServerClaim to be Bound")
		return ctrl.Result{
			RequeueAfter: infrav1alpha1.DefaultReconcilerRequeue,
		}, nil
	}

	machineScope.Info("Patching ProviderID in MetalMachine")
	if err := r.patchMetalMachineProviderID(ctx, machineScope.Logger, machineScope.MetalMachine, serverClaim); err != nil {
		machineScope.Error(err, "failed to patch the MetalMachine with providerid")
		return ctrl.Result{}, err
	}

	machineScope.SetReady()
	machineScope.Logger.Info("MetalMachine is ready")

	return reconcile.Result{}, nil
}

func (r *MetalMachineReconciler) applyIgnitionSecret(ctx context.Context, log *logr.Logger, metalmachine *infrav1alpha1.MetalMachine, capidatasecret *corev1.Secret) (*corev1.Secret, error) {
	dataSecret := capidatasecret
	findAndReplaceIgnition(metalmachine, dataSecret)

	secretObj := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("ignition-%s", capidatasecret.Name),
			Namespace: capidatasecret.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		Data: map[string][]byte{
			DefaultIgnitionSecretKeyName: dataSecret.Data["value"],
		},
	}

	if err := controllerutil.SetControllerReference(capidatasecret, secretObj, r.Client.Scheme()); err != nil {
		return nil, fmt.Errorf("failed to set ControllerReference: %w", err)
	}

	opResult, err := controllerutil.CreateOrPatch(ctx, r.Client, secretObj, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create or patch the IgnitionSecret: %w", err)
	}
	log.Info("Created or Patched IgnitionSecret", "IgnitionSecret", secretObj.Name, "Operation", opResult)

	return secretObj, nil
}

func (r *MetalMachineReconciler) applyServerClaim(ctx context.Context, log *logr.Logger, metalmachine *infrav1alpha1.MetalMachine, ignitionsecret *corev1.Secret) (*metalv1alpha1.ServerClaim, error) {
	serverClaimObj := &metalv1alpha1.ServerClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      metalmachine.Name,
			Namespace: metalmachine.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: metalv1alpha1.GroupVersion.String(),
			Kind:       "ServerClaim",
		},
		Spec: metalv1alpha1.ServerClaimSpec{
			Power: metalv1alpha1.PowerOn,
			IgnitionSecretRef: &corev1.LocalObjectReference{
				Name: ignitionsecret.Name,
			},
			Image:          metalmachine.Spec.Image,
			ServerSelector: metalmachine.Spec.ServerSelector,
		},
	}

	if err := controllerutil.SetControllerReference(metalmachine, serverClaimObj, r.Client.Scheme()); err != nil {
		return nil, fmt.Errorf("failed to set ControllerReference: %w", err)
	}

	opResult, err := controllerutil.CreateOrPatch(ctx, r.Client, serverClaimObj, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create or patch ServerClaim: %w", err)
	}
	log.Info("Created or Patched ServerClaim", "ServerClaim", serverClaimObj.Name, "Operation", opResult)

	return serverClaimObj, nil
}

func (r *MetalMachineReconciler) patchMetalMachineProviderID(ctx context.Context, log *logr.Logger, metalmachine *infrav1alpha1.MetalMachine, serverClaim *metalv1alpha1.ServerClaim) error {
	providerID := fmt.Sprintf("metal://%s/%s", serverClaim.Namespace, serverClaim.Name)

	patch := client.MergeFrom(metalmachine.DeepCopy())
	metalmachine.Spec.ProviderID = &providerID

	if err := r.Client.Patch(ctx, metalmachine, patch); err != nil {
		log.Error(err, "failed to patch MetalMachine with ProviderID")
		return err
	}

	log.Info("Successfully patched MetalMachine with ProviderID", "ProviderID", providerID)
	return nil
}

func (r *MetalMachineReconciler) ensureServerClaimBound(ctx context.Context, serverClaim *metalv1alpha1.ServerClaim) (bool, error) {
	claimObj := &metalv1alpha1.ServerClaim{}
	if err := r.Get(ctx, client.ObjectKeyFromObject(serverClaim), claimObj); err != nil {
		return false, err
	}

	if claimObj.Status.Phase != metalv1alpha1.PhaseBound {
		return false, nil
	}
	return true, nil
}

func findAndReplaceIgnition(metalmachine *infrav1alpha1.MetalMachine, capidatasecret *corev1.Secret) {
	data := capidatasecret.Data["value"]

	// replace $${METAL_HOSTNAME} with machine name
	hostname := "%24%24%7BMETAL_HOSTNAME%7D"
	modifiedData := strings.ReplaceAll(string(data), hostname, metalmachine.Name)

	capidatasecret.Data["value"] = []byte(modifiedData)
}
