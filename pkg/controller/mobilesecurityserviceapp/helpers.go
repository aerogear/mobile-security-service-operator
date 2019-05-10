package mobilesecurityserviceapp

import (
	"encoding/json"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/models"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	"github.com/go-logr/logr"
)

const SDK  = "-sdk"
const FINALIZER = "finalizer.mobile-security-service.aerogear.com"

// Returns an string map with the labels which wil be associated to the kubernetes/openshift objects
// which will be created and managed by this operator
func getAppLabels(name string) map[string]string {
	return map[string]string{"app": "mobilesecurityservice", "mobilesecurityserviceapp_cr": name}
}

//To transform the object into a string with its json
func getSdkConfigStringJsonFormat(sdk *models.SDKConfig) string{
	jsonSdk, _ := json.MarshalIndent(sdk, "", "\t")
	return string(jsonSdk)
}

// return properties for the response SDK
func getConfigMapSDKForMobileSecurityService(m *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, serviceURL string) map[string]string {
	sdk := models.NewSDKConfig(m, serviceURL)
	return map[string]string{
		"SDKConfig": getSdkConfigStringJsonFormat(sdk),
	}
}

// return properties for the response SDK
func getSDKConfigMapName(m *mobilesecurityservicev1alpha1.MobileSecurityServiceApp) string {
	return m.Spec.AppName + SDK
}

//Check if the mandatory specs are filled
func hasMandatorySpecs(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService, reqLogger logr.Logger) bool {
	//Check if the appId was added in the CR
	if len(instance.Spec.AppId) < 1 {
		reqLogger.Info("AppID was not found. Check the App CR configuration.")
		return false
	}

	//Check if the appId was added in the CR
	if len(instance.Spec.AppId) < 1 {
		reqLogger.Info("AppName was not found. Check the App CR configuration.")
		return false
	}

	//Check the values defined for the ClusterProtocol in the MobileSecurityService CR
	if res := utils.CheckClusterProtocol(serviceInstance, reqLogger); res != true {
		return false
	}

	return true
}