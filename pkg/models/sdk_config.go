package models

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
)

type SDKConfig struct{
	Version 			  string     `json:"version"`
	Name                  string     `json:"name"`
	Namespace             string     `json:"namespace"`
	Host             	  string     `json:"host"`
	ClientID              string     `json:"clientId"`
	Services              []SDKConfigService `json:"services,omitempty"`
}

func NewSDKConfig(m *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, pod corev1.Pod) *SDKConfig {

	cfg := new(SDKConfig)
	cfg.Version = "1.0.0"
	cfg.Name = utils.GetAppNameByPodLabel(pod, m)
	cfg.Namespace = pod.Namespace
	cfg.Host = m.Spec.ClusterHost
	cfg.Services = getServices(m)
	return cfg
}

//return the Service data for the SDK ConfigMap
func getServices(m *mobilesecurityservicev1alpha1.MobileSecurityServiceBind) []SDKConfigService{
	service := *NewSDKConfigServices(m)
	res := []SDKConfigService{}
	res = append(res, service)
	return res
}