package mobilesecurityservicedb

import (
	"context"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func getDBLabels(name string) map[string]string {
	return map[string]string{"app": "mobilesecurityservicedb", "mobilesecurityservicedb_cr": name, "name": "mobilesecurityservicedb"}
}

func (r *ReconcileMobileSecurityServiceDB) getDatabaseNameEnvVar(m *mobilesecurityservicev1alpha1.MobileSecurityServiceDB, serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService) corev1.EnvVar {
	if r.hasAppConfigMap(m, serviceInstance) {
		return corev1.EnvVar{
			Name: m.Spec.DatabaseNameParam,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: utils.GetConfigMapName(serviceInstance),
					},
					Key: "PGDATABASE",
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  m.Spec.DatabaseNameParam,
		Value: m.Spec.DatabaseName,
	}
}

//Check if has App Config Map created
func (r *ReconcileMobileSecurityServiceDB) hasAppConfigMap(m *mobilesecurityservicev1alpha1.MobileSecurityServiceDB, serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService) bool {
	//if has not service instance return false
	if len(serviceInstance.Name) < 1 {
		return false
	}

	//Looking for the configMap created by the service instance
	configMap := &corev1.ConfigMap{}
	operatorNamespace, _ := k8sutil.GetOperatorNamespace()
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: utils.GetConfigMapName(serviceInstance), Namespace: operatorNamespace}, configMap)
	if err != nil {
		return false
	}
	return true
}

func (r *ReconcileMobileSecurityServiceDB) getDatabaseUserEnvVar(m *mobilesecurityservicev1alpha1.MobileSecurityServiceDB, serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService) corev1.EnvVar {
	if r.hasAppConfigMap(m, serviceInstance) {
		return corev1.EnvVar{
			Name: m.Spec.DatabaseUserParam,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: utils.GetConfigMapName(serviceInstance),
					},
					Key: "PGUSER",
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  m.Spec.DatabaseUserParam,
		Value: m.Spec.DatabaseUser,
	}
}

func (r *ReconcileMobileSecurityServiceDB) getDatabasePasswordEnvVar(m *mobilesecurityservicev1alpha1.MobileSecurityServiceDB, serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService) corev1.EnvVar {
	if r.hasAppConfigMap(m, serviceInstance) {
		return corev1.EnvVar{
			Name: m.Spec.DatabasePasswordParam,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: utils.GetConfigMapName(serviceInstance),
					},
					Key: "PGPASSWORD",
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  m.Spec.DatabasePasswordParam,
		Value: m.Spec.DatabasePassword,
	}
}
