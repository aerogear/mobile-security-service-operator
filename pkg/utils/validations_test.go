package utils

import (
	"testing"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
)

func TestCheckClusterProtocol(t *testing.T) {
	type args struct {
		serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return false when ClusterProtocol is not defined",
			args: args{
				serviceInstance: &mobilesecurityservicev1alpha1.MobileSecurityService{},
			},
			want: false,
		},
		{
			name: "should return false when ClusterProtocol is invalid",
			args: args{
				serviceInstance: &mobilesecurityservicev1alpha1.MobileSecurityService{
					Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceSpec{
						ClusterProtocol: "ws",
					},
				},
			},
			want: false,
		},
		{
			name: "should return true when ClusterProtocol is valid",
			args: args{
				serviceInstance: &mobilesecurityservicev1alpha1.MobileSecurityService{
					Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceSpec{
						ClusterProtocol: "https",
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reqLogger := log.WithValues("Validations Test")

			if got := CheckClusterProtocol(tt.args.serviceInstance, reqLogger); got != tt.want {
				t.Errorf("CheckClusterProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}
