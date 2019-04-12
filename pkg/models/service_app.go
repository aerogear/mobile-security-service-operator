package models

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
)

type App struct{
	ID                    string     `json:"id"`
	AppID                 string     `json:"appId"`
	AppName               string     `json:"appName,omitempty"`
	DeletedAt             string     `json:"deletedAt,omitempty"`
}

func NewApp(m *mobilesecurityservicev1alpha1.MobileSecurityServiceBind) *App {
	app := new(App)
	app.AppName = m.Spec.AppName
	app.AppID = m.Spec.AppId
	return app
}

//TODO: It should be removed when the PR: https://github.com/aerogear/mobile-security-service/pull/145 be merged