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
TAG= 0.1.0

IMAGE_REGISTRY=quay.io
REGISTRY_ORG=aerogear
REGISTRY_REPO=mobile-security-service-operator
IMAGE_MASTER_TAG=$(IMAGE_REGISTRY)/$(REGISTRY_ORG)/$(REGISTRY_REPO):master
IMAGE_RELEASE_TAG=$(IMAGE_REGISTRY)/$(REGISTRY_ORG)/$(REGISTRY_REPO):$(CIRCLE_TAG)
NAMESPACE=mobile-security-service
APP_NAMESPACES=mobile-security-service-apps

# This follows the output format for goreleaser
BINARY_LINUX_64 = ./dist/linux_amd64/$(BINARY)

LDFLAGS=-ldflags "-w -s -X main.Version=${TAG}"

#########################################
# Local Development                     #
#########################################

.PHONY: setup-githooks
setup-githooks:
	@echo Setting up Git hooks:
	ln -sf $$PWD/.githooks/* $$PWD/.git/hooks/

.PHONY: setup
setup: setup-githooks
	dep ensure

.PHONY: debug-setup
debug-setup:
	@echo Exporting env vars to run operator locally:
	- . ./scripts/export_local_envvars.sh
	@echo Installing ...
	- make install

.PHONY: vet
vet:
	@echo go vet
	go vet $$(go list ./... | grep -v /vendor/)

.PHONY: fmt
fmt:
	@echo go fmt
	go fmt $$(go list ./... | grep -v /vendor/)

.PHONY: refresh-operator-image
refresh-operator-image:
	@echo INFO: Deleting and re-applying the operator ...
	- oc project ${NAMESPACE}
	- kubectl delete -f deploy/operator.yaml
	- kubectl create -f deploy/operator.yaml

#########################################
# CI                                    #
#########################################

.PHONY: test
test:
	@echo Running tests:
	GOCACHE=off go test -cover $(TEST_PKGS)

.PHONY: test-integration-cover
test-integration-cover:
	echo "mode: count" > coverage-all.out
	GOCACHE=off $(foreach pkg,$(PACKAGES),\
		go test -failfast -tags=integration -coverprofile=coverage.out -covermode=count $(addprefix $(PKG)/,$(pkg)) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)

.PHONY: build-linux
build-linux:
	env GOOS=linux GOARCH=amd64 go build $(APP_FILE)

.PHONY: build-master
build-master:
	@echo Building operator with the tag: $(IMAGE_MASTER_TAG)
	operator-sdk build $(IMAGE_MASTER_TAG)

.PHONY: build-release
build-release:
	@echo Building operator with the tag: $(IMAGE_RELEASE_TAG)
	operator-sdk build $(IMAGE_RELEASE_TAG)

.PHONY: push-master
push-master:
	@echo Pushing operator with tag $(IMAGE_MASTER_TAG) to $(IMAGE_REGISTRY)
	@docker login --username $(QUAY_USERNAME) --password $(QUAY_PASSWORD) quay.io
	docker push $(IMAGE_MASTER_TAG)

.PHONY: push-release
push-release:
	@echo Pushing operator with tag $(IMAGE_RELEASE_TAG) to $(IMAGE_REGISTRY)
	@docker login --username $(QUAY_USERNAME) --password $(QUAY_PASSWORD) quay.io
	docker push $(IMAGE_RELEASE_TAG)

#########################################
# Operator                              #
#########################################

.PHONY: apply-app-example
apply-app-example:
	kubectl apply -f deploy/crds/examples/mobile-security-service_v1alpha1_mobilesecurityserviceapp_cr.yaml

.PHONY: delete-app-example
delete-app-example:
	kubectl delete -f deploy/crds/examples/mobile-security-service_v1alpha1_mobilesecurityserviceapp_cr.yaml

.PHONY: run-local
run-local:
	@echo INFO: Exporting env vars to run operator locally ...
	- . ./scripts/export_local_envvars.sh
	@echo INFO: Installing ...
	- make install
	@echo INFO: Starting ...
	- operator-sdk up local

.PHONY: install
install:
	@echo INFO: Creating namespace ${NAMESPACE} ...
	- oc new-project ${NAMESPACE}
	@echo INFO: Applying files into ${NAMESPACE} ...
	- kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	- kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_crd.yaml
	- kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceapp_crd.yaml
	- kubectl apply -f deploy/cluster_role.yaml
	- kubectl apply -f deploy/cluster_role_binding.yaml
	- kubectl apply -f deploy/service_account.yaml
	- kubectl apply -f deploy/operator.yaml
	@echo INFO: Applying Mobile Security Service and Database into ${NAMESPACE} ...:
	- kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_cr.yaml
	- kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_cr.yaml

.PHONY: uninstall
uninstall:
	@echo INFO: Uninstalling all from ${NAMESPACE} ...
	- oc project ${NAMESPACE}
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceapp_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_crd.yaml
	- kubectl delete -f deploy/cluster_role.yaml
	- kubectl delete -f deploy/cluster_role_binding.yaml
	- kubectl delete -f deploy/service_account.yaml
	- kubectl delete -f deploy/operator.yaml
	- make uninstall-monitoring
	- kubectl delete namespace ${NAMESPACE}

.PHONY: install-monitoring
install-monitoring:
	@echo INFO: Installing service monitor in ${NAMESPACE} ...
	- oc project ${NAMESPACE}
	- kubectl label namespace ${NAMESPACE} monitoring-key=middleware
	- kubectl create -f deploy/monitor/service_monitor.yaml
	- kubectl create -f deploy/monitor/operator_service.yaml
	- kubectl create -f deploy/monitor/prometheus-rule.yaml
	- kubectl create -f deploy/monitor/grafana-dashboard.yaml

.PHONY: uninstall-monitoring
uninstall-monitoring:
	@echo INFO: Uninstalling monitor service from ${NAMESPACE} ...
	- oc project ${NAMESPACE}
	- kubectl delete -f deploy/monitor/service_monitor.yaml
	- kubectl delete -f deploy/monitor/operator_service.yaml
	- kubectl delete -f deploy/monitor/prometheus-rule.yaml
	- kubectl delete -f deploy/monitor/grafana-dashboard.yaml
