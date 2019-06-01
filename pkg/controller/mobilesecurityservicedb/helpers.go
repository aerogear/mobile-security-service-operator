package mobilesecurityservicedb

import (
	"context"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	"time"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func getDBLabels(name string) map[string]string {
	return map[string]string{"app": "mobilesecurityservicedb", "mobilesecurityservicedb_cr": name, "name": "mobilesecurityservicedb"}
}

func (r *ReconcileMobileSecurityServiceDB) getDatabaseNameEnvVar(db *mobilesecurityservicev1alpha1.MobileSecurityServiceDB, serviceConfigMapName string) corev1.EnvVar {
	if len(serviceConfigMapName) > 0 {
		return corev1.EnvVar{
			Name: db.Spec.DatabaseNameParam,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: serviceConfigMapName,
					},
					Key: "PGDATABASE",
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  db.Spec.DatabaseNameParam,
		Value: db.Spec.DatabaseName,
	}
}

func (r *ReconcileMobileSecurityServiceDB) getAppConfigMapName(db *mobilesecurityservicev1alpha1.MobileSecurityServiceDB) string {

	serviceConfigMapName := r.fetchServiceConfigMap(db)
	if len(serviceConfigMapName) < 1 {
		// Wait for 30 seconds to check if will be created
		time.Sleep(30 * time.Second)
		// Try again
		serviceConfigMapName = r.fetchServiceConfigMap(db)
	}
	return serviceConfigMapName
}

//Check if has App Config Map created
func (r *ReconcileMobileSecurityServiceDB) fetchServiceConfigMap(db *mobilesecurityservicev1alpha1.MobileSecurityServiceDB) string {
	// It will fetch the service db for the DB type be able to get the configMap config created by it, however,
	// if the Instance cannot be found and/or its configMap was not created than the default values specified in its CR will be used
	serviceInstance := &mobilesecurityservicev1alpha1.MobileSecurityService{}
	r.client.Get(context.TODO(), types.NamespacedName{Name: utils.MobileSecurityServiceCRName, Namespace: db.Namespace}, serviceInstance)

	//if has not service db return false
	if len(serviceInstance.Spec.ConfigMapName) > 1 {
		//Looking for the configMap created by the service db
		configMap := &corev1.ConfigMap{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: serviceInstance.Spec.ConfigMapName, Namespace: db.Namespace}, configMap)
		if err == nil {
			return configMap.Name
		}

	}
	return ""
}

func (r *ReconcileMobileSecurityServiceDB) getDatabaseUserEnvVar(db *mobilesecurityservicev1alpha1.MobileSecurityServiceDB, serviceConfigMapName string) corev1.EnvVar {
	if len(serviceConfigMapName) > 0 {
		return corev1.EnvVar{
			Name: db.Spec.DatabaseUserParam,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: serviceConfigMapName,
					},
					Key: "PGUSER",
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  db.Spec.DatabaseUserParam,
		Value: db.Spec.DatabaseUser,
	}
}

func (r *ReconcileMobileSecurityServiceDB) getDatabasePasswordEnvVar(db *mobilesecurityservicev1alpha1.MobileSecurityServiceDB, serviceConfigMapName string) corev1.EnvVar {
	if len(serviceConfigMapName) > 0 {
		return corev1.EnvVar{
			Name: db.Spec.DatabasePasswordParam,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: serviceConfigMapName,
					},
					Key: "PGPASSWORD",
				},
			},
		}
	}

	return corev1.EnvVar{
		Name:  db.Spec.DatabasePasswordParam,
		Value: db.Spec.DatabasePassword,
	}
}
