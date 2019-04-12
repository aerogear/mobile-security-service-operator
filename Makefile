APP_NAME = mobile-security-service-operator
ORG_NAME = aerogear
PKG = github.com/$(ORG_NAME)/$(APP_NAME)
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

.PHONY: run-local
run-local:
	@echo Running operator locally:
	- export OPERATOR_NAME=mobile-security-service-operator
	- make create-all
	- operator-sdk up local --namespace=mobile-security-service-operator

.PHONY: create-all
create-all:
	@echo Create Mobile Security Service Operator and Service in the namespace "mobile-security-service-operator":
	make create-oper
	make create-app
	make create-bind

.PHONY: delete-all
delete-all:
	@echo Delete Mobile Security Service Operator, Service and namespace "mobile-security-service-operator":
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_cr.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_cr.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicebind_cr.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_crd.yaml
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicebind_crd.yaml
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
	- kubectl create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicebind_crd.yaml
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
	- kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicebind_crd.yaml
	- kubectl delete -f deploy/cluster_role.yaml
	- kubectl delete -f deploy/cluster_role_binding.yaml
	- kubectl delete -f deploy/service_account.yaml
	- kubectl delete -f deploy/operator.yaml
	- kubectl delete -f deploy/role.yaml
	- kubectl delete -f deploy/role_binding.yaml
	- kubectl delete namespace mobile-security-service-operator

.PHONY: create-bind
create-bind:
	@echo Creating Mobile Security Service Bind:
	kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicebind_cr.yaml

.PHONY: delete-bind
delete-bind:
	@echo Deleting Mobile Security Service Bind:
	kubectl delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicebind_cr.yaml

.PHONY: create-app
create-app:
	@echo Creating Mobile Security Service and Database into project:
	kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_cr.yaml
	kubectl apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservicedb_cr.yaml


.PHONY: create-app-only
create-app-only:
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
	@echo Buinding operator with the tag $(TAG):
	operator-sdk build $(DOCKER-ORG)/$(DOCKER-REPO):$(TAG)

.PHONY: publish
publish:
	@echo Publishing operator in $(DOCKER-ORG)/$(DOCKER-REPO) with the tag $(TAG):
	docker push $(DOCKER-ORG)/$(DOCKER-REPO):$(TAG)


.PHONY: build-dev
build-dev:
	@echo Buinding operator with the tag $(TAG)$(DEV):
	operator-sdk build $(DOCKER-ORG)/$(DOCKER-REPO):$(TAG)-$(DEV)


.PHONY: publish-dev
publish-dev:
	@echo Publishing operator in $(DOCKER-ORG)/$(DOCKER-REPO) with the tag $(TAG)-$(DEV):
	docker push $(DOCKER-ORG)/$(DOCKER-REPO):$(TAG)-$(DEV)

.PHONY: vet
vet:
	@echo go vet
	go vet $$(go list ./... | grep -v /vendor/)

.PHONY: fmt
fmt:
	@echo go fmt
	go fmt $$(go list ./... | grep -v /vendor/)