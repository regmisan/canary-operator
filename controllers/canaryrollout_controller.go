package controllers

import (
    "context"
    "strconv"
    "time"

    netv1 "k8s.io/api/networking/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/types"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log"

    rolloutv1alpha1 "github.com/regmisan/canary-operator/api/v1alpha1"
)

// CanaryRolloutReconciler reconciles a CanaryRollout object
type CanaryRolloutReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=rollout.example.com,resources=canaryrollouts,verbs=get;list;watch;update
// +kubebuilder:rbac:groups=rollout.example.com,resources=canaryrollouts/status,verbs=get;update
// +kubebuilder:rbac:groups="",resources=ingresses,verbs=get;update;patch

func (r *CanaryRolloutReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    logger := log.FromContext(ctx)

    // 1. Fetch the CanaryRollout
    var cr rolloutv1alpha1.CanaryRollout
    if err := r.Get(ctx, req.NamespacedName, &cr); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // 2. Determine next step index
    var nextIndex int32 = 0
    if cr.Status.CurrentStep != nil {
        nextIndex = *cr.Status.CurrentStep + 1
    }
    // If beyond last step, mark Completed and exit
    if int(nextIndex) >= len(cr.Spec.Steps) {
        if !cr.Status.Completed {
            cr.Status.Completed = true
            if err := r.Status().Update(ctx, &cr); err != nil {
                return ctrl.Result{}, err
            }
        }
        return ctrl.Result{}, nil
    }

    step := cr.Spec.Steps[nextIndex]

    // 3. Patch the canary Ingress annotation
    var ing netv1.Ingress
    key := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.CanaryIngress}
    if err := r.Get(ctx, key, &ing); err != nil {
        logger.Error(err, "fetching Ingress", "ingress", key)
        return ctrl.Result{}, err
    }
    if ing.Annotations == nil {
        ing.Annotations = make(map[string]string)
    }
    ing.Annotations["nginx.ingress.kubernetes.io/canary-weight"] = strconv.Itoa(int(step.Weight))
    if err := r.Update(ctx, &ing); err != nil {
        logger.Error(err, "patching Ingress annotation")
        return ctrl.Result{}, err
    }

    // 4. Update status.CurrentStep
    cr.Status.CurrentStep = &nextIndex
    if err := r.Status().Update(ctx, &cr); err != nil {
        return ctrl.Result{}, err
    }

    // 5. Requeue after PauseSeconds (if set)
    pause := time.Duration(30) * time.Second
    if step.PauseSeconds != nil {
        pause = time.Duration(*step.PauseSeconds) * time.Second
    }
    return ctrl.Result{RequeueAfter: pause}, nil
}

func (r *CanaryRolloutReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&rolloutv1alpha1.CanaryRollout{}).
        Complete(r)
}
