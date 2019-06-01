package mobilesecurityserviceapp

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"github.com/aerogear/mobile-security-service/pkg/models"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Request object not found, could have been deleted after reconcile request.
// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
func (r *ReconcileMobileSecurityServiceApp) fetchMssAppInstance(reqLogger logr.Logger, request reconcile.Request) (*mobilesecurityservicev1alpha1.MobileSecurityServiceApp, error) {
	mssApp := &mobilesecurityservicev1alpha1.MobileSecurityServiceApp{}
	//Fetch the MobileSecurityServiceApp mssApp
	err := r.client.Get(context.TODO(), request.NamespacedName, mssApp)
	return mssApp, err
}

// fetchConfigMap returns the config map resource created for this instance
func (r *ReconcileMobileSecurityServiceApp) fetchConfigMap(reqLogger logr.Logger, mssApp *mobilesecurityservicev1alpha1.MobileSecurityServiceApp) (*corev1.ConfigMap, error) {
	reqLogger.Info("Checking if the ConfigMap already exists")
	configMap := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: getSDKConfigMapName(mssApp), Namespace: mssApp.Namespace}, configMap)
	return configMap, err
}

// fetchConfigMapListByLabels returns a list of config map resource created for this instance instance with all labels less the namein order to check if the name was changed
func (r *ReconcileMobileSecurityServiceApp) fetchConfigMapListByLabels(reqLogger logr.Logger, mssApp *mobilesecurityservicev1alpha1.MobileSecurityServiceApp) (*corev1.ConfigMapList, error) {
	reqLogger.Info("Checking if the ConfigMap already exists")
	configMapList := &corev1.ConfigMapList{}
	listOps := &client.ListOptions{}
	listOps.InNamespace(mssApp.Namespace)
	listOps.MatchingLabels(getLabelsToFetch(mssApp))
	err := r.client.List(context.TODO(), listOps, configMapList)
	return configMapList, err
}

// fetchBindAppRestServiceByAppID return app struct from Mobile Security Service Project/REST API or error
var fetchBindAppRestServiceByAppID = func(serviceURL string, mssApp *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, reqLogger logr.Logger) (*models.App, error) {
	return service.GetAppFromServiceByRestApi(serviceURL, mssApp.Spec.AppId, reqLogger)
}
