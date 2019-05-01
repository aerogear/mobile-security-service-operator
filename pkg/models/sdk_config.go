package models

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
)

type SDKConfig struct{
	Version 			  string     `json:"version"`
	Name                  string     `json:"name"`
	Namespace             string     `json:"namespace"`
	ClientID              string     `json:"clientId"`
	Services              []SDKConfigService `json:"services,omitempty"`
}

func NewSDKConfig(m *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, serviceUrl string) *SDKConfig {

	cfg := new(SDKConfig)
	cfg.Version = "1.0.0"
	cfg.Name = m.Spec.AppName
	cfg.ClientID = m.Spec.AppId
	cfg.Namespace = m.Namespace
	cfg.Services = getServices(m, serviceUrl)
	return cfg
}

//return the Service data for the SDK ConfigMap
func getServices(m *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, serviceUrl string) []SDKConfigService{
	service := *NewSDKConfigServices(m, serviceUrl)
	res := []SDKConfigService{}
	res = append(res, service)
	return res
}