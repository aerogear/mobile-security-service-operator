package mobilesecurityserviceapp

import (
	"encoding/json"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/models"
)

const SDK = "-security"
const FINALIZER = "finalizer.mobile-security-service.aerogear.com"

// Returns an string map with the labels which wil be associated to the kubernetes/openshift objects
// which will be created and managed by this operator
func getAppLabels(name string) map[string]string {
	return map[string]string{"app": "mobilesecurityservice", "mobilesecurityserviceapp_cr": name}
}

//To transform the object into a string with its json
func getSdkConfigStringJsonFormat(sdk *models.SDKConfig) string {
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
