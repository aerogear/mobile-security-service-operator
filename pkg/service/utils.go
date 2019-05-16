package service

import (
	"fmt"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
)

const ENDPOINT_API = "/api"

//Return REST Service API
func GetServiceAPIURL(mssInstance *mobilesecurityservicev1alpha1.MobileSecurityService) string {
	uri := mssInstance.Spec.ClusterProtocol + "://" + utils.APPLICATION_SERVICE_INSTANCE_NAME + ":" + fmt.Sprint(mssInstance.Spec.Port) + ENDPOINT_API
	return uri
}
