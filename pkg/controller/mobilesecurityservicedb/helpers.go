package mobilesecurityservicedb

import (
	"context"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func getDBLabels(name string) map[string]string {
	return map[string]string{"app": "mobilesecurityservicedb", "mobilesecurityservicedb_cr": name, "name": "mobilesecurityservicedb"}
}

func (r *ReconcileMobileSecurityServiceDB) getDatabaseNameEnvVar(m *mobilesecurityservicev1alpha1.MobileSecurityServiceDB) corev1.EnvVar {
	if r.hasAppConfigMap(m) {
		return corev1.EnvVar{
			Name: "POSTGRES_DB",
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: m.Spec.ConfigMapName,
					},
					Key: "PGDATABASE",
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  "POSTGRES_DB",
		Value: m.Spec.DatabaseName,
	}
}

//Check if has App Config Map created
func (r *ReconcileMobileSecurityServiceDB) hasAppConfigMap(m *mobilesecurityservicev1alpha1.MobileSecurityServiceDB) bool {
	configMap := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: "mobile-security-service-app", Namespace: m.Namespace}, configMap)
	if err != nil {
		return false
	}
	return true
}

func (r *ReconcileMobileSecurityServiceDB) getDatabaseUserEnvVar(m *mobilesecurityservicev1alpha1.MobileSecurityServiceDB) corev1.EnvVar {
	if r.hasAppConfigMap(m) {
		return corev1.EnvVar{
			Name: "POSTGRES_USER",
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: m.Spec.ConfigMapName,
					},
					Key: "PGUSER",
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  "POSTGRES_USER",
		Value: m.Spec.DatabaseUser,
	}
}

func (r *ReconcileMobileSecurityServiceDB) getDatabasePasswordEnvVar(m *mobilesecurityservicev1alpha1.MobileSecurityServiceDB) corev1.EnvVar {
	if r.hasAppConfigMap(m) {
		return corev1.EnvVar{
			Name: "POSTGRES_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: m.Spec.ConfigMapName,
					},
					Key: "PGPASSWORD",
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  "POSTGRES_PASSWORD",
		Value: m.Spec.DatabasePassword,
	}
}