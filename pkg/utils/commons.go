package utils

import (
	"fmt"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"os"
)

// APP_NAMESPACE_ENV_VAR is the constant for env variable APP_NAMESPACES
// which is the namespace where the APP CR can applied.
// The namespaces should be informed split by ";".
const APP_NAMESPACE_ENV_VAR = "APP_NAMESPACES"
const SERVICE_INSTANCE_NAME  = "mobile-security-service"
const SERVICE_INSTANCE_NAMESPACE  = "mobile-security-service-operator"


//GetRouteName returns an string name with the name of the router
func GetRouteName(m *mobilesecurityservicev1alpha1.MobileSecurityService) string{
	if len(m.Spec.RouteName) > 0 {
		return m.Spec.RouteName
	}
	return m.Name
}

//GetConfigMapName returns an string name with the name of the configMap
func GetConfigMapName(m *mobilesecurityservicev1alpha1.MobileSecurityService) string {
	if len(m.Spec.ConfigMapName) > 0 {
		return m.Spec.ConfigMapName
	}
	return m.Name
}

// GetCustomWatchNamespaces returns the namespace the operator should be watching for changes
func GetAppNamespaces() (string, error) {
	ns, found := os.LookupEnv(APP_NAMESPACE_ENV_VAR)
	if !found {
		return "", fmt.Errorf("%s must be set", APP_NAMESPACE_ENV_VAR)
	}
	return ns, nil
}