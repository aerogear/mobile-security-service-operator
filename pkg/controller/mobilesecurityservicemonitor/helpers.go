package mobilesecurityservicemonitor

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//hasLabelSelector return true whn the LabelSelector has values specified in the CR
func hasLabelSelector(m *mobilesecurityservicev1alpha1.MobileSecurityServiceMonitor) bool {
	if len(m.Spec.LabelSelector.MatchLabels) > 0 {
		return true
	}
	return false
}

//hasNamespaceSelector return true when the NamespaceSelector a namespace specified in the CR
func hasNamespaceSelector(m *mobilesecurityservicev1alpha1.MobileSecurityServiceMonitor) bool {
	if len(m.Spec.NamespaceSelector) > 0 {
		return true
	}
	return false
}

//getMonitorListOps return the ListOptions according to the specifications made in the CR
func getListOptionsToFilterResources(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceMonitor, reqLogger logr.Logger) client.ListOptions {
	var listOps client.ListOptions
	if hasLabelSelector(instance) && hasNamespaceSelector(instance) {
		reqLogger.Info("Watching by NamespaceSelector and by the LabelSelector ...", "LabelSelector", instance.Spec.LabelSelector.MatchLabels, "NamespaceSelector", instance.Spec.NamespaceSelector)
		labelSelector := labels.SelectorFromSet(instance.Spec.LabelSelector.MatchLabels)
		listOps = client.ListOptions{Namespace: instance.Spec.NamespaceSelector, LabelSelector: labelSelector}
	} else if hasLabelSelector(instance) && !hasNamespaceSelector(instance) {
		reqLogger.Info("Watching by LabelSelector only ..." , "LabelSelector", instance.Spec.LabelSelector.MatchLabels )
		labelSelector := labels.SelectorFromSet(instance.Spec.LabelSelector.MatchLabels)
		listOps = client.ListOptions{LabelSelector: labelSelector}
	} else if !hasLabelSelector(instance) && hasNamespaceSelector(instance) {
		reqLogger.Info("Watching by NamespaceSelector only ...", "NamespaceSelector", instance.Spec.NamespaceSelector)
		listOps = client.ListOptions{Namespace: instance.Spec.NamespaceSelector}
	} else {
		reqLogger.Info("Not found Specification for LabelSelector or NamespaceSelector. Watching all namespaces ...")
		listOps = client.ListOptions{Namespace: ""}
	}
	return listOps
}
