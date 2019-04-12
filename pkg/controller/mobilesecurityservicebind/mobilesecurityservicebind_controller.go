package mobilesecurityservicebind

import (
	"context"
	"fmt"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_mobilesecurityservicebind")

const (
	CONFIGMAP = "ConfigMap"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new MobileSecurityServiceBind Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMobileSecurityServiceBind{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mobilesecurityservicebind-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MobileSecurityServiceBind
	err = c.Watch(&source.Kind{Type: &mobilesecurityservicev1alpha1.MobileSecurityServiceBind{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	/** Watch for changes to secondary resources and create the owner MobileSecurityService **/
	//ConfigMap
	if err := watchConfigMap(c); err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMobileSecurityServiceBind{}

//ReconcileMobileSecurityServiceBind reconciles a MobileSecurityServiceBind object
type ReconcileMobileSecurityServiceBind struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

//Build the object, cluster resource, and add the object in the queue to reconcile
func create(r *ReconcileMobileSecurityServiceBind, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, kind string, reqLogger logr.Logger, err error) (reconcile.Result, error) {
	obj, errBuildObject := buildObject(reqLogger, instance, r, kind)
	if errBuildObject != nil {
		return reconcile.Result{}, errBuildObject
	}
	if errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ", "kind", kind, "Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
		err = r.client.Create(context.TODO(), obj)
		if err != nil {
			reqLogger.Error(err, "Failed to create new ", "kind", kind, "Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
			return reconcile.Result{}, err
		}
		reqLogger.Info("Created successfully - return and create", "kind", kind, "Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
		return reconcile.Result{Requeue: true}, nil
	}
	reqLogger.Error(err, "Failed to build", "kind", kind, "Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
	return reconcile.Result{}, err
}

//Delete the object and reconcile it
func delete(r *ReconcileMobileSecurityServiceBind, obj runtime.Object, reqLogger logr.Logger ) (reconcile.Result, error) {
	reqLogger.Info("Deleting a new object ", "kind", obj.GetObjectKind())
	err := r.client.Delete(context.TODO(), obj)
	if err != nil {
		reqLogger.Error(err, "Failed to delete ", "kind", obj.GetObjectKind())
		return reconcile.Result{}, err
	}
	reqLogger.Info("Delete with successfully - return and requeue", "kind", obj.GetObjectKind())
	return reconcile.Result{Requeue: true}, nil
}

//Build Objects for MobileSecurityServiceBind
func buildObject(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, r *ReconcileMobileSecurityServiceBind, kind string) (runtime.Object, error) {
	reqLogger.Info("Building Object ", "kind", kind, "Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
	switch kind {
	case CONFIGMAP:
		return r.buildAppBindSDKConfigMap(instance), nil
	default:
		msg := "Failed to recognize type of object" + kind + " into the MobileSecurityServiceBind.Namespace " + instance.Namespace
		panic(msg)
	}
}

//Fetch MobileSecurityServiceBind instance
func fetch(r *ReconcileMobileSecurityServiceBind, reqLogger logr.Logger, err error) (reconcile.Result, error) {
	if errors.IsNotFound(err) {
		// Return and don't create
		reqLogger.Info("Mobile Security Service Bind resource not found. Ignoring since object must be deleted")
		return reconcile.Result{}, nil
	}
	// Error reading the object - create the request.
	reqLogger.Error(err, "Failed to get Mobile Security Service Bind")
	return reconcile.Result{}, err
}

// Reconcile reads that state of the cluster for a MobileSecurityServiceBind object and makes changes based on the state read
// and what is in the MobileSecurityServiceBind.Spec
// Note:
// The Controller will create the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMobileSecurityServiceBind) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MobileSecurityServiceBind")

	instance := &mobilesecurityservicev1alpha1.MobileSecurityServiceBind{}

	//Fetch the MobileSecurityServiceBind instance
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		return fetch(r, reqLogger, err)
	}

	//Check if the cluster host was added in the CR
	if len(instance.Spec.ClusterHost) < 1 || instance.Spec.ClusterHost == "{{clusterHost}}" {
		err := fmt.Errorf("Cluster Host IP was not found.")
		reqLogger.Error(err, "Please check its configuration. See https://github.com/aerogear/mobile-security-service-operator#configuring .")
		return reconcile.Result{}, nil
	}

	reqLogger.Info("Checking if the SDKConfigMap already exists, if not create a new one")
	configMap := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: getConfigMapName(instance), Namespace: instance.Namespace}, configMap)
	if err != nil {
		return create(r, instance,  CONFIGMAP, reqLogger, err)
	}

	app, _ := getAppFromServiceByRestApi(instance, reqLogger)
	if len(app.ID) > 0 {
		if app.AppName != instance.Spec.AppName {
			//Update the name by the REST API
			reqLogger.Info("Calling the Rest api to update the app name ...", "App.newName:", instance.Spec.AppName, "App.oldName", app.AppName, "App.appID:", instance.Spec.AppId , "App.iD", app.ID)
			return updateAppNameByRestAPI(instance, app, reqLogger)
		}
	} else {
		reqLogger.Info("Calling the Rest API for create new app...")
		return createAppByRestAPI(instance, reqLogger)
	}

	reqLogger.Info("Updating the SDKConfigMap status ...")
	configMapStatus := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: getConfigMapName(instance), Namespace: instance.Namespace}, configMapStatus)
	if err != nil {
		reqLogger.Error(err, "Failed to update SDKConfigMap status")
		return reconcile.Result{}, err
	}
	if !reflect.DeepEqual(configMapStatus.Name, instance.Status.ConfigMap) {
		instance.Status.ConfigMap = configMapStatus.Name
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update SDKConfigMap status for MobileSecurityServiceBind")
			return reconcile.Result{}, err
		}
	}

	reqLogger.Info("Updating the Bind App status ...")
	bindApp, err := getAppFromServiceByRestApi(instance, reqLogger)
	if err != nil || len(bindApp.ID) < 0 {
		reqLogger.Error(err, "Failed to update Bind App status")
		return reconcile.Result{}, err
	}
	if !reflect.DeepEqual(bindApp.ID, instance.Status.Bind) {
		instance.Status.Bind = bindApp.ID
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Bind App status for MobileSecurityServiceBind")
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}
