package utils

import (
	"os"
	"testing"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	mssInstance = mobilesecurityservicev1alpha1.MobileSecurityService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service",
			Namespace: "mobile-security-service-operator",
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceSpec{
			Size:            1,
			MemoryLimit:     "512Mi",
			MemoryRequest:   "512Mi",
			ClusterProtocol: "http",
			ConfigMapName:   "mss-config",
			RouteName:       "route",
		},
	}

	route = routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mssInstance.Spec.RouteName,
			Namespace: mssInstance.Namespace,
			Labels:    map[string]string{"app": "mobilesecurityservice", "mobilesecurityservice_cr": mssInstance.Name},
		},
		Status: routev1.RouteStatus{
			Ingress: []routev1.RouteIngress{
				{
					Host: "testhost",
				},
			},
		},
	}
)

func TestGetAppNamespaces(t *testing.T) {
	type fields struct {
		envVar string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Should return namespace",
			want: "apps-namespace",
			fields: fields{
				envVar: "apps-namespace",
			},
		},
		{
			name:    "Should return error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// first, unset any env that may be lying around from the previous case
			os.Unsetenv(AppNamespaceEnvVar)

			if tt.fields.envVar != "" {
				os.Setenv(AppNamespaceEnvVar, tt.fields.envVar)
			}

			got, err := GetAppNamespaces()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAppNamespaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAppNamespaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidAppNamespace(t *testing.T) {
	type fields struct {
		envVar string
	}
	type args struct {
		namespace string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Should be a valid app namespace",
			fields: fields{
				envVar: "apps-namespace",
			},
			args: args{
				namespace: "apps-namespace",
			},
			want: true,
		},
		{
			name: "Should find a valid app namespace in a dlimited string",
			fields: fields{
				envVar: "hello-world;apps-namespace;another-namespace",
			},
			args: args{
				namespace: "apps-namespace",
			},
			want: true,
		},
		{
			name: "Should be an invalid namespace",
			fields: fields{
				envVar: "hello-world",
			},
			args: args{
				namespace: "apps-namespace",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Should return local namespace with no name is set",
			args: args{
				namespace: OperatorNamespaceForLocalEnv,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// clear old env var
			os.Unsetenv(AppNamespaceEnvVar)

			if tt.fields.envVar != "" {
				os.Setenv(AppNamespaceEnvVar, tt.fields.envVar)
			}

			got, err := IsValidAppNamespace(tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidAppNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsValidAppNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidOperatorNamespace(t *testing.T) {
	type args struct {
		namespace string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "should check and return true",
			args: args{
				namespace: "mobile-security-service",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsValidOperatorNamespace(tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidOperatorNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsValidOperatorNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInitPublicURL(t *testing.T) {
	type args struct {
		instance *mobilesecurityservicev1alpha1.MobileSecurityService
		route    *routev1.Route
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should check and return true",
			args: args{
				instance: &mssInstance,
				route:    &route,
			},
			want: "http://testhost/init",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetInitPublicURL(tt.args.route, tt.args.instance)
			if got != tt.want {
				t.Errorf("TestGetInitPublicURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceAPIURL(t *testing.T) {
	type args struct {
		instance *mobilesecurityservicev1alpha1.MobileSecurityService
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should check and return true",
			args: args{
				instance: &mssInstance,
			},
			want: "http://mobile-security-service-application:0/api",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetServiceAPIURL(tt.args.instance)
			if got != tt.want {
				t.Errorf("TestGetServiceAPIURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
