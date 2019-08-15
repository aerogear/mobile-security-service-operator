package mobilesecurityserviceapp

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobile-security-service/v1alpha1"
	"testing"
)

func Test_hasMandatorySpecs(t *testing.T) {
	type args struct {
		instance        *mobilesecurityservicev1alpha1.MobileSecurityServiceApp
		serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return true when instance has an AppId and AppName",
			args: args{
				instance: &instance,
			},
			want: true,
		},
		{
			name: "should return false when instance has no AppId",
			args: args{
				instance: &instanceNoAppName,
			},
			want: false,
		},
		{
			name: "should return false when instance has no AppName",
			args: args{
				instance: &instanceNoAppId,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqLogger := log.WithValues("Request.Namespace", &instance.Namespace, "Request.Name", &instance.Name)

			if got := hasMandatorySpecs(tt.args.instance, reqLogger); got != tt.want {
				t.Errorf("hasMandatorySpecs() = %v, want %v", got, tt.want)
			}
		})
	}
}
