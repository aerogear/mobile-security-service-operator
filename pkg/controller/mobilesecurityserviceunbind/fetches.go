package mobilesecurityserviceunbind

import (
	"context"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"github.com/aerogear/mobile-security-service-operator/pkg/models"
)

// Request object not found, could have been deleted after reconcile request.
// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
func (r *ReconcileMobileSecurityServiceUnbind) fetch(request reconcile.Request, reqLogger logr.Logger) (*mobilesecurityservicev1alpha1.MobileSecurityServiceUnbind, error) {
	instance := &mobilesecurityservicev1alpha1.MobileSecurityServiceUnbind{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if errors.IsNotFound(err) {
		// Return and don't create
		reqLogger.Info( "Return Mobile Security Service Unbind instance")
		return instance, nil
	}
	// Error reading the object - create the request.
	reqLogger.Error(err, "Failed to get Mobile Security Service Unbind")
	return instance, err
}

//fetchBindAppRestServiceByAppID return app struct from Mobile Security Service Project/REST API or error
func fetchBindAppRestServiceByAppID(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceUnbind, reqLogger logr.Logger) (models.App, error){
	return service.GetAppFromServiceByRestApi(instance.Spec.Protocol, instance.Spec.ClusterHost, instance.Spec.HostSufix, instance.Spec.AppId, reqLogger)
}