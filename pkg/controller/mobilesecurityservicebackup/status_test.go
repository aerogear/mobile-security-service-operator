package mobilesecurityservicebackup

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestMobileSecurityServiceBackup_updateBackupStatus(t *testing.T) {
	type fields struct {
		objs   []runtime.Object
		scheme *runtime.Scheme
	}
	type args struct {
		cronJobStatus *v1beta1.CronJob
		request       reconcile.Request
		dbSecret      *corev1.Secret
		awsSecret     *corev1.Secret
		dbPod         *corev1.Pod
		dbService     *corev1.Service
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "should return an error when no name found",
			fields: fields{
				objs:   []runtime.Object{&bkpInstance},
				scheme: scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
				cronJobStatus: &v1beta1.CronJob{},
				dbSecret:      &corev1.Secret{},
				awsSecret:     &corev1.Secret{},
				dbPod:         &corev1.Pod{},
				dbService:     &corev1.Service{},
			},

			wantErr: true,
		},
		{
			name: "should update status without enc secret",
			fields: fields{
				objs:   []runtime.Object{&bkpInstance},
				scheme: scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
				cronJobStatus: &v1beta1.CronJob{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
				dbSecret: &corev1.Secret{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
				awsSecret: &corev1.Secret{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
				dbPod: &corev1.Pod{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
				dbService: &corev1.Service{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
			},
			wantErr: false,
		},
		{
			name: "should return error when not found secret by name",
			fields: fields{
				objs:   []runtime.Object{&bkpInstanceWithSecretNames},
				scheme: scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      bkpInstance.Name,
						Namespace: bkpInstance.Namespace,
					},
				},
				cronJobStatus: &v1beta1.CronJob{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
				dbSecret: &corev1.Secret{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
				dbPod: &corev1.Pod{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
				dbService: &corev1.Service{ObjectMeta: v1.ObjectMeta{
					Name: "test",
				}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs, t)

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			if err := r.updateBackupStatus(reqLogger, tt.args.cronJobStatus, tt.args.dbSecret, tt.args.awsSecret, tt.args.dbPod, tt.args.dbService, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("MobileSecurityServiceBackup.updateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMobileSecurityServiceBackup_updateStatusWithInvalidNamespace(t *testing.T) {
	type args struct {
		instance *mobilesecurityservicev1alpha1.MobileSecurityServiceBackup
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should return without an error when updating status",
			args: args{
				instance: &bkpInstance,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{
				tt.args.instance,
			}

			r := buildReconcileWithFakeClientWithMocks(objs, t)

			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			// mock request to simulate Reconcile() being called on an event for a watched resource
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      bkpInstance.Name,
					Namespace: bkpInstance.Namespace,
				},
			}

			if err := r.updateStatusWithInvalidNamespace(reqLogger, req); (err != nil) != tt.wantErr {
				t.Errorf("MobileSecurityServiceBackup.updateStatusWithInvalidNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
