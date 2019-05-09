package mobilesecurityserviceapp

import (
	"fmt"
	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"github.com/go-logr/logr"
	"context"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
)
//updateFinilizer returns error when the app still not deleted in the REST Service
func (r *ReconcileMobileSecurityServiceApp) updateFinilizer(serviceAPI string,reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp) error {
	if len(instance.GetFinalizers()) > 0 && instance.GetDeletionTimestamp() != nil {
		reqLogger.Info("Removing Finalizer for the MobileSecurityServiceApp")
		if app, err := fetchBindAppRestServiceByAppID(serviceAPI, instance, reqLogger);  err != nil || hasApp(app){
			if hasApp(app) {
				err := fmt.Errorf("App was found in the REST Service API")
				reqLogger.Error(err, "Unable to delete APP", "App.appId", instance.Spec.AppId)
				return err
			}
			return err
		}
		instance.SetFinalizers(nil)
		instance.SetDeletionTimestamp(nil)
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update finalizer for the MobileSecurityService App")
			return err
		}
	}
	return nil
}

//handleFinalizer check if has the finalizer and delete the app
func (r *ReconcileMobileSecurityServiceApp) handleFinalizer(serviceAPI string,reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp) error {
	// set up finalizers
	if len(instance.GetFinalizers()) > 0 && instance.GetDeletionTimestamp() != nil {
		//Check if App is delete into the REST Service
		if app, err := fetchBindAppRestServiceByAppID(serviceAPI, instance, reqLogger); err == nil {
			if hasApp(app) {
				if err := service.DeleteAppFromServiceByRestAPI(serviceAPI,  app.ID, reqLogger); err != nil {
					reqLogger.Error(err, "App was not delete on the Rest Service", "App.id",  app.ID)
					return err
				}
				return nil
			}
		}
	}
	return nil
}
