package mobilesecurityserviceapp

import (
	"context"
	"github.com/aerogear/mobile-security-service-operator/pkg/models"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
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

const ConfigMap = "ConfigMap"

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
func (r *ReconcileMobileSecurityServiceApp) create(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, kind string, serviceURL string, reqLogger logr.Logger, request reconcile.Request) error {
	obj := r.buildFactory(reqLogger, instance, kind, serviceURL)
	reqLogger.Info("Creating a new ", "kind", kind, "Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
	err := r.client.Create(context.TODO(), obj)
	if err != nil {
		reqLogger.Error(err, "Failed to create new ", "kind", kind, "Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
	}
	reqLogger.Info("Created successfully - return and create", "kind", kind, "Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
	return err
}

//buildFactory will return the resource according to the kind defined
func (r *ReconcileMobileSecurityServiceApp) buildFactory(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, kind string, serviceURL string) runtime.Object {
	reqLogger.Info("Building Object ", "kind", kind, "MobileSecurityServiceApp.Namespace", instance.Namespace, "MobileSecurityServiceApp.Name", instance.Name)
	switch kind {
	case ConfigMap:
		return r.buildAppSDKConfigMap(instance, serviceURL)
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

	//Fetch the MobileSecurityService App instance
	instance := &mobilesecurityservicev1alpha1.MobileSecurityServiceApp{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		instance, err = r.fetchInstance(reqLogger, request)
		if errors.IsNotFound(err) {
			// Return and don't create
			reqLogger.Info("Mobile Security Service App resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - create the request.
		reqLogger.Error(err, "Failed to get Mobile Security Service App")
		return reconcile.Result{}, err
	}

	// FIXME: Check if is a valid namespace
	// We should not checked if the namespace is valid or not. It is an workaround since currently is not possible watch/cache a List of Namespaces
	// The impl to allow do it is done and merged in the master branch of the lib but not released in an stable version. It should be removed when this feature be impl.
	// See the PR which we are working on to update the deps and have this feature: https://github.com/operator-framework/operator-sdk/pull/1388
	if isValidNamespace, err := utils.IsValidAppNamespace(instance.Namespace); err != nil || isValidNamespace == false {
		// Stop reconcile
		envVar, _ := utils.GetAppNamespaces()
		reqLogger.Error(err, "Unable to reconcile Mobile Security Service App", "instance.Namespace", instance.Namespace, "isValidNamespace", isValidNamespace, "EnvVar.APP_NAMESPACES", envVar)

		//Update status with Invalid Namespace
		if err := r.updateBindStatusWithInvalidNamespace(reqLogger, request); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	reqLogger.Info("Valid namespace for MobileSecurityServiceApp", "Namespace", request.Namespace)

	reqLogger.Info("Checking for service instance ...")
	mssInstance := &mobilesecurityservicev1alpha1.MobileSecurityService{}

	operatorNamespace, err := k8sutil.GetOperatorNamespace()

	// Check if it is a local env or an unit test
	if err == k8sutil.ErrNoNamespace {
		operatorNamespace = utils.OperatorNamespaceForLocalEnv
	}

	// GET MSS CR
	r.client.Get(context.TODO(), types.NamespacedName{Name: utils.MobileSecurityServiceCRName, Namespace: operatorNamespace}, mssInstance)

	//Get the REST Service Endpoint
	serviceAPI := service.GetServiceAPIURL(mssInstance)

	//Check if the APP CR was marked to be deleted
	isAppMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isAppMarkedToBeDeleted && len(instance.GetFinalizers()) > 0 {

		// If the Service was deleted and/or marked to be deleted
		if err := r.client.Get(context.TODO(), types.NamespacedName{Name: utils.MobileSecurityServiceCRName, Namespace: operatorNamespace}, mssInstance); err != nil || mssInstance.GetDeletionTimestamp() != nil {
			reqLogger.Info("Mobile Security Service instance resource not found. Mobile Security Service Application is required to create the application")

			//Remove finalizer
			instance.SetFinalizers(nil)

			//Update CR
			err := r.client.Update(context.TODO(), instance)
			if err != nil {
				reqLogger.Error(err, "Failed to update MobileSecurityService App CR with finalizer")
				return reconcile.Result{}, err
			}

			//Stop the reconcile
			return reconcile.Result{}, nil
		}

		//If the CR was marked to be deleted before it finalizes the app need to be deleted from the Service
		//Do request to get the app.ID to delete app
		app, err := fetchBindAppRestServiceByAppID(serviceAPI, instance, reqLogger)
		if err != nil {
			return reconcile.Result{}, err
		}

		// If the request works with success and the app was found then
		// Do request to delete it from the service
		if app.ID != "" {
			if err := service.DeleteAppFromServiceByRestAPI(serviceAPI, app.ID, reqLogger); err != nil {
				reqLogger.Error(err, "Unable to delete app from Service", "App.ID", app.ID)
				return reconcile.Result{}, err
			}
			reqLogger.Info("Successfully delete app ...")
		}

		// Check if the finalizer criteria is met and remove finalizer from the CR
		if err := r.removeFinalizer(serviceAPI, reqLogger, request); err != nil {
			return reconcile.Result{}, err
		}

		//Stop the reconcile
		return reconcile.Result{}, nil
	}

	if !hasMandatorySpecs(instance, mssInstance, reqLogger) {
		//Stop reconcile since it has not the mandatory specs
		return reconcile.Result{}, nil
	}

	// Add finalizer for this CR
	if err := r.addFinalizer(reqLogger, instance, request); err != nil {
		return reconcile.Result{}, err
	}

	// Get the route in order to obtain the public Service URL API
	reqLogger.Info("Checking if the route already exists ...")
	route := &routev1.Route{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: utils.GetRouteName(mssInstance), Namespace: operatorNamespace}, route); err != nil {
		return reconcile.Result{}, err
	}

	// Get the Public Service API URL which will be used to build the SDKConfigMap json
	publicServiceURLAPI := utils.GetPublicServiceAPIURL(route, mssInstance)

	//Check if ConfigMap for the app exist, if not create one.
	if _, err := r.fetchSDKConfigMap(reqLogger, instance); err != nil {
		if err := r.create(instance, ConfigMap, publicServiceURLAPI, reqLogger, request); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Fetch app
	app, err := fetchBindAppRestServiceByAppID(serviceAPI, instance, reqLogger)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Bind App in the Service by the REST API
	// NOTE: If the app was soft deleted before it will make the required job as well
	if app.ID == "" {
		newApp := models.NewApp(instance.Spec.AppName, instance.Spec.AppId)
		if err := service.CreateAppByRestAPI(serviceAPI, newApp, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Update the app name if it was changed.
	if app.AppName != instance.Spec.AppName {

		// Re-fetch the app to get the app.ID since now it was created in the Service
		if app.ID == "" {
			app, err = fetchBindAppRestServiceByAppID(serviceAPI, instance, reqLogger)
			if err != nil {
				return reconcile.Result{}, err
			}
		}

		// Update the name by the REST API when exists the app
		app.AppName = instance.Spec.AppName

		//Check if App was update with success
		if err := service.UpdateAppNameByRestAPI(serviceAPI, app, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
	}

	//Update status for SDKConfigMap
	SDKConfigMapStatus, err := r.updateSDKConfigMapStatus(reqLogger, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	//Update status for BindStatus
	if err := r.updateBindStatus(serviceAPI, reqLogger, SDKConfigMapStatus, request); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
