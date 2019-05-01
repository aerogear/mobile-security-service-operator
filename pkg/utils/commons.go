package utils

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
)

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