package mobilesecurityservicebind

import (
	"encoding/json"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/models"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Returns an string map with the labels which wil be associated to the kubernetes/openshift objects
// which will be created and managed by this operator
func getAppLabels(name string) map[string]string {
	return map[string]string{"app": "mobilesecurityservice", "mobilesecurityservicebind_cr": name}
}

func getAppLabelsForSDKConfigMap(name string) map[string]string {
	return map[string]string{"app": "mobilesecurityservice", "mobilesecurityservicebind_cr": name, "name": name+"-sdk-config"}
}

//To transform the object into a string with its json
func getSdkConfigStringJsonFormat(sdk *models.SDKConfig) string{
	jsonSdk, _ := json.MarshalIndent(sdk, "", "\t")
	return string(jsonSdk)
}

// return properties for the response SDK
func getConfigMapSDKForMobileSecurityService(m *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, pod corev1.Pod) map[string]string {
	sdk := models.NewSDKConfig(m, pod)
	return map[string]string{
		"SDKConfig": getSdkConfigStringJsonFormat(sdk),
	}
}
// return true when the key and label spec are filled
func hasWatchLabelSelectors(m *mobilesecurityservicev1alpha1.MobileSecurityServiceBind) bool {
	if len(m.Spec.WatchKeyLabelSelector) > 0 && len(m.Spec.WatchValueLabelSelector) > 0  {
		return true
	}
	return false
}

// return true when the key and label spec are filled
func hasBindLabelSelectors(m *mobilesecurityservicev1alpha1.MobileSecurityServiceBind) bool {
	if len(m.Spec.AppKeyLabelSelector) > 0 && len(m.Spec.AppValueLabelSelector) > 0  {
		return true
	}
	return false
}

// return true when the namespace spec is filled
func hasWatchNamespaceSelector(m *mobilesecurityservicev1alpha1.MobileSecurityServiceBind) bool {
	if len(m.Spec.WatchNamespaceSelector) > 0 {
		return true
	}
	return false
}

//Return the List with Ops to tell what resources the Bind should watch
func getAppWatchListOps(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, reqLogger logr.Logger) client.ListOptions {
	labelSelector := labels.SelectorFromSet(map[string]string{instance.Spec.WatchKeyLabelSelector: instance.Spec.WatchValueLabelSelector})
	return getListOptionsToFilterResources(instance, reqLogger, labelSelector)
}

func getListOptionsToFilterResources(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, reqLogger logr.Logger, labelSelector labels.Selector) client.ListOptions {
	var listOps client.ListOptions
	if hasWatchLabelSelectors(instance) && hasWatchNamespaceSelector(instance) {
		reqLogger.Info("Watching by WatchNamespaceSelector and by the WatchLabelSelectors ...")
		listOps = client.ListOptions{Namespace: instance.Spec.WatchNamespaceSelector, LabelSelector: labelSelector}
	} else if hasWatchLabelSelectors(instance) && !hasWatchNamespaceSelector(instance) {
		reqLogger.Info("Watching by WatchLabelSelectors only ...")
		listOps = client.ListOptions{LabelSelector: labelSelector}
	} else if !hasWatchLabelSelectors(instance) && hasWatchNamespaceSelector(instance) {
		reqLogger.Info("Watching by WatchNamespaceSelector only ...")
		listOps = client.ListOptions{Namespace: instance.Spec.WatchNamespaceSelector}
	} else {
		reqLogger.Info("Not found Specification for Watch Labels or Namespace. Watching in the same namespace where the BIND is installed ...")
		listOps = client.ListOptions{Namespace: instance.Namespace}
	}
	return listOps
}

// return true when the pod has the labels to tell that the pod is bind to the service
func isBind(pod corev1.Pod, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind) bool {
	var isBind = false
	if val, ok := pod.Labels[instance.Spec.AppKeyLabelSelector]; ok {
		if val == instance.Spec.AppValueLabelSelector {
			isBind = true
		}
	}
	return isBind
}