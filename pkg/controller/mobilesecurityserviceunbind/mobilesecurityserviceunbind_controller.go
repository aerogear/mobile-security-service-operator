package mobilesecurityserviceunbind

import (
	"context"
	"fmt"
	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"github.com/go-logr/logr"
	"reflect"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_mobilesecurityserviceunbind")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new MobileSecurityServiceUnbind Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMobileSecurityServiceUnbind{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mobilesecurityserviceunbind-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MobileSecurityServiceUnbind
	err = c.Watch(&source.Kind{Type: &mobilesecurityservicev1alpha1.MobileSecurityServiceUnbind{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMobileSecurityServiceUnbind{}

// ReconcileMobileSecurityServiceUnbind reconciles a MobileSecurityServiceUnbind object
type ReconcileMobileSecurityServiceUnbind struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

//Fetch ReconcileMobileSecurityServiceUnbind instance
func fetch(r *ReconcileMobileSecurityServiceUnbind, reqLogger logr.Logger, err error) (reconcile.Result, error) {
	if errors.IsNotFound(err) {
		// Return and don't create
		reqLogger.Info("Mobile Security Service Unbind resource not found. Ignoring since object must be deleted")
		return reconcile.Result{}, nil
	}
	// Error reading the object - create the request.
	reqLogger.Error(err, "Failed to get Mobile Security Service Unbind")
	return reconcile.Result{}, err
}

// Reconcile reads that state of the cluster for a ReconcileMobileSecurityServiceUnbind object and makes changes based on the state read
// and what is in the ReconcileMobileSecurityServiceUnbind.Spec
// Note:
// The Controller will create the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMobileSecurityServiceUnbind) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MobileSecurityServiceUnbind")

	// Fetch the MobileSecurityServiceUnbind instance
	instance := &mobilesecurityservicev1alpha1.MobileSecurityServiceUnbind{}

	//Fetch the MobileSecurityServiceBind instance
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		return fetch(r, reqLogger, err)
	}

	//Check if the cluster host was added in the CR
	if len(instance.Spec.ClusterHost) < 1 || instance.Spec.ClusterHost == "{{clusterHost}}" {
		err := fmt.Errorf("Cluster Host IP was not found.")
		reqLogger.Error(err, "Please check its configuration. See https://github.com/aerogear/mobile-security-service-operator#unbinding-an-app .")
		return reconcile.Result{}, nil
	}

	reqLogger.Info("Calling Rest API to get app ID to delete ...")
	bindApp, err := service.GetAppFromServiceByRestApi(instance.Spec.Protocol, instance.Spec.ClusterHost, instance.Spec.HostSufix, instance.Spec.AppId, reqLogger)
	if err != nil && len(bindApp.ID) < 1 {
		reqLogger.Error(err, "Failed to update get App ID to delete")
		return reconcile.Result{}, err
	}

	reqLogger.Info("Calling Rest Service API to delete app by ID ...", "app.id", bindApp.ID)
	err = service.DeleteAppFromServiceByRestAPI(instance.Spec.Protocol, instance.Spec.ClusterHost, instance.Spec.HostSufix,  bindApp.ID, reqLogger)
	if err != nil {
		reqLogger.Error(err, "Failed to delete unbind app with id", "App.id",  bindApp.ID)
		return reconcile.Result{}, err
	}

	reqLogger.Info("Updating the Unbind App status ...")
	bindAppCheck, err := service.GetAppFromServiceByRestApi(instance.Spec.Protocol, instance.Spec.ClusterHost, instance.Spec.HostSufix, instance.Spec.AppId, reqLogger)
	if err != nil || len(bindAppCheck.ID) > 1 {
		reqLogger.Error(err, "Failed to update Unbind App status")
		return reconcile.Result{}, err
	}
	if !reflect.DeepEqual(instance.Spec.AppId, instance.Status.Unbind) {
		instance.Status.Unbind = instance.Spec.AppId
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Unbind app status for MobileSecurityServiceUnbind")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
