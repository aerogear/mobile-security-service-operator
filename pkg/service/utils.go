package service

import (
	"fmt"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
)

const ApiEndpoint = "/api"

//Return REST Service API
func GetServiceAPIURL(mssInstance *mobilesecurityservicev1alpha1.MobileSecurityService) string {
	return mssInstance.Spec.ClusterProtocol + "://" + utils.ApplicationServiceInstanceName + ":" + fmt.Sprint(mssInstance.Spec.Port) + ApiEndpoint
}
