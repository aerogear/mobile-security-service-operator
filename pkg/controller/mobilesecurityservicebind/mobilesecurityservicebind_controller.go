package mobilesecurityservicebind

import (
	"context"
	"encoding/json"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service/pkg/models"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
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

	//ConfigMap
	if err := watchPod(c); err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMobileSecurityServiceBind{}

// ReconcileMobileSecurityServiceBind reconciles a MobileSecurityServiceBind object
type ReconcileMobileSecurityServiceBind struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

//Update the object and reconcile it
func update(r *ReconcileMobileSecurityServiceBind, obj runtime.Object, reqLogger logr.Logger) (reconcile.Result, error) {
	err := r.client.Update(context.TODO(), obj)
	if err != nil {
		reqLogger.Error(err, "Failed to update Spec")
		return reconcile.Result{}, err
	}
	reqLogger.Info("Spec updated - return and create")
	return reconcile.Result{Requeue: true}, nil
}

func create(r *ReconcileMobileSecurityServiceBind, pod corev1.Pod, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, reqLogger logr.Logger, kind string, err error) (reconcile.Result, error) {
	obj, errBuildObject := buildObject(reqLogger, pod, instance, r, kind)
	if errBuildObject != nil {
		return reconcile.Result{}, errBuildObject
	}
	if errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ", "kind", kind, "Namespace", pod.Namespace)
		err = r.client.Create(context.TODO(), obj)
		if err != nil {
			reqLogger.Error(err, "Failed to create new ", "kind", kind, "Namespace", pod.Namespace)
			return reconcile.Result{}, err
		}
		reqLogger.Info("Created successfully - return and create", "kind", kind, "Namespace", pod.Namespace)
		return reconcile.Result{Requeue: true}, nil
	}
	reqLogger.Error(err, "Failed to get", "kind", kind, "Namespace", pod.Namespace)
	return reconcile.Result{}, err
}

func buildObject(reqLogger logr.Logger, pod corev1.Pod, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, r *ReconcileMobileSecurityServiceBind, kind string) (runtime.Object, error) {
	reqLogger.Info("Check "+kind, "into the namespace", pod.Namespace)
	switch kind {
	case SDK_CONFIGMAP:
		return r.buildAppBindSDKConfigMap(instance, pod), nil
	default:
		msg := "Failed to recognize type of object" + kind + " into the Namespace " + pod.Namespace
		panic(msg)
	}
}

// Reconcile reads that state of the cluster for a MobileSecurityServiceBind object and makes changes based on the state read
// and what is in the MobileSecurityServiceBind.Spec
// Note:
// The Controller will create the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMobileSecurityServiceBind) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MobileSecurityServiceBind")


	// Fetch the MobileSecurityServiceBind instance
	instance := &mobilesecurityservicev1alpha1.MobileSecurityServiceBind{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}


	//Check the key:labels and/or namespace which should be watched
	listOps := getWatchListOps(instance,reqLogger)

	//Check all Deployments in the Namespace where the Bind was installed and with the labelSelector and ValueSelector Specified
	reqLogger.Info("Watching pods by specs ...")
	podList := &corev1.PodList{}
	err = r.client.List(context.TODO(), &listOps, podList)
	if err == nil {
		reqLogger.Info("Listing all pods by specs ...")
		if len(podList.Items) > 0 {
			reqLogger.Info("Found pods by specs ...")
		}
		for i := 0; i < len(podList.Items); i++ {
			//TODO: In progress:
			pod := podList.Items[i]
			reqLogger.Info("*** PodName"+pod.Name)
			log.WithValues("Request.Namespace", request.Namespace, "Pod:", pod)

			//Check if the SDK ConfigMap already exists, if not create a new one
			configmapsdk := &corev1.ConfigMap{}
			configmapsdk_name := pod.Name + "-sdk"
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: configmapsdk_name, Namespace: pod.Namespace}, configmapsdk)
			if err != nil {
				return create(r, pod, instance, reqLogger, SDK_CONFIGMAP, err)
			}

			//Check if the APP was created
			appId := string(pod.UID)
			app, _ := json.Marshal(models.App{AppID:appId, AppName:pod.Name})
			log.WithValues("Request.Namespace", request.Namespace, "** BODY:", strings.NewReader(string(app)))
			req, err := http.NewRequest(http.MethodPost, "http://mobile-security-service-app.192.168.64.16.nip.io/api/apps", strings.NewReader(string(app)))

			if err != nil{
				log.WithValues("Request.Namespace", request.Namespace, "Error to POST app from Service REST API:", err)
			} else {
				log.WithValues("Request.Namespace", request.Namespace, "Response from POST app from Service REST API:", req)
			}

			//Check if the APP changed the name
		}
	} else {
		log.WithValues("Request.Namespace", request.Namespace, "*** Error to List pods by spec:", err)
	}

	return reconcile.Result{}, nil
}

