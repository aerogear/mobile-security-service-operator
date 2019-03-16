package mobilesecurityservice

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"k8s.io/apimachinery/pkg/util/intstr"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)


// getDeploymentForMobileSecurityService returns a MobileSecurityService Deployment resource
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
