package mobilesecurityservicebind

import (
	"context"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"github.com/aerogear/mobile-security-service-operator/pkg/models"
)

// Request object not found, could have been deleted after reconcile request.
// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
func (r *ReconcileMobileSecurityServiceBind) fetch(request reconcile.Request, reqLogger logr.Logger) (*mobilesecurityservicev1alpha1.MobileSecurityServiceBind, error) {
	instance := &mobilesecurityservicev1alpha1.MobileSecurityServiceBind{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if errors.IsNotFound(err) {
		// Return and don't create
		reqLogger.Info( "Return Mobile Security Service Bind instance")
		return instance, nil
	}
	// Error reading the object - create the request.
	reqLogger.Error(err, "Failed to get Mobile Security Service Bind")
	return instance, err
}

//fetchSDKConfigMap returns the config map resource created for this instance
func (r *ReconcileMobileSecurityServiceBind) fetchSDKConfigMap(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind) (*corev1.ConfigMap, error) {
	reqLogger.Info("Checking if the ConfigMap already exists")
	configMap := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: getConfigMapName(instance), Namespace: instance.Namespace}, configMap)
	return configMap, err
}

//fetchBindAppRestServiceByAppID return app struct from Mobile Security Service Project/REST API or error
func fetchBindAppRestServiceByAppID(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, reqLogger logr.Logger) (models.App, error){
	return service.GetAppFromServiceByRestApi(instance.Spec.Protocol, instance.Spec.ClusterHost, instance.Spec.HostSufix, instance.Spec.AppId, reqLogger)
}