package service

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"

)

const ENDPOINT_API = "/api"

//Return the router host URL
func GetServiceURL(route *routev1.Route, serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService) string{
	return serviceInstance.Spec.ClusterProtocol + "://" + route.Status.Ingress[0].Host
}

//Return REST Service API
func GetServiceAPIURL(route *routev1.Route, serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService) string{
	return GetServiceURL(route,serviceInstance) + ENDPOINT_API
}

