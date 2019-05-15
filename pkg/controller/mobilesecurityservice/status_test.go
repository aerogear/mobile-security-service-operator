package mobilesecurityservice

// func TestReconcileMobileSecurityService_updateStatus(t *testing.T) {
// 	type fields struct {
// 		client client.Client
// 		scheme *runtime.Scheme
// 	}
// 	type args struct {
// 		reqLogger        logr.Logger
// 		configMapStatus  *corev1.ConfigMap
// 		deploymentStatus *v1beta1.Deployment
// 		serviceStatus    *corev1.Service
// 		routeStatus      *routev1.Route
// 		request          reconcile.Request
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &ReconcileMobileSecurityService{
// 				client: tt.fields.client,
// 				scheme: tt.fields.scheme,
// 			}
// 			if err := r.updateStatus(tt.args.reqLogger, tt.args.configMapStatus, tt.args.deploymentStatus, tt.args.serviceStatus, tt.args.routeStatus, tt.args.request); (err != nil) != tt.wantErr {
// 				t.Errorf("ReconcileMobileSecurityService.updateStatus() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestReconcileMobileSecurityService_updateConfigMapStatus(t *testing.T) {
// 	type fields struct {
// 		instance *mobilesecurityservicev1alpha1.MobileSecurityService
// 		scheme   *runtime.Scheme
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		want    *corev1.ConfigMap
// 		wantErr bool
// 	}{
// 		{
// 			name: "should update the ConfigMap status",
// 			fields: fields{
// 				instance: &instance,
// 				scheme:   scheme.Scheme,
// 			},
// 			want: &corev1.ConfigMap{
// 				TypeMeta: metav1.TypeMeta{
// 					APIVersion: "v1",
// 					Kind:       "ConfigMap",
// 				},
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name:      utils.GetConfigMapName(&instance),
// 					Namespace: instance.Namespace,
// 					Labels:    getAppLabels(instance.Name),
// 				},
// 				Data: getAppEnvVarsMap(&instance),
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			objs := []runtime.Object{tt.fields.instance}

// 			r, _ := getReconcile(objs)

// 			req := reconcile.Request{
// 				NamespacedName: types.NamespacedName{
// 					Name:      tt.fields.instance.Name,
// 					Namespace: tt.fields.instance.Namespace,
// 				},
// 			}

// 			reqLogger := log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)

// 			// call reconcile to create the MobileSecurityService instance
// 			_, err := r.Reconcile(req)
// 			if err != nil {
// 				t.Fatalf("reconcile: (%v)", err)
// 			}

// 			got, err := r.updateConfigMapStatus(reqLogger, req)

// 			// if (err != nil) != tt.wantErr {
// 			// 	t.Errorf("ReconcileMobileSecurityService.updateConfigMapStatus() error = %v, wantErr %v", err, tt.wantErr)
// 			// 	return
// 			// }
// 			// if !reflect.DeepEqual(got, tt.want) {
// 			// 	t.Errorf("ReconcileMobileSecurityService.updateConfigMapStatus() = %v, want %v", got, tt.want)
// 			// }
// 		})
// 	}
// }
