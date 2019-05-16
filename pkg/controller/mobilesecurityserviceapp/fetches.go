package mobilesecurityserviceapp

import (
	"context"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/models"
	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Request object not found, could have been deleted after reconcile request.
// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
func (r *ReconcileMobileSecurityServiceApp) fetchInstance(reqLogger logr.Logger, request reconcile.Request) (*mobilesecurityservicev1alpha1.MobileSecurityServiceApp, error) {
	instance := &mobilesecurityservicev1alpha1.MobileSecurityServiceApp{}
	//Fetch the MobileSecurityServiceApp instance
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	return instance, err
}

//fetchSDKConfigMap returns the config map resource created for this instance
func (r *ReconcileMobileSecurityServiceApp) fetchSDKConfigMap(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp) (*corev1.ConfigMap, error) {
	reqLogger.Info("Checking if the ConfigMap already exists")
	configMap := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: getSDKConfigMapName(instance), Namespace: instance.Namespace}, configMap)
	return configMap, err
}

//fetchBindAppRestServiceByAppID return app struct from Mobile Security Service Project/REST API or error
var fetchBindAppRestServiceByAppID = func(serviceURL string, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, reqLogger logr.Logger) (*models.App, error) {
	return service.GetAppFromServiceByRestApi(serviceURL, instance.Spec.AppId, reqLogger)
}
