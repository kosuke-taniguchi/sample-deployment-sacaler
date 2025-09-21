package controller

import (
	"context"
	"fmt"
	"time"

	scalingv1 "github.com/kosuke-taniguchi/deployment-scaler/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// RBAC: manage CR and read/update Deployments
//+kubebuilder:rbac:groups=scaling.example.com,resources=deploymentscalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=scaling.example.com,resources=deploymentscalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=scaling.example.com,resources=deploymentscalers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch

type DeploymentScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *DeploymentScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var ds scalingv1.DeploymentScaler
	if err := r.Get(ctx, req.NamespacedName, &ds); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if !ds.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}
	// resolve target namespace (default to CR's)
	targetNS := ds.Namespace
	if ds.Spec.Target.Namespace != nil && *ds.Spec.Target.Namespace != "" {
		targetNS = *ds.Spec.Target.Namespace
	}
	targetName := ds.Spec.Target.Name
	// fetch target Deployment
	var dep appsv1.Deployment
	if err := r.Get(ctx, types.NamespacedName{Namespace: targetNS, Name: targetName}, &dep); err != nil {
		if apierrors.IsNotFound(err) {
			cond := metav1.Condition{
				Type:               "TargetAvailable",
				Status:             metav1.ConditionFalse,
				Reason:             "NotFound",
				Message:            fmt.Sprintf("Deployment %s/%s not found", targetNS, targetName),
				LastTransitionTime: metav1.Now(),
			}
			setCondition(&ds, cond)
			_ = r.Status().Update(ctx, &ds)
			// しばらくして再確認
			return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
		}
		return ctrl.Result{}, err
	}
	cur := int32(0)
	if dep.Spec.Replicas != nil {
		cur = *dep.Spec.Replicas
	}
	desired := ds.Spec.Replicas
	// 実際の状態とdesired stateが異なる場合はscale
	if cur != desired {
		dep.Spec.Replicas = &desired
		if err := r.Update(ctx, &dep); err != nil {
			return ctrl.Result{}, err
		}
		setCondition(&ds, metav1.Condition{
			Type:               "Synchronized",
			Status:             metav1.ConditionTrue,
			Reason:             "Scaled",
			Message:            fmt.Sprintf("Scaled %s/%s from %d to %d", targetNS, targetName, cur, desired),
			LastTransitionTime: metav1.Now(),
		})
	} else {
		setCondition(&ds, metav1.Condition{
			Type:               "Synchronized",
			Status:             metav1.ConditionTrue,
			Reason:             "InSync",
			Message:            "Deployment replicas already match",
			LastTransitionTime: metav1.Now(),
		})
	}
	if err := r.Status().Update(ctx, &ds); err != nil {
		logger.Error(err, "failed to update status")
	}
	return ctrl.Result{}, nil
}

func setCondition(ds *scalingv1.DeploymentScaler, cond metav1.Condition) {
	found := false
	for i, c := range ds.Status.Conditions {
		if c.Type == cond.Type {
			ds.Status.Conditions[i] = cond
			found = true
			break
		}
	}
	if !found {
		ds.Status.Conditions = append(ds.Status.Conditions, cond)
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scalingv1.DeploymentScaler{}).
		Named("deploymentscaler").
		Complete(r)
}
