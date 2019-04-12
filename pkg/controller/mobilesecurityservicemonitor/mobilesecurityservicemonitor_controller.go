package mobilesecurityservicemonitor

import (
	"context"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"
)

var log = logf.Log.WithName("controller_mobilesecurityservicemonitor")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new MobileSecurityServiceMonitor Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMobileSecurityServiceMonitor{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mobilesecurityservicemonitor-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MobileSecurityServiceMonitor
	err = c.Watch(&source.Kind{Type: &mobilesecurityservicev1alpha1.MobileSecurityServiceMonitor{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch all changes regards the MobileSecurityServiceBind in the cluster
	err = c.Watch(&source.Kind{Type: &mobilesecurityservicev1alpha1.MobileSecurityServiceBind{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMobileSecurityServiceMonitor{}

// ReconcileMobileSecurityServiceMonitor reconciles a MobileSecurityServiceMonitor object
type ReconcileMobileSecurityServiceMonitor struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

func fetch(r *ReconcileMobileSecurityServiceMonitor, reqLogger logr.Logger, err error) (reconcile.Result, error) {
	if errors.IsNotFound(err) {
		// Return and don't create
		reqLogger.Info("Mobile Security Service Monitor resource not found. Ignoring since object must be deleted")
		return reconcile.Result{}, nil
	}
	// Error reading the object - create the request.
	reqLogger.Error(err, "Failed to get Mobile Security Service Monitor")
	return reconcile.Result{}, err
}

// Reconcile reads that state of the cluster for a MobileSecurityServiceMonitor object and makes changes based on the state read
// and what is in the MobileSecurityServiceMonitor.Spec
// Note:
// The Controller will create the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMobileSecurityServiceMonitor) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MobileSecurityServiceMonitor")

	instance := &mobilesecurityservicev1alpha1.MobileSecurityServiceMonitor{}

	//Fetch the MobileSecurityServiceMonitor instance
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		return fetch(r, reqLogger, err)
	}

	filter := getListOptionsToFilterResources(instance, reqLogger)

	reqLogger.Info("Checking bind applied ...")
	bindList := &mobilesecurityservicev1alpha1.MobileSecurityServiceBindList{}
	err = r.client.List(context.TODO(), &filter, bindList)
	if err == nil {
		reqLogger.Info("Total of binds by specs defined found ...", "total", len(bindList.Items))

		//Check all apps and if has no bind then remove
		allApps, err := service.GetAllAppsFromServiceByRestApi(instance.Spec.Protocol, instance.Spec.ClusterHost, instance.Spec.HostSufix, reqLogger)
		if err != nil {
			reqLogger.Error(err, "Failed to get all apps")
			return reconcile.Result{}, err
		}

		reqLogger.Info("Total of apps found in the Rest Service API ...", "total", len(allApps))

		// Looking for the apps which are in the Service but has no longer bind in the cluster
		var unbindAppList []string
		for i := 0; i < len(allApps); i++ {
			appID := allApps[i].AppID
			found := false
			for i := 0; i < len(bindList.Items); i++ {
				bindAppID := bindList.Items[i].Spec.AppId
				if bindAppID == appID {
					found = true
					break
				}
			}
			if found == false {
				unbindAppList = append(unbindAppList, appID)
			}
		}

		reqLogger.Info("Total of unbind apps found to be deleted ...", "total", len(unbindAppList))

		for i := 0; i < len(unbindAppList); i++ {
			id := unbindAppList[i]
			reqLogger.Info("Calling Rest Service API to delete app...", "app.id", id)
			err := service.DeleteAppFromServiceByRestAPI(instance.Spec.Protocol, instance.Spec.ClusterHost, instance.Spec.HostSufix, id, reqLogger)
			if err != nil {
				reqLogger.Error(err, "Failed to delete unbind app with id", "App.id", id)
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{RequeueAfter: time.Second * 10}, nil
}

