package mobilesecurityserviceapp

import (
	"context"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//addFinalizer will add the Finalizer metadata in the Mobile Security Service App CR
func (r *ReconcileMobileSecurityServiceApp) addFinalizer(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp,  request reconcile.Request) error {
	if len(instance.GetFinalizers()) < 1 && instance.GetDeletionTimestamp() == nil {
		reqLogger.Info("Adding Finalizer for the MobileSecurityServiceApp")

		// Get the latest version of CR
		instance, err := r.fetchInstance(reqLogger, request)
		if err != nil {
			return err
		}

		//Set finalizer string/metadata
		instance.SetFinalizers([]string{FINALIZER})

		//Update CR
		err = r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update MobileSecurityService App CR with  finalizer")
			return err
		}
	}
	return nil
}


//removeFinalizer returns error when the app still not deleted in the REST Service
func (r *ReconcileMobileSecurityServiceApp) removeFinalizer(serviceAPI string,reqLogger logr.Logger, request reconcile.Request) error {

	// Get the latest version of CR
	instance, err := r.fetchInstance(reqLogger, request)
	if err != nil {
		return err
	}

	if len(instance.GetFinalizers()) > 0 && instance.GetDeletionTimestamp() != nil {
		reqLogger.Info("Removing Finalizer for the MobileSecurityServiceApp")
		if app, err := fetchBindAppRestServiceByAppID(serviceAPI, instance, reqLogger);  err != nil || app.ID != "" {
			reqLogger.Error(err, "Unable to delete app", "App.appId", instance.Spec.AppId, "app.ID", app.ID)
			return err
		}

		//Remove finalizer
		instance.SetFinalizers(nil)

		//Update CR
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update MobileSecurityService App CR with finalizer")
			return err
		}
	}
	return nil
}
