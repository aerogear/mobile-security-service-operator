package mobilesecurityservicebind

import (
	"context"
	"fmt"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
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
	SDK_CONFIGMAP = "SDKConfigMap"
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

	//Pod
	if err := watchPod(c); err != nil {
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
func create(r *ReconcileMobileSecurityServiceBind, pod corev1.Pod, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, kind string, reqLogger logr.Logger, err error) (reconcile.Result, error) {
	obj, errBuildObject := buildObject(reqLogger, pod, instance, r, kind)
	if errBuildObject != nil {
		return reconcile.Result{}, errBuildObject
	}
	if errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ", "kind", kind, "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), obj)
		if err != nil {
			reqLogger.Error(err, "Failed to create new ", "kind", kind, "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
			return reconcile.Result{}, err
		}
		reqLogger.Info("Created successfully - return and create", "kind", kind, "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		return reconcile.Result{Requeue: true}, nil
	}
	reqLogger.Error(err, "Failed to build", "kind", kind, "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
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
func buildObject(reqLogger logr.Logger, pod corev1.Pod, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, r *ReconcileMobileSecurityServiceBind, kind string) (runtime.Object, error) {
	reqLogger.Info("Building Object ", "kind", kind, "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	switch kind {
	case SDK_CONFIGMAP:
		return r.buildAppBindSDKConfigMap(instance, pod), nil
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
		reqLogger.Error( err,"Please check its configuration. See https://github.com/aerogear/mobile-security-service-operator#configuring .")
		return reconcile.Result{}, err
	}

	//Check the key:labels and/or namespace which should be watched to get the bind apps
	listAppOps := getAppWatchListOps(instance,reqLogger)

	//Check all pods which are bind to the service
	reqLogger.Info("Checking pods application according to the specs ...")
	appPodList := &corev1.PodList{}
	err = r.client.List(context.TODO(), &listAppOps, appPodList)
	if err == nil {
		if len(appPodList.Items) > 0 {
			reqLogger.Info("Found pods by specs defined ...")
		}

		for i := 0; i < len(appPodList.Items); i++ {
			// Get required data
			pod := appPodList.Items[i]
			appName := utils.GetAppNameByPodLabel(pod, instance)
			isBind := isBind(pod, instance)
			reqLogger.Info("Reconcile Mobile Security Service for the data:", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name, "Pod.isBind", isBind, "Pod.App.Name", appName, "Pod.Labels", pod.Labels)

			// Reconcile SDKConfigMap
			reqLogger.Info("Checking if has SDK ConfigMap already exists for the pod ...")
			configmapsdk := &corev1.ConfigMap{}
			configMapName := utils.GetAppIdByPodLabel(pod, instance) + "-sdk"
			reqLogger.Info("Search for the SDKConfigMap in:", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name, "ConfigMap.Name", configMapName)
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: configMapName, Namespace: pod.Namespace}, configmapsdk)

			reqLogger.Info("Checking if the app already exists in the REST Service API ...")
			foundApp, _ := getAppFromServiceByRestApi(instance, pod, reqLogger)

			//If the app/pod is bind to the service
			if isBind {
				//Create SDKConfig because it was not found
				if err != nil {
					reqLogger.Info("Creating the SDKConfigMap ...")
					return create(r, pod, instance, SDK_CONFIGMAP, reqLogger, err)
				}

				if len(foundApp.ID) > 0 {
					if foundApp.AppName != appName {
						//Remove the config-sdk with the old name
						outdatedConfigMapName := foundApp.AppName + "-sdk"
						outdatedConfigMap := &corev1.ConfigMap{}
						reqLogger.Info("Search for the outdated SDKConfigMap in:", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name, "ConfigMap.Name", outdatedConfigMapName)
						errOutdatedConfigMap := r.client.Get(context.TODO(), types.NamespacedName{Name: outdatedConfigMapName, Namespace: pod.Namespace}, outdatedConfigMap)
						if errOutdatedConfigMap == nil {
							reqLogger.Info("Deleting the outdated SDKConfigMap ...")
							return delete(r, configmapsdk, reqLogger)
						}

						//Update the name by the REST API
						reqLogger.Info("Updating the app name ...", "NewName:", appName, "OldName", foundApp.AppName)
						return updateAppNameByRestAPI(instance, foundApp, pod, reqLogger)
					}
				} else {
					reqLogger.Info("Calling the Rest API for create app...")
					return createAppByRestAPI(instance, pod, reqLogger)
				}
				reqLogger.Info("Unable to reconcile the bind operations with the REST API")
				return reconcile.Result{Requeue: true}, nil
			}

			//If the app/pod is NOT bind to the service
			if !isBind {
				//Delete SDKConfigMap because it was found and the app/pod is no long bind to the service
				if err == nil {
					reqLogger.Info("Deleting the SDKConfigMap ...")
					return delete(r, configmapsdk, reqLogger)
				}

				//Delete the app from the Service by the REST API since it is no longer bind to the service
				if len(foundApp.ID) > 0 {
					reqLogger.Info("Calling the Rest API for deleting app...")
					return deleteAppFromServiceByRestAPI(instance,foundApp, reqLogger)
				}
				reqLogger.Info("Unable to reconcile the bind operations with the REST API")
				return reconcile.Result{Requeue: true}, nil
			}

		}
	}

	// The following code will store all appName label of the applications/pods which are bind to the service
	// The reconcile will be recalled when have any change on it.
	reqLogger.Info("Updating the MobileSecurityServiceBind status with the pod names")
	podList := &corev1.PodList{}
	err = r.client.List(context.TODO(), &listAppOps, podList)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods", "MobileSecurityServiceBind.Namespace", instance.Namespace, "MobileSecurityServiceBind.Name", instance.Name)
		return reconcile.Result{}, err
	}
	reqLogger.Info("Get pod names")
	podNames := utils.GetPodNames(podList.Items)
	reqLogger.Info("Update status.Nodes if needed")
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update MobileSecurityServiceBind status")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

