package mobilesecurityserviceapp

import (
	"context"
	"github.com/aerogear/mobile-security-service/pkg/models"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
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

// Update the factory object and requeue
func (r *ReconcileMobileSecurityServiceApp) delete(obj runtime.Object, reqLogger logr.Logger) error {
	err := r.client.Delete(context.TODO(), obj)
	if err != nil {
		reqLogger.Error(err, "Failed to delete obj", "obj:", obj)
		return err
	}
	reqLogger.Info("Deleted successfully", "obj:", obj)
	return nil
}

// ReconcileMobileSecurityServiceApp reconciles a MobileSecurityServiceApp object
type ReconcileMobileSecurityServiceApp struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Build the object, cluster resource, and add the object in the queue to reconcile
func (r *ReconcileMobileSecurityServiceApp) create(app *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, kind string, serviceURL string, reqLogger logr.Logger, request reconcile.Request) error {
	obj := r.buildFactory(reqLogger, app, kind, serviceURL)
	reqLogger.Info("Creating a new ", "kind", kind, "Instance.Namespace", app.Namespace, "Instance.Name", app.Name)
	err := r.client.Create(context.TODO(), obj)
	if err != nil {
		reqLogger.Error(err, "Failed to create new ", "kind", kind, "Instance.Namespace", app.Namespace, "Instance.Name", app.Name)
	}
	reqLogger.Info("Created successfully", "kind", kind, "Instance.Namespace", app.Namespace, "Instance.Name", app.Name)
	return err
}

// buildFactory will return the resource according to the kind defined
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

	//Fetch the MobileSecurityService App app
	app := &mobilesecurityservicev1alpha1.MobileSecurityServiceApp{}
	app, err := r.fetchAppInstance(reqLogger, request)
	if err != nil {
		reqLogger.Error(err, "Failed to get Mobile Security Service App")
		return reconcile.Result{}, err
	}

	// Ensure that the APP CR will be installed and performed just in the namespaces setup in the ENV VAR APP NAMESPACES
	if isValidNamespace, err := utils.IsValidAppNamespace(app.Namespace); err != nil || isValidNamespace == false {
		// Stop reconcile
		envVar, _ := utils.GetAppNamespaces()
		reqLogger.Error(err, "Unable to reconcile Mobile Security Service App", "app.Namespace", app.Namespace, "isValidNamespace", isValidNamespace, "EnvVar.APP_NAMESPACES", envVar)

		//Update status with Invalid Namespace
		if err := r.updateBindStatusWithInvalidNamespace(reqLogger, request); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	reqLogger.Info("Valid namespace for MobileSecurityServiceApp", "Namespace", request.Namespace)
	reqLogger.Info("Checking for service app ...")
	mssInstance := &mobilesecurityservicev1alpha1.MobileSecurityService{}

	operatorNamespace, err := k8sutil.GetOperatorNamespace()

	// Check if it is a local env or an unit test
	if err == k8sutil.ErrNoNamespace {
		operatorNamespace = utils.OperatorNamespaceForLocalEnv
	}

	// GET MSS CR
	r.client.Get(context.TODO(), types.NamespacedName{Name: utils.MobileSecurityServiceCRName, Namespace: operatorNamespace}, mssInstance)

	// Get the REST Service Endpoint
	serviceAPI := utils.GetServiceAPIURL(mssInstance)

	// Check if has Conditionals to be deleted and perform the actions required to allow it.
	if hasConditionsToBeDeleted(app, mssInstance) {
		// Try to fetch MSS and check if the Service was deleted and/or marked to be deleted
		if r.isMobileSecurityServiceDeleted(operatorNamespace, mssInstance) {
			reqLogger.Info("Mobile Security Service app resource not found. Mobile Security Service Application is required to create the application")

			if err := r.removeFinalizerFromCR(app); err != nil {
				reqLogger.Error(err, "Failed to update MobileSecurityService App CR with finalizer")
				return reconcile.Result{}, err
			}

			//Stop the reconcile
			return reconcile.Result{}, nil
		}

		// If the CR was marked to be deleted before it finalizes the appService need to be deleted from the Service
		// Do request to get the appService.ID to delete appService
		app, err := fetchBindAppRestServiceByAppID(serviceAPI, app, reqLogger)
		if err != nil {
			return reconcile.Result{}, err
		}

		// If the request works with success and the appService was found then
		// Do request to delete it from the service
		if app.ID != "" {
			if err := service.DeleteAppFromServiceByRestAPI(serviceAPI, app.ID, reqLogger); err != nil {
				reqLogger.Error(err, "Unable to delete appService from Service", "App.ID", app.ID)
				return reconcile.Result{}, err
			}
			reqLogger.Info("Successfully delete appService ...")
		}

		// Check if the finalizer criteria is met and remove finalizer from the CR
		if err := r.handleFinalizer(serviceAPI, reqLogger, request); err != nil {
			return reconcile.Result{}, err
		}

		//Stop the reconcile
		return reconcile.Result{}, nil
	}

	if !hasMandatorySpecs(app, reqLogger) {
		//Stop reconcile since it has not the mandatory specs
		return reconcile.Result{}, nil
	}

	// Add finalizer for this CR
	if err := r.addFinalizer(reqLogger, app, request); err != nil {
		return reconcile.Result{}, err
	}

	// Get the route in order to obtain the public Service URL API
	reqLogger.Info("Checking if the route already exists ...")
	route := &routev1.Route{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: mssInstance.Spec.RouteName, Namespace: operatorNamespace}, route); err != nil {
		return reconcile.Result{}, err
	}

	// Get the Public Service API URL which will be used to build the SDKConfigMap json
	publicServiceURLAPI := utils.GetInitPublicURL(route, mssInstance)

	reqLogger.Info("Checking if the configMap already exists ...")
	// Check if ConfigMap for the appService exist, if not create one.
	if _, err := r.fetchConfigMap(reqLogger, app); err != nil {
		if err := r.create(app, ConfigMap, publicServiceURLAPI, reqLogger, request); err != nil {
			return reconcile.Result{}, err
		}
	}

	reqLogger.Info("Checking if a configMap for the same appId already exists with another name ...")
	// Check if has already a confgMap for the same appId, if yes and the name is not the same then remove it
	if list, err := r.fetchConfigMapListByLabels(reqLogger, app); err == nil && len(list.Items) > 0 {
		for i := 0; i < len(list.Items); i++ {
			cmItem := list.Items[i]

			// If found a ConfigMap with the same appID but different name then remove.
			if cmItem.Name != getSDKConfigMapName(app) {
				if err := r.delete(&cmItem, reqLogger); err != nil {
					return reconcile.Result{}, err
				}
			}
		}
	}

	reqLogger.Info("Checking if has more than one configMap for the same appId ...")
	//Ensure that it has always just one configMap for each appId
	if list, err := r.fetchConfigMapListByLabels(reqLogger, app); err == nil && len(list.Items) > 1 {
		// Remove all and leave just one
		for i := 0; i < len(list.Items)-1; i++ {
			cmItem := list.Items[i]
			if err := r.delete(&cmItem, reqLogger); err != nil {
				return reconcile.Result{}, err
			}
		}
		//Requeu in order to do a full validation
		return reconcile.Result{Requeue: true}, nil
	}

	// Fetch appService
	appService, err := fetchBindAppRestServiceByAppID(serviceAPI, app, reqLogger)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Bind App in the Service by the REST API
	// NOTE: If the appService was soft deleted before it will make the required job as well
	if appService.ID == "" {
		newApp := models.NewAppByNameAndAppID(app.Spec.AppName, app.Spec.AppId)
		if err := service.CreateAppByRestAPI(serviceAPI, newApp, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Update the appService name if it was changed.
	if appService.AppName != app.Spec.AppName {

		// Re-fetch the appService to get the appService.ID since now it was created in the Service
		if appService.ID == "" {
			appService, err = fetchBindAppRestServiceByAppID(serviceAPI, app, reqLogger)
			if err != nil {
				return reconcile.Result{}, err
			}
		}

		// Update the name by the REST API when exists the appService
		appService.AppName = app.Spec.AppName

		//Check if App was update with success
		if err := service.UpdateAppNameByRestAPI(serviceAPI, appService, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
	}

	//Update status for SDKConfigMap
	SDKConfigMapStatus, err := r.updateConfigMapStatus(reqLogger, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	//Update status for BindStatus
	if err := r.updateBindStatus(serviceAPI, reqLogger, SDKConfigMapStatus, request); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
