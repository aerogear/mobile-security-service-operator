package service

import (
	"encoding/json"
	"fmt"
	"github.com/aerogear/mobile-security-service-operator/pkg/models"
	"github.com/go-logr/logr"
	"net/http"
	"strings"
)

//DeleteAppFromServiceByRestAPI delete the app object in the service
func DeleteAppFromServiceByRestAPI(serviceAPI string, id string, reqLogger logr.Logger ) error {
	reqLogger.Info( "Calling REST Service to DELETE app", "serviceAPI", serviceAPI, "App.id", id )
	//Create the DELETE request
	url:= serviceAPI + "/apps/" + id
	req, err := http.NewRequest(http.MethodDelete, url ,nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		reqLogger.Error(err, "Unable to create DELETE request", "HTTPMethod", http.MethodDelete, "url", url)
		return err
	}

	//Do the request
	client := &http.Client{}
	response, err := client.Do(req)
	defer response.Body.Close()

	if err != nil || 204 != response.StatusCode {
		reqLogger.Error(err, "HTTP StatusCode not expected", "HTTPMethod", http.MethodDelete, "url", url, "response.StatusCode", response.StatusCode)
		return err
	}

	reqLogger.Info("Successfully deleted app  ...",  "App.Id:", id)
	return nil
}


//CreateAppByRestAPI create the app object in the service
func CreateAppByRestAPI(serviceAPI string, app models.App, reqLogger logr.Logger) error {
	reqLogger.Info( "Calling Service to POST app", "serviceAPI", serviceAPI, "App", app )

	// Create the object and parse for JSON
	appJSON, err := json.Marshal(app)
	if err != nil {
		reqLogger.Error(err, "Error to transform the app object in JSON", "App", app, "Error", err)
		return err
	}

	//Create the POST request
	url:= serviceAPI + "/apps"
	req, err := http.NewRequest(http.MethodPost, url , strings.NewReader(string(appJSON)))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		reqLogger.Error(err, "Unable to create POST request", "HTTPMethod", http.MethodPost, "url", url)
		return err
	}

	//Do the request
	client := &http.Client{}
	response, err := client.Do(req)
	defer response.Body.Close()

	if err != nil || 201 != response.StatusCode {
		reqLogger.Error(err, "HTTP StatusCode not expected", "HTTPMethod", http.MethodPost, "url", url, "response.StatusCode", response.StatusCode )
		return err
	}

	reqLogger.Info("Successfully created app  ...",  "App:", app)
	return nil
}

//GetAppFromServiceByRestApi returns the app object from the service
func GetAppFromServiceByRestApi(serviceAPI string, appId string, reqLogger logr.Logger) (models.App, error) {
	reqLogger.Info( "Calling REST API to GET app", "serviceAPI", serviceAPI, "App.appId", appId )

	// Fill the record with the data from the JSON
	// Transform the body request in the version struct
	got := models.App{}

	if appId == "" {
		err := fmt.Errorf( "App without AppId", "App.AppId", appId)
		return got, err
	}

	//Create the GET request
	url:= serviceAPI + "/apps" +"?appId="+appId
	req, err := http.NewRequest(http.MethodGet,  url, nil)
	reqLogger.Info("URL to get", "HTTPMethod", http.MethodGet, "url", url)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		reqLogger.Error(err, "Unable to create GET request", "HTTPMethod", http.MethodGet, "Request", req, "url", url)
		return got, err
	}

	//Do the request
	client := &http.Client{}
	response, err := client.Do(req)
	defer response.Body.Close()
	if err != nil {
		reqLogger.Error(err, "Unable to execute GET request", "HTTPMethod", http.MethodGet, "url", url, "response.StatusCode", response.StatusCode)
		return got, err
	}

	if 200 != response.StatusCode && 204 != response.StatusCode {
		err := fmt.Errorf( "HTTP StatusCode not expected")
		reqLogger.Error(err, "HTTP StatusCode not expected", "HTTPMethod", http.MethodGet, "url", url, "response.StatusCode", response.StatusCode)
		return got, err
	}

	var obj []models.App
	err = json.NewDecoder(response.Body).Decode(&obj)
	if err != nil {
		reqLogger.Error(err, "Error when try to do decode the body response", "HTTPMethod", http.MethodGet, "url", url, "Response.Body", response.Body )
		return got, err
	}

	if 204 == response.StatusCode || got.ID == "" {
		reqLogger.Info("The app was not found in the REST Service API", "HTTPMethod", http.MethodGet, "url", url)
		return got, nil
	}

	reqLogger.Info("App found in the Service", "App", obj[0])
	return obj[0], nil
}

//UpdateAppNameByRestAPI will update name of the APP in the Service
func UpdateAppNameByRestAPI(serviceAPI string, app models.App, reqLogger logr.Logger) error {
	reqLogger.Info( "Calling Service to update app name", "serviceAPI", serviceAPI, "App", app )

	//Create the DELETE request
	url:= serviceAPI + "/apps" +app.ID
	appJSON, err := json.Marshal(app)

	if err != nil {
		reqLogger.Error(err, "Error to transform the app object in JSON", "AppJSON", appJSON, "App", app, "Error", err)
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, url ,strings.NewReader(string(appJSON)))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		reqLogger.Error(err, "Unable to create PATCH request to update app name", "HTTPMethod", http.MethodPatch, "url", url)
		return err
	}

	//Do the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil || 204 != response.StatusCode {
		reqLogger.Error(err, "HTTP StatusCode not expected", "HTTPMethod", http.MethodPatch, "url", url, "response.StatusCode", response.StatusCode)
		return err
	}
	defer response.Body.Close()

	reqLogger.Info("Successfully updated app name ...",  "App", app)
	return nil
}