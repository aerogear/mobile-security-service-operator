package mobilesecurityservice

import (
	"context"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	"github.com/go-logr/logr"
	"k8s.io/api/extensions/v1beta1"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	CONFIGMAP     = "ConfigMap"
	DEEPLOYMENT   = "Deployment"
	SERVICE       = "Service"
	INGRESS       = "Ingress"
)

var log = logf.Log.WithName("controller_mobilesecurityservice")

// Add creates a new MobileSecurityService Controller and adds it to the Manager.
// The Manager will set fields on the Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// Returns the a new Reconciler for this operator and controller
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMobileSecurityService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// Add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mobilesecurityservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MobileSecurityService
	err = c.Watch(&source.Kind{Type: &mobilesecurityservicev1alpha1.MobileSecurityService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	/** Watch for changes to secondary resources and create the owner MobileSecurityService **/

	//ConfigMap
	if err := watchConfigMap(c); err != nil {
		return err
	}

	//Deployment
	if err := watchDeployment(c); err != nil {
		return err
	}

	//Service
	if err := watchService(c); err != nil {
		return err
	}

	//Ingress
	if err:= watchIngress(c); err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMobileSecurityService{}

// ReconcileMobileSecurityService reconciles a MobileSecurityService object
type ReconcileMobileSecurityService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

//Update the object and reconcile it
func update(r *ReconcileMobileSecurityService, obj runtime.Object, reqLogger logr.Logger) (reconcile.Result, error) {
	err := r.client.Update(context.TODO(), obj)
	if err != nil {
		reqLogger.Error(err, "Failed to update Spec")
		return reconcile.Result{}, err
	}
	reqLogger.Info("Spec updated - return and create")
	return reconcile.Result{Requeue: true}, nil
}

func create(r *ReconcileMobileSecurityService, instance *mobilesecurityservicev1alpha1.MobileSecurityService, reqLogger logr.Logger, kind string, err error) (reconcile.Result, error) {
	obj, errBuildObject := buildObject(reqLogger, instance, r, kind)
	if errBuildObject != nil {
		return reconcile.Result{}, errBuildObject
	}
	if errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ", "kind", kind, "Namespace", instance.Namespace)
		err = r.client.Create(context.TODO(), obj)
		if err != nil {
			reqLogger.Error(err, "Failed to create new ", "kind", kind, "Namespace", instance.Namespace)
			return reconcile.Result{}, err
		}
		reqLogger.Info("Created successfully - return and create", "kind", kind, "Namespace", instance.Namespace)
		return reconcile.Result{Requeue: true}, nil
	}
	reqLogger.Error(err, "Failed to get", "kind", kind, "Namespace", instance.Namespace)
	return reconcile.Result{}, err

}

func buildObject(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityService, r *ReconcileMobileSecurityService, kind string) (runtime.Object, error) {
	reqLogger.Info("Check "+kind, "into the namespace", instance.Namespace)
	switch kind {
	case CONFIGMAP:
		return r.buildAppConfigMap(instance), nil
	case DEEPLOYMENT:
		return r.buildAppDeployment(instance), nil
	case SERVICE:
		return r.buildAppService(instance), nil
	case INGRESS:
		return r.buildAppIngress(instance), nil
	default:
		msg := "Failed to recognize type of object" + kind + " into the Namespace " + instance.Namespace
		panic(msg)
	}
}

func fetch(r *ReconcileMobileSecurityService, reqLogger logr.Logger, err error) (reconcile.Result, error) {
	if errors.IsNotFound(err) {
		// Return and don't create
		reqLogger.Info("Mobile Security Service App resource not found. Ignoring since object must be deleted")
		return reconcile.Result{}, nil
	}
	// Error reading the object - create the request.
	reqLogger.Error(err, "Failed to get Mobile Security Service App")
	return reconcile.Result{}, err
}

// Reconcile reads that state of the cluster for a MobileSecurityService object and makes changes based on the state read
// and what is in the MobileSecurityService.Spec
// Note:
// The Controller will create the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMobileSecurityService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Mobile Security Service App")

	instance := &mobilesecurityservicev1alpha1.MobileSecurityService{}

	//Fetch the MobileSecurityService instance
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		return fetch(r, reqLogger, err)
	}

	reqLogger.Info("Checkiing if the ConfigMap already exists, if not create a new one")
	configMap := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.ConfigMapName, Namespace: instance.Namespace}, configMap)
	if err != nil {
		return create(r, instance, reqLogger, CONFIGMAP, err)
	}

	reqLogger.Info("Checking if the deployment already exists, if not create a new one")
	deployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, deployment)
	if err != nil {
		return create(r, instance, reqLogger, DEEPLOYMENT, err)
	}

	reqLogger.Info("Checking if the service already exists, if not create a new one")
	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, service)
	if err != nil {
		return create(r, instance, reqLogger, SERVICE, err)
	}

	reqLogger.Info("Checking if the ingress already exists, if not create a new one")
	ingress := &v1beta1.Ingress{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, ingress)
	if err != nil {
		return create(r, instance, reqLogger, INGRESS, err)
	}

	reqLogger.Info("Ensuring the Mobile Security Service deployment size is the same as the spec")
	size := instance.Spec.Size
	if *deployment.Spec.Replicas != size {
		deployment.Spec.Replicas = &size
		return update(r, deployment, reqLogger)
	}

	reqLogger.Info("Updating the MobileSecurityService status with the pod names")
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(getAppLabels(instance.Name))
	listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOps, podList)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods", "MobileSecurityService.Namespace", instance.Namespace, "MobileSecurityService.Name", instance.Name)
		return reconcile.Result{}, err
	}
	reqLogger.Info("Get pod names")
	podNames := utils.GetPodNames(podList.Items)
	reqLogger.Info("Update status.Nodes if needed")
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update MobileSecurityService status")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
