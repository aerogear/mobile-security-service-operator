package mobilesecurityservice

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/go-logr/logr"
)

//Check if the mandatory specs are filled
func hasMandatorySpecs(serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService, reqLogger logr.Logger) bool {
	//Check the values defined for the ClusterProtocol in the MobileSecurityService CR
	if res := checkClusterProtocol(serviceInstance, reqLogger); res != true {
		return false
	}

	return true
}

// CheckClusterProtocol will return true if the cluster protocol was added correctly.
func checkClusterProtocol(serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService, reqLogger logr.Logger) bool {
	//Check if the cluster protocol was defined
	if len(serviceInstance.Spec.ClusterProtocol) < 1 {
		reqLogger.Info("Unable to get the config of the ClusterProtocol setup for the service. Check its property in the MobileSecurityService CR")
		return false
	}

	//Check if the cluster protocol was defined properly
	if serviceInstance.Spec.ClusterProtocol != "http" && serviceInstance.Spec.ClusterProtocol != "https" {
		reqLogger.Info("Invalid config for ClusterProtocol setup in the service %s. Check its property in the MobileSecurityService CR", serviceInstance.Spec.ClusterProtocol)
		return false
	}
	return true
}
