package mobilesecurityservicebind

import (
	"encoding/json"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/models"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	"github.com/go-logr/logr"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

//createAppByRestAPI create the app object in the service
func createAppByRestAPI(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, reqLogger logr.Logger) (reconcile.Result, error) {
	// Create the object and parse for JSON
	app := models.NewApp(instance)
	appJSON, err := json.Marshal(models.NewApp(instance))
	if err != nil {
		reqLogger.Error(err, "Error to transform the app object in JSON", "Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name, "AppJSON", appJSON, "App", app, "Error", err)
		return reconcile.Result{}, err
	}

	//Create the POST request
	url:= utils.GetRestAPIForApps(instance)
	req, err := http.NewRequest(http.MethodPost, url , strings.NewReader(string(appJSON)))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		reqLogger.Error(err, "Error when try to create the request", "HTTPMethod", http.MethodPost, "Request", req, "url", url, "error", err)
		return reconcile.Result{}, err
	}

	//Do the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		reqLogger.Error(err, "Error when try to do the request", "HTTPMethod", http.MethodPost, "Request", req, "Response", response, "url", url, "error", err)
		return reconcile.Result{}, err
	}
	defer response.Body.Close()

	reqLogger.Info("Created successfully app object in REST Service API",  "App:", app)
	return reconcile.Result{Requeue: true}, nil
}

//getAppFromServiceByRestApi returns the app object from the service
func getAppFromServiceByRestApi(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, reqLogger logr.Logger) (models.App, error) {
	// Create the object
	app := models.NewApp(instance)

	// Fill the record with the data from the JSON
	// Transform the body request in the version struct
	got := models.App{}

	if len(app.AppID) < 1 {
		reqLogger.Info( "Unable to get the AppID", "Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
		return got, nil
	}

	//Create the GET request
	url:= utils.GetRestAPIForApps(instance)+"?appId="+app.AppID
	req, err := http.NewRequest(http.MethodGet,  url, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		reqLogger.Error(err, "Error when try to create the request", "HTTPMethod", http.MethodGet, "Request", req, "url", url, "error", err)
		return got, err
	}

	//Do the request
	client := &http.Client{}
	response, err := client.Do(req)
	defer response.Body.Close()

	if err != nil {
		reqLogger.Error(err, "Error when try to do the request", "HTTPMethod", http.MethodGet, "Request", req, "url", url, "error", err)
		return got, err
	}

	if 200 != response.StatusCode {
		reqLogger.Info("The app was not found in the REST Service API", "HTTPMethod", http.MethodGet, "url", url)
		return got, nil
	}

	var obj []models.App
	err = json.NewDecoder(response.Body).Decode(&obj)

	if err != nil {
		reqLogger.Error(err, "Error when try to do decoder the body response", "HTTPMethod", http.MethodGet, "Request", req, "Response", response, "url", url, "error", err)
		return got, err
	}

	reqLogger.Info("App found in the Rest Service", "App", got)
	return obj[0], nil
}

//deleteAppFromServiceByRestAPI delete the app object in the service
func deleteAppFromServiceByRestAPI(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, app models.App, reqLogger logr.Logger ) (reconcile.Result, error) {
	//Create the DELETE request
	url:= utils.GetRestAPIForApps(instance)+"/"+app.ID
	req, err := http.NewRequest(http.MethodDelete, url ,nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		reqLogger.Error(err, "Error when try to create the request", "HTTPMethod", http.MethodDelete, "Request", req, "url", url, "error", err)
		return reconcile.Result{}, err
	}

	//Do the request
	client := &http.Client{}
	response, err := client.Do(req)
	defer response.Body.Close()

	if err != nil || 204 != response.StatusCode {
		reqLogger.Error(err, "Error when try to do the request", "HTTPMethod", http.MethodDelete, "Request", req, "url", url, "error", err)
		return reconcile.Result{}, err
	}

	reqLogger.Info("Deleted successfully app object in REST Service API",  "App:", app)
	return reconcile.Result{Requeue: true}, nil
}

//deleteAppFromServiceByRestAPI delete the app object in the service
func updateAppNameByRestAPI(instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBind, app models.App, reqLogger logr.Logger) (reconcile.Result, error) {
	//Create the DELETE request
	url:= utils.GetRestAPIForApps(instance)+"/"+app.ID
	appJSON, err := json.Marshal(models.NewApp(instance))

	if err != nil {
		reqLogger.Error(err, "Error to transform the app object in JSON", "AppJSON", appJSON, "App", app, "Error", err)
		return reconcile.Result{}, err
	}

	req, err := http.NewRequest(http.MethodPatch, url ,strings.NewReader(string(appJSON)))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		reqLogger.Error(err, "Error when try to create the request", "HTTPMethod", http.MethodPatch, "Request", req, "url", url, "error", err)
		return reconcile.Result{}, err
	}

	//Do the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil || 204 != response.StatusCode {
		reqLogger.Error(err, "Error when try to do the request", "HTTPMethod", http.MethodPatch, "Request", req, "url", url, "error", err)
		return reconcile.Result{}, err
	}
	defer response.Body.Close()

	reqLogger.Info("Deleted successfully app object in REST Service API",  "App:", app)
	return reconcile.Result{Requeue: true}, nil
}