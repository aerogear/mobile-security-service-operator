package mobilesecurityservicedb

import (
	"context"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/labels"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

var log = logf.Log.WithName("controller_mobilesecurityservicedb")

const (
	DEEPLOYMENT = "Deployment"
	PVC         = "PersistentVolumeClaim"
	SERVICE     = "Service"
)

// Add creates a new MobileSecurityServiceDB Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMobileSecurityServiceDB{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mobilesecurityservicedb-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	// Watch for changes to primary resource MobileSecurityServiceDB
	err = c.Watch(&source.Kind{Type: &mobilesecurityservicev1alpha1.MobileSecurityServiceDB{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	/** Watch for changes to secondary resources and create the owner MobileSecurityService **/

	//Deployment
	if err := watchDeployment(c); err != nil {
		return err
	}

	//Service
	if err := watchService(c); err != nil {
		return err
	}

	//PersistenceVolume
	if err := watchPersistenceVolumeClaim(c); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMobileSecurityServiceDB{}

// ReconcileMobileSecurityServiceDB reconciles a MobileSecurityServiceDB object
type ReconcileMobileSecurityServiceDB struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

//Update the object and reconcile it
func update(r *ReconcileMobileSecurityServiceDB, obj runtime.Object, reqLogger logr.Logger) (reconcile.Result, error) {
	err := r.client.Update(context.TODO(), obj)
	if err != nil {
		reqLogger.Error(err, "Failed to update Spec")
		return reconcile.Result{}, err
	}
	reqLogger.Info("Spec updated - return and create")
	return reconcile.Result{Requeue: true}, nil
}

func create(r *ReconcileMobileSecurityServiceDB, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceDB, reqLogger logr.Logger, kind string, err error) (reconcile.Result, error) {
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

func buildObject(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceDB, r *ReconcileMobileSecurityServiceDB, kind string) (runtime.Object, error) {
	reqLogger.Info("Check "+kind, "into the namespace", instance.Namespace)
	switch kind {
	case PVC:
		return r.buildPVCForDB(instance), nil
	case DEEPLOYMENT:
		return r.buildDBDeployment(instance), nil
	case SERVICE:
		return r.buildDBService(instance), nil
	default:
		msg := "Failed to recognize type of object" + kind + " into the Namespace " + instance.Namespace
		panic(msg)
	}
}

// Request object not found, could have been deleted after reconcile request.
// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
func fetch(r *ReconcileMobileSecurityServiceDB, reqLogger logr.Logger, err error) (reconcile.Result, error) {
	if errors.IsNotFound(err) {
		// Return and don't create
		reqLogger.Info("Mobile Security Service Database resource not found. Ignoring since object must be deleted")
		return reconcile.Result{}, nil
	}
	// Error reading the object - create the request.
	reqLogger.Error(err, "Failed to get Mobile Security Service Database")
	return reconcile.Result{}, err
}

// Reconcile reads that state of the cluster for a MobileSecurityServiceDB object and makes changes based on the state read
// and what is in the MobileSecurityServiceDB.Spec
// Note:
// The Controller will create the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMobileSecurityServiceDB) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Mobile Security Service Database")

	// Fetch the MobileSecurityServiceDB instance
	instance := &mobilesecurityservicev1alpha1.MobileSecurityServiceDB{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		return fetch(r, reqLogger, err)
	}

	reqLogger.Info("Checking if the DB deployment already exists, if not create a new one")
	deployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, deployment)
	if err != nil {
		return create(r, instance, reqLogger, DEEPLOYMENT, err)
	}

	reqLogger.Info("Checking if the DB Service already exists, if not create a new one")
	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, service)
	if err != nil {
		return create(r, instance, reqLogger, SERVICE, err)
	}

	reqLogger.Info("Checking if the DB PersistentVolumeClaim already exists, if not create a new one")
	pvc := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, pvc)
	if err != nil {
		return create(r, instance, reqLogger, PVC, err)
	}

	//Ensure the deployment size is the same as the spec
	reqLogger.Info("Ensuring the MobileSecurityServiceDB deployment size is the same as the spec")
	size := instance.Spec.Size
	if *deployment.Spec.Replicas != size {
		deployment.Spec.Replicas = &size
		return update(r, deployment, reqLogger)
	}

	reqLogger.Info("Updating the MobileSecurityServiceDB status with the pod names")
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(getDBLabels(instance.Name))
	listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOps, podList)
	if err != nil {
		reqLogger.Error(err, "Failed to list DB pods", "MobileSecurityServiceDB.Namespace", instance.Namespace, "MobileSecurityServiceDB.Name", instance.Name)
		return reconcile.Result{}, err
	}
	reqLogger.Info("Get DB pod names")
	podNames := utils.GetPodNames(podList.Items)
	reqLogger.Info("Update status.Nodes if needed regards MobileSecurityServiceDB")
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update MobileSecurityServiceDB status")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
