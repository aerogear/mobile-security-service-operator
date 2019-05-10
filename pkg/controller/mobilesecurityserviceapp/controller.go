package mobilesecurityserviceapp

import (
	"context"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/models"
	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
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

var log = logf.Log.WithName("controller_mobilesecurityserviceapp")

const CONFIGMAP = "ConfigMap"

// Add creates a new MobileSecurityServiceApp Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMobileSecurityServiceApp{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mobilesecurityserviceapp-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MobileSecurityServiceApp
	err = c.Watch(&source.Kind{Type: &mobilesecurityservicev1alpha1.MobileSecurityServiceApp{}}, &handler.EnqueueRequestForObject{})
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

var _ reconcile.Reconciler = &ReconcileMobileSecurityServiceApp{}

//ReconcileMobileSecurityServiceApp reconciles a MobileSecurityServiceApp object
type ReconcileMobileSecurityServiceApp struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

//Build the object, cluster resource, and add the object in the queue to reconcile
func (r *ReconcileMobileSecurityServiceApp) create(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, kind string, serviceURL string, reqLogger logr.Logger, err error) (reconcile.Result, error) {
	obj, errBuildObject := r.buildFactory(reqLogger, instance, kind, serviceURL)
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

//buildFactory will return the resource according to the kind defined
func (r *ReconcileMobileSecurityServiceApp) buildFactory(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, kind string, serviceURL string) (runtime.Object, error) {
	reqLogger.Info("Building Object ", "kind", kind, "MobileSecurityServiceApp.Namespace", instance.Namespace, "MobileSecurityServiceApp.Name", instance.Name)
	switch kind {
	case CONFIGMAP:
		return r.buildAppSDKConfigMap(instance, serviceURL), nil
	default:
		msg := "Failed to recognize type of object" + kind + " into the MobileSecurityServiceApp.Namespace " + instance.Namespace
		panic(msg)
	}
}

// Reconcile reads that state of the cluster for a MobileSecurityServiceApp object and makes changes based on the state read
// and what is in the MobileSecurityServiceApp.Spec
// Note:
// The Controller will create the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMobileSecurityServiceApp) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MobileSecurityServiceApp")

	// FIXME: Check if is a valid namespace
	// We should not checked if the namespace is valid or not. It is an workaround since currently is not possible watch/cache a List of Namespaces
	// The impl to allow do it is done and merged in the master branch of the lib but not released in an stable version. It should be removed when this feature be impl.
	// See the PR which we are working on to update the deps and have this feature: https://github.com/operator-framework/operator-sdk/pull/1388
	if isValidNamespace, err := utils.IsValidAppNamespace(request.Namespace); err != nil || isValidNamespace == false {
		// Stop reconcile
		envVar, _ := utils.GetAppNamespaces()
		reqLogger.Error(err, "Unable to reconcile Mobile Security Service App", "Request.Namespace", request.Namespace, "isValidNamespace", isValidNamespace, "EnvVar.APP_NAMESPACES", envVar)
		return reconcile.Result{}, nil
	}

	reqLogger.Info("Valid namespace for MobileSecurityServiceApp", "Namespace", request.Namespace)
	reqLogger.Info("Start Reconciling MobileSecurityServiceApp ...")

	instance := &mobilesecurityservicev1alpha1.MobileSecurityServiceApp{}

	//Fetch the MobileSecurityService instance
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		return fetch(r, reqLogger, err)
	}

	reqLogger.Info("Checking for service instance ...")
	serviceInstance := &mobilesecurityservicev1alpha1.MobileSecurityService{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: utils.SERVICE_INSTANCE_NAME, Namespace: utils.SERVICE_INSTANCE_NAMESPACE}, serviceInstance); err != nil {
		// Return and don't create
		reqLogger.Info("Mobile Security Service instance resource not found. Ignoring since object must be deleted")
		return reconcile.Result{}, nil
	}

	//Check specs
	if !hasMandatorySpecs(instance, serviceInstance, reqLogger) {
		return reconcile.Result{Requeue: true}, nil
	}

	reqLogger.Info("Checking if the route already exists ...")
	route := &routev1.Route{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: utils.GetRouteName(serviceInstance), Namespace: utils.SERVICE_INSTANCE_NAMESPACE}, route); err != nil {
		return reconcile.Result{}, err
	}

	//Get the REST Service Endpoint
	serviceAPI := service.GetServiceAPIURL(route, serviceInstance)

	//Check if ConfigMap for the app exist, if not create one.
	if _, err := r.fetchSDKConfigMap(reqLogger, instance); err != nil {
		return r.create(instance, CONFIGMAP, serviceAPI, reqLogger, err)
	}

	//Check if App is Bind in the REST Service, if not then bind it
	if app, err := fetchBindAppRestServiceByAppID(serviceAPI, instance, reqLogger); err == nil {
		//Update the name by the REST API
		if hasApp(app) && app.AppName != instance.Spec.AppName {
			app.AppName = instance.Spec.AppName
			//Check if App was update with success

			if err := service.UpdateAppNameByRestAPI(serviceAPI, app, reqLogger); err != nil {
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		}
		// Bind App in the Service by the REST API
		if !hasApp(app) {
			newApp := models.NewApp(instance.Spec.AppName, instance.Spec.AppId)
			if err := service.CreateAppByRestAPI(serviceAPI, newApp, reqLogger); err != nil {
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		}
	}

	//Update status for SDKConfigMap
	SDKConfigMapStatus, err := r.updateSDKConfigMapStatus(reqLogger, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	//Update status for BindStatus
	if err := r.updateBindStatus(serviceAPI, reqLogger, SDKConfigMapStatus, instance); err != nil {
		return reconcile.Result{}, err
	}

	//Handle the finalizer when the CR is deleted to unbind the app
	if err := r.handleFinalizer(serviceAPI, reqLogger, instance); err != nil {
		return reconcile.Result{}, err
	}

	//Update finalizer timestamp to allow delete CR
	if err := r.updateFinilizer(serviceAPI, reqLogger, instance); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
