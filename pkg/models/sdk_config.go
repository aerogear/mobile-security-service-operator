package models

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
)

type SDKConfig struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	URL           string        `json:"url"`
	ConfigService ConfigService `json:"config,omitempty"`
}

func NewSDKConfig(m *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, serviceUrl string) *SDKConfig {
	cfg := new(SDKConfig)
	cfg.ID = m.Spec.AppId
	cfg.Name = m.Spec.AppName
	cfg.Type = m.Name
	cfg.URL = serviceUrl
	cfg.ConfigService = *NewConfigService(serviceUrl)
	return cfg
}

type ConfigService struct {
	MobileSecurityServiceURL string `json:"mobile-security-server-url"`
}

func NewConfigService(url string) *ConfigService {
	service := new(ConfigService)
	service.MobileSecurityServiceURL = url
	return service
}
