APP_NAME = mobile-security-service-operator
ORG_NAME = aerogear
PKG = github.com/$(ORG_NAME)/$(APP_NAME)
TOP_SRC_DIRS = pkg
PACKAGES ?= $(shell sh -c "find $(TOP_SRC_DIRS) -name \\*_test.go \
              -exec dirname {} \\; | sort | uniq")
TEST_PKGS = $(addprefix $(PKG)/,$(PACKAGES))
APP_FILE=./cmd/manager/main.go
BIN_DIR := $(GOPATH)/bin
BINARY ?= mobile-security-service-operator
IMAGE_REGISTRY=quay.io
REGISTRY_ORG=aerogear
REGISTRY_REPO=mobile-security-service-operator
IMAGE_LATEST_TAG=$(IMAGE_REGISTRY)/$(REGISTRY_ORG)/$(REGISTRY_REPO):latest
IMAGE_MASTER_TAG=$(IMAGE_REGISTRY)/$(REGISTRY_ORG)/$(REGISTRY_REPO):master
IMAGE_RELEASE_TAG=$(IMAGE_REGISTRY)/$(REGISTRY_ORG)/$(REGISTRY_REPO):$(CIRCLE_TAG)
NAMESPACE=mobile-security-service
APP_NAMESPACES=mobile-security-service-apps
CODE_COMPILE_OUTPUT = build/_output/bin/mobile-security-service-operator
TEST_COMPILE_OUTPUT = build/_output/bin/mobile-security-service-operator-test
# This follows the output format for goreleaser
BINARY_LINUX_64 = ./dist/linux_amd64/$(BINARY)

LDFLAGS=-ldflags "-w -s -X main.Version=${TAG}"

##############################
# INSTALL/UNINSTALL          #
##############################

.PHONY: install
install:
	@echo ....... Creating namespace ....... 
	- oc new-project ${NAMESPACE}
	@echo ....... Applying Mobile Security Service CRDS and Operator .......
	- kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	- kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_crd.yaml
	- kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceapp_crd.yaml
	- kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicebackup_crd.yaml
	@echo ....... Applying Rules and Service Account .......
	- kubectl apply -f deploy/role.yaml
	- kubectl apply -f deploy/role_binding.yaml
	- kubectl apply -f deploy/service_account.yaml
	@echo ....... Applying Mobile Security Service Operator .......
	- kubectl apply -f deploy/operator.yaml
	@echo ....... Creating the Mobile Security Service and Database .......
	- kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_cr.yaml
	- kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_cr.yaml
	@echo ....... Creating namespace ${APP_NAMESPACES} .......
	- oc new-project ${APP_NAMESPACES}

.PHONY: uninstall
uninstall:
	@echo ....... Uninstalling .......
	- oc project ${NAMESPACE}
	@echo ....... Deleting the Mobile Security Service and Database .......
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_cr.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_cr.yaml
	@echo ....... Deleting CRDs.......
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceapp_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicebackup_crd.yaml
	@echo ....... Deleting Rules and Service Account .......
	- kubectl delete -f deploy/role.yaml
	- kubectl delete -f deploy/role_binding.yaml
	- kubectl delete -f deploy/service_account.yaml
	@echo ....... Deleting Mobile Security Service Operator .......
	- kubectl delete -f deploy/operator.yaml
	@echo ....... Deleting namespace ${NAMESPACE}.......
	- kubectl delete namespace ${NAMESPACE}
	@echo ....... Delete namespace ${APP_NAMESPACES} .......
	- kubectl delete namespace ${APP_NAMESPACES}

.PHONY: refresh-operator-image
refresh-operator-image:
	@echo ....... Deleting and applying the operator in order to refresh the image when a tag is not changed \(development use\).......
	- oc project ${NAMESPACE}
	- kubectl delete -f deploy/operator.yaml
	- kubectl create -f deploy/operator.yaml

.PHONY: example-app/apply
example-app/apply:
	@echo ....... Applying the MobileSecurityServiceApp example in the current namespace  ......
	@echo ....... An APP CR can only be applied in the namespaces configured in the operator\'s EnvVar APP_NAMESPACES.
	- kubectl apply -f deploy/crds/examples/mobile-security-service_v1alpha1_mobilesecurityserviceapp_cr.yaml
	
.PHONY: example-app/delete
example-app/delete:
	@echo ....... Deleting the MobileSecurityServiceApp example from the current app namespace  ......
	- kubectl delete -f deploy/crds/examples/mobile-security-service_v1alpha1_mobilesecurityserviceapp_cr.yaml

.PHONY: monitoring/install
monitoring/install:
	@echo Installing service monitor in ${NAMESPACE} :
	- oc project ${NAMESPACE}
	- kubectl label namespace ${NAMESPACE} monitoring-key=middleware
	- kubectl create -f deploy/monitor/mss_service_monitor.yaml
	- kubectl create -f deploy/monitor/prometheus_rule.yaml
	- kubectl create -f deploy/monitor/mss_prometheus_rule.yaml
	- kubectl create -f deploy/monitor/grafana_dashboard.yaml
	- kubectl create -f deploy/monitor/mss_grafana_dashboard.yaml

.PHONY: monitoring/uninstall
monitoring/uninstall:
	@echo Uninstalling monitor service from ${NAMESPACE} :
	- oc project ${NAMESPACE}
	- kubectl delete -f deploy/monitor/mss_service_monitor.yaml
	- kubectl delete -f deploy/monitor/prometheus_rule.yaml
	- kubectl delete -f deploy/monitor/mss_prometheus_rule.yaml
	- kubectl delete -f deploy/monitor/grafana_dashboard.yaml
	- kubectl delete -f deploy/monitor/mss_grafana_dashboard.yaml

.PHONY: backup/install
backup/install:
	@echo Installing backup service in ${NAMESPACE} :
	- oc project ${NAMESPACE}
	- kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicebackup_cr.yaml

.PHONY: backup/uninstall
backup/uninstall:
	@echo Uninstalling backup service from ${NAMESPACE} :
	- oc project ${NAMESPACE}
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicebackup_cr.yaml

##############################
# CI                         #
##############################

.PHONY: code/build/linux
code/build/linux:
	env GOOS=linux GOARCH=amd64 go build $(APP_FILE)

.PHONY: code/compile
code/compile: code/gen
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o=$(CODE_COMPILE_OUTPUT) ./cmd/manager/main.go
	
.PHONY: image/build/master
image/build/master:
	@echo Building operator with the tag: $(IMAGE_MASTER_TAG)
	operator-sdk build $(IMAGE_MASTER_TAG)

.PHONY: image/build/release
image/build/release:
	@echo Building operator with the tag: $(IMAGE_RELEASE_TAG)
	operator-sdk build $(IMAGE_RELEASE_TAG)
	operator-sdk build $(IMAGE_LATEST_TAG)

.PHONY: image/push/master
image/push/master:
	@echo Pushing operator with tag $(IMAGE_MASTER_TAG) to $(IMAGE_REGISTRY)
	@docker login --username $(QUAY_USERNAME) --password $(QUAY_PASSWORD) quay.io
	docker push $(IMAGE_MASTER_TAG)

.PHONY: image/push/release
image/push/release:
	@echo Pushing operator with tag $(IMAGE_RELEASE_TAG) to $(IMAGE_REGISTRY)
	@docker login --username $(QUAY_USERNAME) --password $(QUAY_PASSWORD) quay.io
	docker push $(IMAGE_RELEASE_TAG)
	@echo Pushing operator with tag $(IMAGE_LATEST_TAG) to $(IMAGE_REGISTRY)
	docker push $(IMAGE_LATEST_TAG)


##############################
# Local Development          #
##############################

.PHONY: setup/debug
setup/debug:
	@echo Exporting env vars to run operator locally:
	- . ./scripts/export_local_envvars.sh
	@echo Installing ...
	- make install

.PHONY: setup/githooks
setup/githooks:
	@echo Installing errcheck dependence:
	go get -u github.com/kisielk/errcheck
	@echo Setting up Git hooks:
	ln -sf $$PWD/.githooks/* $$PWD/.git/hooks/

.PHONY: setup
setup: setup/githooks
	dep ensure
	make code/gen

.PHONY: code/run/local
code/run/local:
	@echo Exporting env vars to run operator locally:
	- . ./scripts/export_local_envvars.sh
	@echo  ....... Installing ...
	- make install
	@echo Starting ...
	- operator-sdk up local

.PHONY: code/vet
code/vet:
	@echo go vet
	go vet $$(go list ./... | grep -v /vendor/)

.PHONY: code/fmt
code/fmt:
	@echo go fmt
	go fmt $$(go list ./... | grep -v /vendor/)

.PHONY: code/gen
code/gen:
	operator-sdk generate k8s
	operator-sdk generate openapi

##############################
# Tests                      #
##############################

.PHONY: test/run
test/run:
	@echo Running tests:
	GOCACHE=off go test -cover $(TEST_PKGS)
.PHONY: test/compile
test/compile:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -c -o=$(TEST_COMPILE_OUTPUT) ./test/e2e/...
.PHONY: test/integration-cover
test/integration-cover:
	echo "mode: count" > coverage-all.out
	GOCACHE=off $(foreach pkg,$(PACKAGES),\
		go test -failfast -tags=integration -coverprofile=coverage.out -covermode=count $(addprefix $(PKG)/,$(pkg)) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)
