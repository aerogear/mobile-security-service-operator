package mobilesecurityserviceapp

import (
	"context"
	"encoding/json"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/models"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/types"
)

const SdkConfigMapSufix = "-security"
const FinalizerMetadata = "finalizer.mobile-security-service.aerogear.com"

// Returns an string map with the labels which wil be associated to the kubernetes/openshift objects
// which will be created and managed by this operator
func getSDKAppLabels(m *mobilesecurityservicev1alpha1.MobileSecurityServiceApp) map[string]string {
	return map[string]string{"app": "mobilesecurityservice", "mobilesecurityserviceapp_cr": m.Name, "appname": m.Spec.AppName, "appid": m.Spec.AppId}
}

//  Returns an string map with the labels which wil be used to check if has an outdated configMap created
func getLabelsToFetch(m *mobilesecurityservicev1alpha1.MobileSecurityServiceApp) map[string]string {
	return map[string]string{"app": "mobilesecurityservice", "mobilesecurityserviceapp_cr": m.Name, "appid": m.Spec.AppId}
}

//To transform the object into a string with its json
func getSdkConfigStringJsonFormat(sdk *models.SDKConfig) string {
	jsonSdk, _ := json.MarshalIndent(sdk, "", "\t")
	return string(jsonSdk)
}

// return properties for the response SdkConfigMapSufix
func getConfigMapSDKForMobileSecurityService(m *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, serviceURL string) map[string]string {
	sdk := models.NewSDKConfig(serviceURL)
	return map[string]string{
		"SDKConfig": getSdkConfigStringJsonFormat(sdk),
	}
}

// return properties for the response SdkConfigMapSufix
func getSDKConfigMapName(m *mobilesecurityservicev1alpha1.MobileSecurityServiceApp) string {
	return m.Spec.AppName + SdkConfigMapSufix
}

// hasConditionsToBeDeleted will return true if the Service instance was not found and/or is marked to be deleted
// OR
// if the APP CR was marked to be deleted
func hasConditionsToBeDeleted(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, mssInstance *mobilesecurityservicev1alpha1.MobileSecurityService) bool {
	//Check if the APP CR was marked to be deleted
	isAppMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	hasFinalizer := len(instance.GetFinalizers()) > 0
	isMssInstanceDeleted := mssInstance == nil
	isMssInstanceMarkedToBeDeleted := mssInstance.GetDeletionTimestamp() != nil
	return (isAppMarkedToBeDeleted && hasFinalizer) || isMssInstanceDeleted || isMssInstanceMarkedToBeDeleted
}

// isMobileSecurityServiceDeleted return true if it is not found because was deleted and/or was marked to be deleted
func (r *ReconcileMobileSecurityServiceApp) isMobileSecurityServiceDeleted(operatorNamespace string, mssInstance *mobilesecurityservicev1alpha1.MobileSecurityService) bool {
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: utils.MobileSecurityServiceCRName, Namespace: operatorNamespace}, mssInstance); err != nil || mssInstance.GetDeletionTimestamp() != nil {
		return true
	}
	return false
}
