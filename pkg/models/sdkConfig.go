package models

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
)

type SDKConfig struct{
	Version 			  string     `json:"name"`
	Name                  string     `json:"name"`
	Namespace             string     `json:"namespace"`
	Host             	  string     `json:"host"`
	Services              []SDKConfigService `json:"services,omitempty"`
}

func NewSDKConfig(m *mobilesecurityservicev1alpha1.MobileSecurityService, host string, services []SDKConfigService) *SDKConfig {
	cfg := new(SDKConfig)
	cfg.Name = m.Name
	cfg.Namespace = m.Namespace
	cfg.Host = host
	cfg.Services = services
	return cfg
}