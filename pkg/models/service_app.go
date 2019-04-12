package models

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
)

type App struct{
	ID                    string     `json:"id"`
	AppID                 string     `json:"appId"`
	AppName               string     `json:"appName,omitempty"`
	DeletedAt             string     `json:"deletedAt,omitempty"`
}

func NewApp(m *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, pod corev1.Pod) *App {
	app := new(App)
	app.AppName = utils.GetAppNameByPodLabel(pod, m)
	app.AppID = utils.GetAppIdByPodLabel(pod, m)
	return app
}

//TODO: It should be removed when the PR: https://github.com/aerogear/mobile-security-service/pull/145 be merged