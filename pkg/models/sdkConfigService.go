package models

type SDKConfigService struct{
	Name                  string     `json:"name"`
	Host             	  string     `json:"host"`
}

func NewSDKConfigServices(serviceName, serviceHost string) *SDKConfigService {
	service := new(SDKConfigService)
	service.Name = serviceName
	service.Host = serviceHost
	return service
}
