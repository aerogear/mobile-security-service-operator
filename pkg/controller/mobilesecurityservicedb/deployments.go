package mobilesecurityservicedb

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//Returns the Deployment object for the Mobile Security Service Database
func (r *ReconcileMobileSecurityServiceDB) buildDBDeployment(m *mobilesecurityservicev1alpha1.MobileSecurityServiceDB) *appsv1.Deployment {
	ls := getDBLabels(m.Name)
	auto := true
	replicas := m.Spec.Size
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
		Spec: appsv1.DeploymentSpec{
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
						Image:           m.Spec.Image,
						Name:            m.Spec.ContainerName,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{{
							ContainerPort: m.Spec.DatabasePort,
							Protocol:      "TCP",
						}},
						Env: []corev1.EnvVar{
							r.getDatabaseNameEnvVar(m),
							r.getDatabaseUserEnvVar(m),
							r.getDatabasePasswordEnvVar(m),
							{
								Name:  "PGDATA",
								Value: "/var/lib/pgsql/data/pgdata",
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      m.Name,
								MountPath: "/var/lib/pgsql/data",
							},
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								Exec: &corev1.ExecAction{
									Command: []string{
										"pg_isready",
										"-h",
										"localhost",
										"-U",
										m.Spec.DatabaseUser,
									},
								},
							},
							InitialDelaySeconds: 5,
							TimeoutSeconds:      1,
						},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								Exec: &corev1.ExecAction{
									Command: []string{
										"pg_isready",
										"-h",
										"localhost",
										"-U",
										m.Spec.DatabaseUser,
									},
								},
							},
							InitialDelaySeconds: 120,
							TimeoutSeconds:      10,
						},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(m.Spec.DatabaseMemoryLimit),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(m.Spec.DatabaseMemoryRequest),
							},
						},
						TerminationMessagePath: "/dev/termination-log",
					}},
					DNSPolicy:     corev1.DNSClusterFirst,
					RestartPolicy: corev1.RestartPolicyAlways,
					Volumes: []corev1.Volume{
						{
							Name: m.Name,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: m.Name,
								},
							},
						},
					},
					AutomountServiceAccountToken: &auto,
				},
			},
		},
	}
	// Set MobileSecurityService instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}
