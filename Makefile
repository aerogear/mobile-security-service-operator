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
DOCKER-ORG=aerogear
DOCKER-REPO=mobile-security-service-operator

# This follows the output format for goreleaser
BINARY_LINUX_64 = ./dist/linux_amd64/$(BINARY)

LDFLAGS=-ldflags "-w -s -X main.Version=${TAG}"

.PHONY: setup
setup:
	dep ensure

.PHONY: test
test:
	@echo Running tests:
	@if [ -z $(TEST_PKGS) ]; then \
	  echo No test files found to test; \
	else \
	  GOCACHE=off go test -cover $(TEST_PKGS); \
	fi

.PHONY: build_linux
build_linux:
	env GOOS=linux GOARCH=amd64 go build $(APP_FILE)

.PHONY: create-app
create-app:
	- kubectl delete -f deploy/crds/examples/mobile-security-service_v1alpha1_mobilesecurityserviceunbind_cr.yaml
	- kubectl create -f deploy/crds/examples/mobile-security-service_v1alpha1_mobilesecurityserviceapp_cr.yaml

.PHONY: unbind
unbind:
	- kubectl delete -f deploy/crds/examples/mobile-security-service_v1alpha1_mobilesecurityserviceapp_cr.yaml
	- kubectl apply -f deploy/crds/examples/mobile-security-service_v1alpha1_mobilesecurityserviceunbind_cr.yaml

.PHONY: run-local
run-local:
	@echo Installing the operator in a cluster and run it locally:
	- export OPERATOR_NAME=mobile-security-service-operator
	- make create-all
	- operator-sdk up local --namespace=mobile-security-service-operator

.PHONY: create-all
create-all:
	@echo Create Mobile Security Service Operator and Service in the namespace "mobile-security-service-operator":
	make create-oper
	make create-service

.PHONY: delete-all
delete-all:
	@echo Delete Mobile Security Service Operator, Service and namespace "mobile-security-service-operator":
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_cr.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_cr.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceapp_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceunbind_crd.yaml
	- kubectl delete -f deploy/cluster_role.yaml
	- kubectl delete -f deploy/cluster_role_binding.yaml
	- kubectl delete -f deploy/service_account.yaml
	- kubectl delete -f deploy/operator.yaml
	- kubectl delete -f deploy/role.yaml
	- kubectl delete -f deploy/role_binding.yaml
	- kubectl delete namespace mobile-security-service-operator

.PHONY: create-oper
create-oper:
	@echo Create Mobile Security Service Operator:
	- kubectl create namespace mobile-security-service-operator
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_crd.yaml
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceapp_crd.yaml
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceunbind_crd.yaml
	- kubectl create -f deploy/cluster_role.yaml
	- kubectl create -f deploy/cluster_role_binding.yaml
	- kubectl create -f deploy/role.yaml
	- kubectl create -f deploy/role_binding.yaml
	- kubectl create -f deploy/service_account.yaml
	- kubectl create -f deploy/operator.yaml

.PHONY: delete-oper
delete-oper:
	@echo Deleting Mobile Security Service Operator and namespace:
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceapp_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceunbind_crd.yaml
	- kubectl delete -f deploy/cluster_role.yaml
	- kubectl delete -f deploy/cluster_role_binding.yaml
	- kubectl delete -f deploy/service_account.yaml
	- kubectl delete -f deploy/operator.yaml
	- kubectl delete -f deploy/role.yaml
	- kubectl delete -f deploy/role_binding.yaml
	- kubectl delete namespace mobile-security-service-operator

.PHONY: create-service
create-service:
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

.PHONY: delete-app
delete-app:
	@echo Deleting Mobile Security Service and Database:
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_cr.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_cr.yaml

.PHONY: delete-app-only
delete-app-only:
	@echo Deleting Mobile Security Service App only:
	kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_cr.yaml

.PHONY: delete-db-only
delete-db-only:
	@echo Deleting Mobile Security Service Database only:
	kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_cr.yaml

.PHONY: build
build:
	@echo Building operator with the tag $(TAG):
	operator-sdk build $(DOCKER-ORG)/$(DOCKER-REPO):$(TAG)

.PHONY: publish
publish:
	@echo Publishing operator in $(DOCKER-ORG)/$(DOCKER-REPO) with the tag $(TAG):
	docker push $(DOCKER-ORG)/$(DOCKER-REPO):$(TAG)


.PHONY: build-dev
build-dev:
	@echo Buinding operator with the tag $(TAG)-$(DEV):
	operator-sdk build $(DOCKER-ORG)/$(DOCKER-REPO):$(TAG)-$(DEV)

.PHONY: publish-dev
publish-dev:
	@echo Publishing operator in $(DOCKER-ORG)/$(DOCKER-REPO) with the tag $(TAG)-$(DEV):
	docker push $(DOCKER-ORG)/$(DOCKER-REPO):$(TAG)-$(DEV)

.PHONY: debug-setup
debug-setup:
	@echo Exporting WATCH_NAMESPACE=default:
	- export WATCH_NAMESPACE=default
	@echo Create Namespace:
	- kubectl create namespace mobile-security-service-operator
	@echo Installing the CRD:
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_crd.yaml
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceapp_crd.yaml
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityserviceunbind_crd.yaml
	- kubectl create -f deploy/crds/examples/mobile-security-service_v1alpha1_mobilesecurityserviceapp_crd.yaml

.PHONY: vet
vet:
	@echo go vet
	go vet $$(go list ./... | grep -v /vendor/)

.PHONY: fmt
fmt:
	@echo go fmt
	go fmt $$(go list ./... | grep -v /vendor/)
