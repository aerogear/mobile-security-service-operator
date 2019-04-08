package mobilesecurityservicebind

import (
	"bytes"
	"encoding/json"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/models"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"unsafe"
)


// Returns an string map with the labels which wil be associated to the kubernetes/openshift objects
// which will be created and managed by this operator
func getAppLabels(name string) map[string]string {
	return map[string]string{"app": "mobilesecurityservice", "mobilesecurityservicebind_cr": name}
}

func getAppLabelsSelectorMDCApp(labelSelector, valueSelector string) map[string]string {
	return map[string]string{labelSelector: valueSelector}
}

func getAppLabelsForSDKConfigMap(name string) map[string]string {
	return map[string]string{"app": "mobilesecurityservice", "mobilesecurityservicebind_cr": name, "name": name+"-sdk-config"}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

//To transform the object into a string with its json
func getSdkConfigStringJsonFormat(sdk *models.SDKConfig) string{
	jsonSdk, _ := json.Marshal(sdk)
	res:= strings.NewReader(string(jsonSdk))
	buf := new(bytes.Buffer)
	buf.ReadFrom(res)
	b := buf.Bytes()
	return *(*string)(unsafe.Pointer(&b))
}

//TODO: Implement this func to get all services available for this project when/if it started to have services
func getServices() []models.SDKConfigService{
	//service := *models.NewSDKConfigServices("","")
	res := []models.SDKConfigService{}
	//res = append(res, service)
	return res
}

// return properties for the response SDK
func getConfigMapSDKForMobileSecurityService(m *mobilesecurityservicev1alpha1.MobileSecurityServiceBind) map[string]string {
	url:= "http://" + getAppIngressHost(m)
	sdk := models.NewSDKConfig(m, url, getServices())
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

// return true when the namespace spec is filled
func hasWatchNamespaceSelector(m *mobilesecurityservicev1alpha1.MobileSecurityServiceBind) bool {
	if len(m.Spec.WatchNamespaceSelector) > 0 {
		return true
	}
	return false
}

//Return the List with Ops to tell what resources the Bind should watch
func getWatchListOps(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, reqLogger logr.Logger) client.ListOptions {
	var listOps client.ListOptions
	labelSelector := labels.SelectorFromSet(getAppLabelsSelectorMDCApp(instance.Spec.WatchKeyLabelSelector, instance.Spec.WatchValueLabelSelector))
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

//TODO: Centralized
// It will build the HOST for the router/ingress created for the Mobile Security Service App
func getAppIngressHost(m *mobilesecurityservicev1alpha1.MobileSecurityServiceBind) string {
	hostName := "mobile-security-service-app" + "." + m.Spec.ClusterHost + m.Spec.HostSufix
	return hostName;
}

func getURLServiceRestAPI(m *mobilesecurityservicev1alpha1.MobileSecurityServiceBind) string {
	return getAppIngressHost(m) + "/api/"
}

