package mobilesecurityservice

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobile-security-service/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	mssInstance = mobilesecurityservicev1alpha1.MobileSecurityService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.MobileSecurityServiceCRName,
			Namespace: utils.OperatorNamespaceForLocalEnv,
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceSpec{
			Size:            1,
			MemoryLimit:     "512Mi",
			MemoryRequest:   "512Mi",
			ClusterProtocol: "http",
			ConfigMapName:   "mss-config",
			RouteName:       "route",
		},
	}

	mssInstance2 = mobilesecurityservicev1alpha1.MobileSecurityService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.MobileSecurityServiceCRName,
			Namespace: utils.OperatorNamespaceForLocalEnv + "2",
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceSpec{
			Size:            1,
			MemoryLimit:     "512Mi",
			MemoryRequest:   "512Mi",
			ClusterProtocol: "http",
			ConfigMapName:   "mss-config",
			RouteName:       "route",
		},
	}

	configMap = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mssInstance.Spec.ConfigMapName,
			Namespace: mssInstance.Namespace,
			Labels:    getAppLabels(mssInstance.Name),
		},
		Data: getAppEnvVarsMap(&mssInstance),
	}

	mssInstanceWithoutSpec = mobilesecurityservicev1alpha1.MobileSecurityService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.MobileSecurityServiceCRName,
			Namespace: utils.OperatorNamespaceForLocalEnv,
		},
	}
)
