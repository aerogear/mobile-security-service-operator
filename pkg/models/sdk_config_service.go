package models

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
)

const (
	ID = "security"
	Name = "security"
	Type = "security"
)

type SDKConfigService struct{
	ID					  string     `json:"id"`
	Name                  string     `json:"name"`
	Type                  string     `json:"type"`
	URL             	  string     `json:"url"`
	ConfigService         ConfigService `json:"config,omitempty"`
}

func NewSDKConfigServices(m *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, serviceUrl string) *SDKConfigService {
	service := new(SDKConfigService)
	service.ID = ID
	service.Name = Name
	service.Type = Type
	service.URL = serviceUrl
	service.ConfigService = *NewConfigService(service.URL)
	return service
}





