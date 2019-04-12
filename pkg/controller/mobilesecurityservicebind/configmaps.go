package mobilesecurityservicebind

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Returns the ConfigMap with the properties used to setup/config the Mobile Security Service Project
func (r *ReconcileMobileSecurityServiceBind) buildAppBindSDKConfigMap(m *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, pod corev1.Pod) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.GetAppIdByPodLabel(pod, m) + "-sdk",
			Namespace: pod.Namespace,
			Labels:    getAppLabelsForSDKConfigMap(m.Name),
		},
		Data: getConfigMapSDKForMobileSecurityService(m, pod),
	}
	// Set MobileSecurityService instance as the owner and controller
	controllerutil.SetControllerReference(m, configMap, r.scheme)
	return configMap
}