package utils

import (
	"os"
	"testing"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetRouteName(t *testing.T) {
	type args struct {
		m *mobilesecurityservicev1alpha1.MobileSecurityService
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should use RouteName",
			args: args{
				m: &mobilesecurityservicev1alpha1.MobileSecurityService{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mss-operator",
					},
					Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceSpec{
						RouteName: "mss-route",
					},
				},
			},
			want: "mss-route",
		},
		{
			name: "should use Instance name",
			args: args{
				m: &mobilesecurityservicev1alpha1.MobileSecurityService{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mss-operator",
					},
				},
			},
			want: "mss-operator",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRouteName(tt.args.m); got != tt.want {
				t.Errorf("GetRouteName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConfigMapName(t *testing.T) {
	type args struct {
		m *mobilesecurityservicev1alpha1.MobileSecurityService
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should use ConfigMapName",
			args: args{
				m: &mobilesecurityservicev1alpha1.MobileSecurityService{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mss-operator",
					},
					Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceSpec{
						ConfigMapName: "mss-config",
					},
				},
			},
			want: "mss-config",
		},
		{
			name: "should use Instance name",
			args: args{
				m: &mobilesecurityservicev1alpha1.MobileSecurityService{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mss-operator",
					},
				},
			},
			want: "mss-operator",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConfigMapName(tt.args.m); got != tt.want {
				t.Errorf("GetConfigMapName() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			os.Unsetenv(APP_NAMESPACE_ENV_VAR)

			if tt.fields.envVar != "" {
				os.Setenv(APP_NAMESPACE_ENV_VAR, tt.fields.envVar)
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
			name: "Should return an error if no namespace is set",
			args: args{
				namespace: "apps-namespace",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// clear old env var
			os.Unsetenv(APP_NAMESPACE_ENV_VAR)

			if tt.fields.envVar != "" {
				os.Setenv(APP_NAMESPACE_ENV_VAR, tt.fields.envVar)
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
		skipCheck bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "should skip the check and return true",
			args: args{
				namespace: "mobile-security-service-operator",
				skipCheck: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "should return false as cannot read from the file system",
			args: args{
				namespace: "mobile-security-service-operator",
				skipCheck: false,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsValidOperatorNamespace(tt.args.namespace, tt.args.skipCheck)
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
