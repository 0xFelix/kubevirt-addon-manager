# Image URL to use in building/pushing/deploying targets
IMG_REPOSITORY ?= quay.io/fmatouschek/kubevirt-addon-manager
IMG_TAG ?= latest
IMG ?= ${IMG_REPOSITORY}:${IMG_TAG}

PROJECT_DIR := ${shell dirname ${abspath ${lastword ${MAKEFILE_LIST}}}}
LOCALBIN ?= ${PROJECT_DIR}/bin
${LOCALBIN}:
	mkdir -p ${LOCALBIN}

# Use oc if available
ifeq (,${shell which oc})
KUBECTL = kubectl
else
KUBECTL = oc
endif

KUBECONFIG ?= ~/.kube/config

.PHONY: vendor
vendor: ## Vendor dependencies
	go mod vendor
	go mod tidy

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: lint
lint: vet golangci-lint ## Lint source code.
	${GOLANGCILINT} run --timeout 4m0s ./...

.PHONY: build
build: ## Build manager binary
	go build -mod vendor -o bin/manager main.go

.PHONY: run
run: manager ## Run the manager controller outside the cluster
	${LOCALBIN}/manager controller --kubeconfig ${KUBECONFIG}

.PHONY: container-build
container-build: ## Build the container image
	podman build -t ${IMG} .

.PHONY: container-push
container-push: ## Push the container image
	podman push ${IMG}

.PHONY: deploy
deploy: kustomize ## Deploy controller to the K8s cluster specified in KUBECONFIG.
	cp config/manager/deployment.yaml config/manager/deployment.yaml.tmp
	cd config/manager && ${KUSTOMIZE} edit set image manager=${IMG}
	${KUSTOMIZE} build config/default | ${KUBECTL} apply -f -
	mv config/manager/deployment.yaml.tmp config/manager/deployment.yaml

.PHONY: undeploy
undeploy: kustomize ## Undeploy controller from the K8s cluster specified in KUBECONFIG.
	${KUSTOMIZE} build config/default | ${KUBECTL} delete -f -

.PHONY: golangci-lint
GOLANGCILINT := ${LOCALBIN}/golangci-lint
golangci-lint: ${GOLANGCILINT} ## Download golangci-lint
${GOLANGCILINT}: ${LOCALBIN}
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${LOCALBIN}

.PHONY: kustomize
KUSTOMIZE := ${LOCALBIN}/kustomize
kustomize: ${KUSTOMIZE} ## Download kustomize
${KUSTOMIZE}: ${LOCALBIN}
	GOBIN=${LOCALBIN} go install sigs.k8s.io/kustomize/kustomize/v4@latest
