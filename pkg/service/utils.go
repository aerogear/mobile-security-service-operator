package service

import (
	"fmt"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
)

const ENDPOINT_API = "/api"

//Return REST Service API
func GetServiceAPIURL(serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService) string {
	return serviceInstance.Spec.ClusterProtocol + "://" + utils.SERVER_SERVICE_INSTANCE_NAME + ":" + fmt.Sprint(serviceInstance.Spec.Port) + ENDPOINT_API
}
