package mobilesecurityservice

import (
	errs "errors"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcileMobileSecurityService_update(t *testing.T) {
	type fields struct {
		instance *v1alpha1.MobileSecurityService
		scheme   *runtime.Scheme
	}
	tests := []struct {
		name    string
		fields  fields
		want    reconcile.Result
		wantErr bool
	}{
		{
			name: "should requeue",
			fields: fields{
				instance: &mobilesecurityservicev1alpha1.MobileSecurityService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "mobile-security-service-app",
						Namespace: "mobile-security-service-operator",
					},
					Spec: v1alpha1.MobileSecurityServiceSpec{
						Size: 1,
					},
				},
				scheme: scheme.Scheme,
			},
			want:    reconcile.Result{Requeue: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.fields.instance}

			tt.fields.scheme.AddKnownTypes(mobilesecurityservicev1alpha1.SchemeGroupVersion, tt.fields.instance)

			cl := fake.NewFakeClient(objs...)

			r := &ReconcileMobileSecurityService{cl, tt.fields.scheme}

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      tt.fields.instance.Name,
					Namespace: tt.fields.instance.Namespace,
				},
			}

			reqLogger := log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)

			got, err := r.update(objs[0], reqLogger)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityService.update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReconcileMobileSecurityService.update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReconcileMobileSecurityService_create(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
	}
	type args struct {
		instance *mobilesecurityservicev1alpha1.MobileSecurityService
		kind     string
		err      error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    reconcile.Result
		wantErr bool
	}{
		{
			name: "should return an error",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance: &mobilesecurityservicev1alpha1.MobileSecurityService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "mobile-security-service-app",
						Namespace: "mobile-security-service-operator",
					},
					Spec: v1alpha1.MobileSecurityServiceSpec{
						Size:          1,
						MemoryLimit:   "512Mi",
						MemoryRequest: "512Mi",
					},
				},
				kind: DEEPLOYMENT,
				err:  errors.NewInternalError(errs.New("Internal Server Error")),
			},
			want:    reconcile.Result{},
			wantErr: true,
		},
		{
			name: "should create and return a new deployment",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance: &mobilesecurityservicev1alpha1.MobileSecurityService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "mobile-security-service-app",
						Namespace: "mobile-security-service-operator",
					},
					Spec: v1alpha1.MobileSecurityServiceSpec{
						Size:          1,
						MemoryLimit:   "512Mi",
						MemoryRequest: "512Mi",
					},
				},
				kind: DEEPLOYMENT,
				err:  errors.NewNotFound(schema.GroupResource{Group: "api/v1", Resource: "ResourceName"}, "Not Found"),
			},
			want:    reconcile.Result{Requeue: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.args.instance}

			tt.fields.scheme.AddKnownTypes(mobilesecurityservicev1alpha1.SchemeGroupVersion, tt.args.instance)

			cl := fake.NewFakeClient(objs...)

			r := &ReconcileMobileSecurityService{
				client: cl,
				scheme: tt.fields.scheme,
			}

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      tt.args.instance.Name,
					Namespace: tt.args.instance.Namespace,
				},
			}

			reqLogger := log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)

			got, err := r.create(tt.args.instance, reqLogger, tt.args.kind, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityService.create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReconcileMobileSecurityService.create() = %v, want %v", got, tt.want)
			}
		})
	}
}
