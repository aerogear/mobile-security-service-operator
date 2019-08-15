package mobilesecurityservicebackup

import (
	"github.com/aerogear/mobile-security-service-operator/pkg/apis/mobile-security-service/v1alpha1"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobile-security-service/v1alpha1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestReconcileMobileSecurityServiceBackup_update(t *testing.T) {
	type fields struct {
		createdInstance  *v1alpha1.MobileSecurityServiceBackup
		instanceToUpdate *v1alpha1.MobileSecurityServiceBackup
		scheme           *runtime.Scheme
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "should successfully update the instance",
			fields: fields{
				createdInstance:  &bkpInstance,
				instanceToUpdate: &bkpInstance,
			},
			wantErr: false,
		},
		{
			name: "should give an error when the namespace is not found",
			fields: fields{
				createdInstance:  &bkpInstance,
				instanceToUpdate: &bkpInstanceNonDefaultNamespace,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.fields.createdInstance}

			r := buildReconcileWithFakeClientWithMocks(objs, t)

			reqLogger := log.WithValues("Request.Namespace", tt.fields.instanceToUpdate.Namespace, "Request.Name", tt.fields.createdInstance.Name)

			err := r.update(tt.fields.instanceToUpdate, reqLogger)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityServiceBackup.update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReconcileMobileSecurityServiceBackup_create(t *testing.T) {
	objs := []runtime.Object{&bkpInstance}
	r := buildReconcileWithFakeClientWithMocks(objs, t)
	dataDBSecret, _ := r.buildDBSecretData(&bkpInstance)
	awsDataSecret := buildAwsSecretData(&bkpInstance)
	type fields struct {
		scheme *runtime.Scheme
	}
	type args struct {
		instance   *mobilesecurityservicev1alpha1.MobileSecurityServiceBackup
		kind       string
		secretData map[string][]byte
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "should create and return a new CronJob",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       CronJob,
				secretData: nil,
			},
			wantErr: false,
		},
		{
			name: "should create and return a new DB Secret",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       DBSecret,
				secretData: dataDBSecret,
			},
			wantErr: false,
		},
		{
			name: "should create and return a new Aws CronJob",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       AwsSecret,
				secretData: awsDataSecret,
			},
			wantErr: false,
		},
		{
			name: "should fail to create an unknown type",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance: &bkpInstance,
				kind:     "UNKNOWN",
			},
			wantErr:   false,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ReconcileMobileSecurityServiceBackup.create() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			err := r.create(tt.args.instance, tt.args.kind, reqLogger)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityServiceBackup.create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReconcileMobileSecurityServiceBackup_buildFactory(t *testing.T) {
	objs := []runtime.Object{&bkpInstance}
	r := buildReconcileWithFakeClientWithMocks(objs, t)
	dataDBSecret, _ := r.buildDBSecretData(&bkpInstance)
	awsDataSecret := buildAwsSecretData(&bkpInstance)

	type fields struct {
		scheme *runtime.Scheme
	}
	type args struct {
		instance   *mobilesecurityservicev1alpha1.MobileSecurityServiceBackup
		kind       string
		secretData map[string][]byte
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      reflect.Type
		wantPanic bool
	}{
		{
			name: "should create a CronJob",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       CronJob,
				secretData: nil,
			},
			want: reflect.TypeOf(&v1beta1.CronJob{}),
		},
		{
			name: "should create a DB Secret",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       DBSecret,
				secretData: dataDBSecret,
			},
			want: reflect.TypeOf(&corev1.Secret{}),
		},
		{
			name: "should create a Aws Secret",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:   &bkpInstance,
				kind:       AwsSecret,
				secretData: awsDataSecret,
			},
			want: reflect.TypeOf(&corev1.Secret{}),
		},
		{
			name: "Should panic when trying to create unrecognized object type",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance: &bkpInstance,
				kind:     "UNDEFINED",
			},
			wantPanic: true,
			want:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			// testing if the panic() function is executed
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("MobileSecurityServiceBackup.buildFactory() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			got, _ := r.buildFactory(tt.args.instance, tt.args.kind, reqLogger)

			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("MobileSecurityServiceBackup.buildFactory() = %v, want %v", gotType, tt.want)
			}
		})
	}
}

func TestReconcileMobileSecurityServiceBackup_Reconcile_NotFound(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&bkpInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs, t)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      bkpInstance.Name,
			Namespace: "unknown",
		},
	}

	res, err := r.Reconcile(req)
	if err == nil {
		t.Error("should fail since the instance do not exist in the <unknown> nammespace")
	}

	if res.Requeue {
		t.Fatalf("did not expected reconcile to requeue.")
	}
}
