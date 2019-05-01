package mobilesecurityserviceunbind

import (
	"context"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	routev1 "github.com/openshift/api/route/v1"
)

var log = logf.Log.WithName("controller_mobilesecurityserviceunbind")

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
	if err = c.Watch(&source.Kind{Type: &mobilesecurityservicev1alpha1.MobileSecurityServiceUnbind{}}, &handler.EnqueueRequestForObject{});err != nil  {
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


// Reconcile reads that state of the cluster for a ReconcileMobileSecurityServiceUnbind object and makes changes based on the state read
// and what is in the ReconcileMobileSecurityServiceUnbind.Spec
// Note:
// The Controller will create the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMobileSecurityServiceUnbind) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MobileSecurityServiceUnbind")

	instance := &mobilesecurityservicev1alpha1.MobileSecurityServiceUnbind{}

	//Fetch the MobileSecurityService instance
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		return fetch(r, reqLogger, err)
	}

	//Check specs
	if !hasSpecs(instance, reqLogger) {
		return reconcile.Result{Requeue: true}, nil
	}

	reqLogger.Info("Checking if the route already exists ...")
	route := &routev1.Route{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: "mobile-security-service", Namespace: "mobile-security-service-operator"}, route); err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("Checking for service ...")
	serviceInstance := &mobilesecurityservicev1alpha1.MobileSecurityService{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: "mobile-security-service", Namespace: "mobile-security-service-operator"}, serviceInstance); err != nil {
		return reconcile.Result{}, err
	}

	//Get the REST Service Endpoint
	serviceAPI := service.GetServiceAPIURL(route, serviceInstance)

	//Check if App is UnBind in the REST Service, if not then unbind it
	if app, err := fetchBindAppRestServiceByAppID(serviceAPI, instance, reqLogger); err == nil {
		if hasApp(app) {
			if err := service.DeleteAppFromServiceByRestAPI(serviceAPI,  app.ID, reqLogger); err != nil {
				reqLogger.Error(err, "Failed to delete unbind app with id", "App.id",  app.ID)
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		}
	}

	//Update status for UnBindStatus
	if err := r.updateUnbindStatus(serviceAPI, instance, reqLogger); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
