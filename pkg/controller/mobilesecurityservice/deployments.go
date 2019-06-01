package mobilesecurityservice

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//buildDeployment returns the Deployment object using as image the MobileSecurityService App ( UI + REST API)
func (r *ReconcileMobileSecurityService) buildDeployment(service *mobilesecurityservicev1alpha1.MobileSecurityService) *v1beta1.Deployment {

	ls := getAppLabels(service.Name)
	replicas := service.Spec.Size
	dep := &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
			Labels:    ls,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: &replicas,
			Strategy: v1beta1.DeploymentStrategy{
				Type: v1beta1.RecreateDeploymentStrategyType,
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: service.Name,
					Containers:         getDeploymentContainers(service),
				},
			},
		},
	}

	// Set MobileSecurityService service as the owner and controller
	controllerutil.SetControllerReference(service, dep, r.scheme)
	return dep
}

func getDeploymentContainers(service *mobilesecurityservicev1alpha1.MobileSecurityService) []corev1.Container {
	var containers []corev1.Container
	containers = append(containers, buildOAuthContainer(service))
	containers = append(containers, buildApplicationContainer(service))
	return containers
}

func buildOAuthContainer(service *mobilesecurityservicev1alpha1.MobileSecurityService) corev1.Container {
	return corev1.Container{
		Image:           service.Spec.OAuthImage,
		Name:            service.Spec.OAuthContainerName,
		ImagePullPolicy: service.Spec.OAuthContainerImagePullPolicy,
		Ports: []corev1.ContainerPort{{
			ContainerPort: service.Spec.OAuthPort,
			Name:          "public",
			Protocol:      "TCP",
		}},
		Args:                     getOAuthArgsMap(service),
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
	}
}

func buildApplicationContainer(service *mobilesecurityservicev1alpha1.MobileSecurityService) corev1.Container {
	environment := buildAppEnvVars(service)
	return corev1.Container{
		Image:           service.Spec.Image,
		Name:            service.Spec.ContainerName,
		ImagePullPolicy: service.Spec.ContainerImagePullPolicy,
		Ports: []corev1.ContainerPort{{
			ContainerPort: service.Spec.Port,
			Name:          "http",
			Protocol:      "TCP",
		}},
		// Get the value from the ConfigMap
		Env: *environment,
		ReadinessProbe: &corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/api/healthz",
					Port: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: service.Spec.Port,
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
						IntVal: service.Spec.Port,
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
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(service.Spec.MemoryLimit),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(service.Spec.MemoryRequest),
			},
		},
	}
}
