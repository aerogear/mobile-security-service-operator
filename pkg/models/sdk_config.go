package models

type SDKConfig struct {
	URL string `json:"url"`
}

func NewSDKConfig(serviceUrl string) *SDKConfig {
	return &SDKConfig{
		URL: serviceUrl,
	}
}
