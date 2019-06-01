package mobilesecurityserviceapp

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Returns the ConfigMap with the properties used to setup/config the Mobile Security Service Project
func (r *ReconcileMobileSecurityServiceApp) buildAppSDKConfigMap(mssApp *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, serviceURL string) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getSDKConfigMapName(mssApp),
			Namespace: mssApp.Namespace,
			Labels:    getSDKAppLabels(mssApp),
		},
		Data: getConfigMapSDKForMobileSecurityService(serviceURL),
	}
	// Set MobileSecurityService mssApp as the owner and controller
	controllerutil.SetControllerReference(mssApp, configMap, r.scheme)
	return configMap
}
