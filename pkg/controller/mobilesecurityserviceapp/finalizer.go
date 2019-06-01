package mobilesecurityserviceapp

import (
	"context"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// addFinalizer will add the Finalizer metadata in the Mobile Security Service App CR
func (r *ReconcileMobileSecurityServiceApp) addFinalizer(reqLogger logr.Logger, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, request reconcile.Request) error {
	if len(instance.GetFinalizers()) < 1 && instance.GetDeletionTimestamp() == nil {
		reqLogger.Info("Adding Finalizer for the MobileSecurityServiceApp")

		// Get the latest version of CR
		instance, err := r.fetchAppInstance(reqLogger, request)
		if err != nil {
			return err
		}

		//Set finalizer string/metadata
		instance.SetFinalizers([]string{FinalizerMetadata})

		//Update CR
		err = r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update MobileSecurityService App CR with  finalizer")
			return err
		}
	}
	return nil
}

// handleFinalizer returns error when the app still not deleted in the REST Service
func (r *ReconcileMobileSecurityServiceApp) handleFinalizer(serviceAPI string, reqLogger logr.Logger, request reconcile.Request) error {

	// Get the latest version of CR
	app, err := r.fetchAppInstance(reqLogger, request)
	if err != nil {
		return err
	}

	if len(app.GetFinalizers()) > 0 && app.GetDeletionTimestamp() != nil {
		reqLogger.Info("Removing Finalizer for the MobileSecurityServiceApp")
		if appService, err := fetchBindAppRestServiceByAppID(serviceAPI, app, reqLogger); err != nil || appService.ID != "" {
			reqLogger.Error(err, "Unable to delete app", "App.appId", app.Spec.AppId, "app.ID", appService.ID)
			return err
		}

		if err := r.removeFinalizerFromCR(app); err != nil {
			reqLogger.Error(err, "Failed to update MobileSecurityService App CR with finalizer")
			return err
		}
	}
	return nil
}

// removeFinalizerFromCR return an error when is not possible remove the finalizer metadata from the app instance
func (r *ReconcileMobileSecurityServiceApp) removeFinalizerFromCR(app *mobilesecurityservicev1alpha1.MobileSecurityServiceApp) error {
	//Remove finalizer
	app.SetFinalizers(nil)

	//Update CR
	err := r.client.Update(context.TODO(), app)
	if err != nil {
		return err
	}
	return nil
}
