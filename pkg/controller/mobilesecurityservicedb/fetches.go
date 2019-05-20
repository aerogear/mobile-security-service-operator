package mobilesecurityservicedb

import (
	"context"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Request object not found, could have been deleted after reconcile request.
// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
func (r *ReconcileMobileSecurityServiceDB) fetchInstance(reqLogger logr.Logger, request reconcile.Request) (*mobilesecurityservicev1alpha1.MobileSecurityServiceDB, error) {
	instance := &mobilesecurityservicev1alpha1.MobileSecurityServiceDB{}
	//Fetch the MobileSecurityServiceDB instance
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	return instance, err
}

//fetchDBService returns the service resource created for this instance
func (r *ReconcileMobileSecurityServiceDB) fetchDBService(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceDB) (*corev1.Service, error) {
	reqLogger.Info("Checking if the service already exists")
	service := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, service)
	return service, err
}

//fetchDBDeployment returns the deployment resource created for this instance
func (r *ReconcileMobileSecurityServiceDB) fetchDBDeployment(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceDB) (*v1beta1.Deployment, error) {
	reqLogger.Info("Checking if the deployment already exists")
	deployment := &v1beta1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, deployment)
	return deployment, err
}

//fetchDBPersistentVolumeClaim returns the PersistentVolumeClaim resource created for this instance
func (r *ReconcileMobileSecurityServiceDB) fetchDBPersistentVolumeClaim(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceDB) (*corev1.PersistentVolumeClaim, error) {
	reqLogger.Info("Checking if the DB PersistentVolumeClaim already exists")
	pvc := &corev1.PersistentVolumeClaim{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, pvc)
	return pvc, err
}
