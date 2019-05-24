package mobilesecurityserviceapp

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//updateSDKConfigMapStatus returns error when status regards the ConfigMap resource could not be updated
func (r *ReconcileMobileSecurityServiceApp) updateSDKConfigMapStatus(reqLogger logr.Logger, request reconcile.Request) (*corev1.ConfigMap, error) {
	reqLogger.Info("Updating SDKConfigMap Status for the MobileSecurityServiceApp")

	// Get the latest version of CR
	instance, err := r.fetchInstance(reqLogger, request)
	if err != nil {
		return &corev1.ConfigMap{}, err
	}

	// Get SDKConfigMap object
	SDKConfigMapStatus, err := r.fetchConfigMap(reqLogger, instance)
	if err != nil {
		reqLogger.Error(err, "Failed to get SDKConfigMap for Status", "MobileSecurityServiceApp.Namespace", instance.Namespace, "MobileSecurityServiceApp.Name", instance.Name)
		return SDKConfigMapStatus, err
	}

	//Update CR Status with SDKConfigMap name
	if SDKConfigMapStatus.Name != instance.Status.SDKConfigMapName {
		// Get the latest version of CR
		instance, err := r.fetchInstance(reqLogger, request)
		if err != nil {
			return &corev1.ConfigMap{}, err
		}

		// Set the data
		instance.Status.SDKConfigMapName = SDKConfigMapStatus.Name

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update SDKConfigMap Status for the MobileSecurityServiceApp")
			return SDKConfigMapStatus, err
		}
	}
	return SDKConfigMapStatus, nil
}

//updateAppStatus returns error when status regards the all required resources could not be updated
func (r *ReconcileMobileSecurityServiceApp) updateBindStatus(serviceURL string, reqLogger logr.Logger, SDKConfigMapStatus *corev1.ConfigMap, request reconcile.Request) error {
	reqLogger.Info("Updating Bind App Status for the MobileSecurityServiceApp")

	// Get the latest version of CR
	instance, err := r.fetchInstance(reqLogger, request)
	if err != nil {
		return err
	}

	// Get App created in the Rest Service
	app, err := fetchBindAppRestServiceByAppID(serviceURL, instance, reqLogger)
	if err != nil {
		reqLogger.Error(err, "Failed to get App for Status", "MobileSecurityServiceApp.Namespace", instance.Namespace, "MobileSecurityServiceApp.Name", instance.Name)
		return err
	}

	// Check if the ConfigMap and the App is created in the Rest Service
	if len(SDKConfigMapStatus.UID) < 1 && app.ID == "" {
		err := fmt.Errorf("Failed to get OK Status for MobileSecurityService Bind.")
		reqLogger.Error(err, "One of the resources are not created", "MobileSecurityServiceApp.Namespace", instance.Namespace, "MobileSecurityServiceApp.Name", instance.Name)
		return err
	}
	status := "OK"

	//Update Bind CR Status with OK
	if !reflect.DeepEqual(status, instance.Status.BindStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchInstance(reqLogger, request)
		if err != nil {
			return err
		}

		// Set the data
		instance.Status.BindStatus = status

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Status for the MobileSecurityService Bind")
			return err
		}
	}
	return nil
}

// updateBindStatusWithInvalidNamespace returns error when status regards the all required resources could not be updated
// DEPRECATED
func (r *ReconcileMobileSecurityServiceApp) updateBindStatusWithInvalidNamespace(reqLogger logr.Logger, request reconcile.Request) error {
	reqLogger.Info("Updating Bind App Status for the MobileSecurityServiceApp")

	// Get the latest version of CR
	instance, err := r.fetchInstance(reqLogger, request)
	if err != nil {
		return err
	}

	status := "Invalid Namespace"

	//Update Bind CR Status with OK
	if !reflect.DeepEqual(status, instance.Status.BindStatus) {
		// Get the latest version of the CR in order to try to avoid errors when try to update the CR
		instance, err := r.fetchInstance(reqLogger, request)
		if err != nil {
			return err
		}

		// Set the data
		instance.Status.BindStatus = status

		// Update the CR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Status for the MobileSecurityService Bind")
			return err
		}
	}
	return nil
}
