package models

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
)

type SDKConfig struct {
	URL string `json:"url"`
}

func NewSDKConfig(m *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, serviceUrl string) *SDKConfig {
	cfg := new(SDKConfig)
	cfg.URL = serviceUrl
	return cfg
}
