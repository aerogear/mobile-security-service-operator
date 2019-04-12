package utils

import (
	corev1 "k8s.io/api/core/v1"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
)

const APP_URL =  "mobile-security-service-app"

func GetAppIngressURL(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind) string {
	return instance.Spec.Protocol +"://" + APP_URL + "." + instance.Spec.ClusterHost + instance.Spec.HostSufix
}

func GetAppIngress(host, hostSufix string) string {
	return APP_URL + "." + host + hostSufix
}

// getPodNames returns the pod names of the array of pods passed in
func GetPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func GetRestAPIForApps(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind) string {
	return GetAppIngressURL(instance) + "/api/apps"
}