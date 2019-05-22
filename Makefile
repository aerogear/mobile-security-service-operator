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
DEV= dev
IMAGE_REGISTRY=quay.io
REGISTRY_ORG=aerogear
REGISTRY_REPO=mobile-security-service-operator
IMAGE_DEV_TAG=$(IMAGE_REGISTRY)/$(REGISTRY_ORG)/$(REGISTRY_REPO):$(TAG)-$(DEV)
IMAGE_LATEST_TAG=$(IMAGE_REGISTRY)/$(REGISTRY_ORG)/$(REGISTRY_REPO):latest
IMAGE_MASTER_TAG=$(IMAGE_REGISTRY)/$(REGISTRY_ORG)/$(REGISTRY_REPO):master
IMAGE_RELEASE_TAG=$(IMAGE_REGISTRY)/$(REGISTRY_ORG)/$(REGISTRY_REPO):$(CIRCLE_TAG)
NAMESPACE=mobile-security-service
APP_NAMESPACES=mobile-security-service-apps

# This follows the output format for goreleaser
BINARY_LINUX_64 = ./dist/linux_amd64/$(BINARY)

LDFLAGS=-ldflags "-w -s -X main.Version=${TAG}"


.PHONY: setup-githooks
setup-githooks:
	@echo Installing errcheck dependence:
	go get -u github.com/kisielk/errcheck
	@echo Setting up Git hooks:
	ln -sf $$PWD/.githooks/* $$PWD/.git/hooks/

.PHONY: setup
setup: setup-githooks
	dep ensure

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

.PHONY: build_linux
build_linux:
	env GOOS=linux GOARCH=amd64 go build $(APP_FILE)

.PHONY: create-app-ns
create-app-ns:
	@echo Creating the namespace ${APP_NAMESPACES}:
	oc new-project ${APP_NAMESPACES}

.PHONY: create-app
create-app:
	kubectl create -f deploy/crds/examples/mobile-security-service_v1alpha1_mobilesecurityserviceapp_cr.yaml

.PHONY: delete-app
delete-app:
	kubectl delete -f deploy/crds/examples/mobile-security-service_v1alpha1_mobilesecurityserviceapp_cr.yaml

.PHONY: run-local
run-local:
	@echo Exporting env vars to run operator locally:
	- . ./scripts/export_local_envvars.sh
	@echo Installing ...
	- make create-all
	@echo Starting ...
	- operator-sdk up local

.PHONY: create-all
create-all:
	@echo Create Mobile Security Service Operator and Service in the namespace ${NAMESPACE}:
	make create-oper
	make create-service-and-db

.PHONY: delete-all
delete-all:
	@echo Delete Mobile Security Service Operator, Service and namespace ${NAMESPACE}:
	- kubectl delete -f deploy/service_monitor.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceapp_crd.yaml
	- kubectl delete -f deploy/cluster_role.yaml
	- kubectl delete -f deploy/cluster_role_binding.yaml
	- kubectl delete -f deploy/service_account.yaml
	- kubectl delete -f deploy/operator.yaml
	- kubectl delete namespace ${NAMESPACE}

.PHONY: create-oper
create-oper:
	@echo Create Mobile Security Service Operator:
	- oc new-project ${NAMESPACE}
	- kubectl label namespace ${NAMESPACE} monitoring-key=middleware
	- kubectl create -f deploy/service_monitor.yaml
	- kubectl create -f deploy/operator_service.yaml -n ${NAMESPACE}
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_crd.yaml
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceapp_crd.yaml
	- kubectl create -f deploy/cluster_role.yaml
	- kubectl create -f deploy/cluster_role_binding.yaml
	- kubectl create -f deploy/service_account.yaml
	- kubectl create -f deploy/operator.yaml

.PHONY: delete-oper
delete-oper:
	@echo Deleting Mobile Security Service Operator and namespace:
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceapp_crd.yaml
	- kubectl delete -f deploy/cluster_role.yaml
	- kubectl delete -f deploy/cluster_role_binding.yaml
	- kubectl delete -f deploy/service_account.yaml
	- kubectl delete -f deploy/operator.yaml
	- kubectl delete namespace ${NAMESPACE}

.PHONY: create-service-and-db
create-service-and-db:
	@echo Creating Mobile Security Service and Database into project:
	kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_cr.yaml
	kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_cr.yaml

.PHONY: create-service-only
create-service-only:
	@echo Creating Mobile Security Service App only:
	kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_cr.yaml

.PHONY: create-db-only
create-db-only:
	@echo Creating Mobile Security Service Database only:
	kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_cr.yaml

.PHONY: delete-service-and-db
delete-service-and-db:
	@echo Deleting Mobile Security Service and Database:
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_cr.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_cr.yaml

.PHONY: delete-service-only
delete-service-only:
	@echo Deleting Mobile Security Service App only:
	kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_cr.yaml

.PHONY: delete-db-only
delete-db-only:
	@echo Deleting Mobile Security Service Database only:
	kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_cr.yaml

.PHONY: build-dev
build-dev:
	@echo Building operator with the tag: $(IMAGE_DEV_TAG)
	operator-sdk build $(IMAGE_DEV_TAG)

.PHONY: build-master
build-master:
	@echo Building operator with the tag: $(IMAGE_MASTER_TAG)
	operator-sdk build $(IMAGE_MASTER_TAG)

.PHONY: build-release
build-release:
	@echo Building operator with the tag: $(IMAGE_RELEASE_TAG)
	operator-sdk build $(IMAGE_RELEASE_TAG)

.PHONY: build-latest
build-latest:
	@echo Building operator with the tag: $(IMAGE_LATEST_TAG)
	operator-sdk build $(IMAGE_LATEST_TAG)
	
.PHONY: push-dev
push-dev:
	@echo Pushing operator with tag $(IMAGE_DEV_TAG) to $(IMAGE_REGISTRY)
	@docker login $(IMAGE_REGISTRY)
	docker push $(IMAGE_DEV_TAG)

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

.PHONY: push-latest
push-latest:
	@echo Pushing operator with tag $(IMAGE_LATEST_TAG) to $(IMAGE_REGISTRY)
	@docker login --username $(QUAY_USERNAME) --password $(QUAY_PASSWORD) quay.io
	docker push $(IMAGE_LATEST_TAG)

.PHONY: debug-setup
debug-setup:
	@echo Exporting WATCH_NAMESPACE=default:
	- export WATCH_NAMESPACE=default
	@echo Create Namespace:
	- oc new-project ${NAMESPACE}
	@echo Installing the CRD:
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_crd.yaml
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceapp_crd.yaml

.PHONY: vet
vet:
	@echo go vet
	go vet $$(go list ./... | grep -v /vendor/)

.PHONY: fmt
fmt:
	@echo go fmt
	go fmt $$(go list ./... | grep -v /vendor/)
