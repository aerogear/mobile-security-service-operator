package mobilesecurityservice

import (
	"context"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_mobilesecurityservice")


/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new MobileSecurityService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMobileSecurityService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mobilesecurityservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MobileSecurityService
	err = c.Watch(&source.Kind{Type: &mobilesecurityservicev1alpha1.MobileSecurityService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	/** Watch for changes to secondary resources and requeue the owner MobileSecurityService **/

	//ConfigMap
	if err:= watchConfigMap(c); err != nil {
		return err
	}

	//Deployment
	if err:= watchDeployment(c); err != nil {
		return err
	}

	//Service
	if err:= watchService(c); err != nil {
		return err
	}

	return nil
}


var _ reconcile.Reconciler = &ReconcileMobileSecurityService{}

// ReconcileMobileSecurityService reconciles a MobileSecurityService object
type ReconcileMobileSecurityService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a MobileSecurityService object and makes changes based on the state read
// and what is in the MobileSecurityService.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMobileSecurityService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MobileSecurityService")

	//Fetch the MobileSecurityService instance
	instance := &mobilesecurityservicev1alpha1.MobileSecurityService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("MobileSecurityService resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get MobileSecurityService")
		return reconcile.Result{}, err
	}

	//Check if the ConfigMap already exists, if not create a new one
	configmap := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, configmap)
	if err != nil && errors.IsNotFound(err) {
		cfg := r.getConfigMapForMobileSecurityService(instance)
		reqLogger.Info("Creating a new ConfigMap", "ConfigMap.Namespace", cfg.Namespace, "ConfigMap.Name", cfg.Name)
		err = r.client.Create(context.TODO(), cfg)
		if err != nil {
			reqLogger.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", cfg.Namespace, "ConfigMap.Name", cfg.Name)
			return reconcile.Result{}, err
		}
		reqLogger.Info("ConfigMap created successfully - return and requeue")
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get ConfigMap")
		return reconcile.Result{}, err
	}

	//Check if the deployment already exists, if not create a new one
	deployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, deployment)
	if err != nil && errors.IsNotFound(err) {
		dep := r.getDeploymentForMobileSecurityService(instance)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		reqLogger.Info("Deployment created successfully - return and requeue")
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}

	//Ensure the deployment size is the same as the spec
	reqLogger.Info("Ensure the deployment size is the same as the spec")
	size := instance.Spec.Size
	if *deployment.Spec.Replicas != size {
		deployment.Spec.Replicas = &size
		err = r.client.Update(context.TODO(), deployment)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return reconcile.Result{}, err
		}
		reqLogger.Info("Spec updated - return and requeue")
		return reconcile.Result{Requeue: true}, nil
	}

	//Check if the Service already exists, if not create a new one
	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Define a new service")
		ser := r.getServiceForMobileSecurityService(instance)
		reqLogger.Info("Creating a new Service", "Service.Namespace", ser.Namespace, "Service.Name", ser.Name)
		err = r.client.Create(context.TODO(), ser)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", ser.Namespace, "Service.Name", ser.Name)
			return reconcile.Result{}, err
		}
		reqLogger.Info("Service created successfully - return and requeue")
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service")
		return reconcile.Result{}, err
	}

	//Update the MobileSecurityService status with the pod names
	//List the pods for this MobileSecurityService's deployment
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(labelsForMobileSecurityService(instance.Name))
	listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOps, podList)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods", "MobileSecurityService.Namespace", instance.Namespace, "MobileSecurityService.Name", instance.Name)
		return reconcile.Result{}, err
	}
	reqLogger.Info("Get pod names")
	podNames := getPodNames(podList.Items)
	reqLogger.Info("Update status.Nodes if needed")
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update MobileSecurityService status")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// getServiceForMobileSecurityService returns a MobileSecurityService Service object
func (r *ReconcileMobileSecurityService) getConfigMapForMobileSecurityService(m *mobilesecurityservicev1alpha1.MobileSecurityService) *corev1.ConfigMap {
	ls := labelsForMobileSecurityService(m.Name)
	ser := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Data: getConfigMapForMobileSecurityService(),
	}
	// Set MobileSecurityService instance as the owner and controller
	controllerutil.SetControllerReference(m, ser, r.scheme)
	return ser
}

// getDeploymentForMobileSecurityService returns a MobileSecurityService Deployment object
func (r *ReconcileMobileSecurityService) getDeploymentForMobileSecurityService(m *mobilesecurityservicev1alpha1.MobileSecurityService) *appsv1.Deployment {
	ls := labelsForMobileSecurityService(m.Name)
	replicas := m.Spec.Size
	envinronment := getAllEnvVarsToSetupMobileSecurityService(m)
	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec:appsv1.DeploymentSpec{
			Replicas: &replicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},

				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   m.Spec.Image,
						Name:    "security-service",
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{{
							ContainerPort: 3000,
							Name:          "http",
							Protocol:      "TCP",
						}},
						// Get the value from the ConfigMap

						Env: *envinronment,
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/api/healthz",
									Port: intstr.IntOrString{
										Type:   intstr.Int,
										IntVal: int32(3000),
									},
									Scheme: corev1.URISchemeHTTP,
								},
							},
							InitialDelaySeconds: 10,
							FailureThreshold:    3,
							TimeoutSeconds:      10,
							PeriodSeconds:       10,
							SuccessThreshold:    1,
						},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/api/ping",
									Port: intstr.IntOrString{
										Type:   intstr.Int,
										IntVal: int32(3000),
									},
									Scheme: corev1.URISchemeHTTP,
								},
							},
							InitialDelaySeconds: 10,
							FailureThreshold:    3,
							TimeoutSeconds:      10,
							PeriodSeconds:       10,
							SuccessThreshold:    1,

						},
					}},

				},

			},
		},
	}
	// Set MobileSecurityService instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}


// getServiceForMobileSecurityService returns a MobileSecurityService Service object
func (r *ReconcileMobileSecurityService) getServiceForMobileSecurityService(m *mobilesecurityservicev1alpha1.MobileSecurityService) *corev1.Service {
	ls := labelsForMobileSecurityService(m.Name)
	ser := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Type:corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					TargetPort:  intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(3000),
					},
					Port:       3000,
					Protocol:   "TCP",
				},
			},
		},

	}
	// Set MobileSecurityService instance as the owner and controller
	controllerutil.SetControllerReference(m, ser, r.scheme)
	return ser
}

