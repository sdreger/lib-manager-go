SHELL := /bin/bash
IMAGE_REGISTRY := gitea.dreger.lan/sdreger/lib-manager-go
BUILD_REF = $(shell git rev-parse --short HEAD)
KIND_CLUSTER_NAME := lib-manager-cluster

# lint: run lint checks (https://golangci-lint.run/welcome/quick-start/)
.PHONY: lint
lint:
	go mod tidy -diff
	go mod verify
	test -z "$(gofmt -l .)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@2025.1 -checks=all,-ST1000,-U1000 ./...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6 run -E testifylint -v ./...

# vulncheck: reports known vulnerabilities that affect Go code (https://go.googlesource.com/vuln)
.PHONY: vulncheck
vulncheck:
	go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 ./...

# test: run all application tests with data race detector
.PHONY: test
test:
	go test -v -race -shuffle=on -buildvcs ./...

# cover: run all application tests and generate test coverage report
.PHONY: cover
cover:
	go test -v -race -shuffle=on -buildvcs -coverprofile=/tmp/cover.out.tmp ./...
	grep -v "_mock.go" /tmp/cover.out.tmp > /tmp/cover.out
	go tool cover -html=/tmp/cover.out

# audit: perform full audit check
.PHONY: audit
audit: lint test vulncheck

# docker/build: build the Docker image
.PHONY: docker/build
docker/build:
	docker build -t ${IMAGE_REGISTRY}:${BUILD_REF} -t ${IMAGE_REGISTRY}:latest -f deploy/docker/Dockerfile .

# docker/run: start Docker development environment
.PHONY: docker/run
docker/run:
	docker compose --env-file ./deploy/docker/.env.dev -f ./deploy/docker/compose.yaml up --build

.PHONY: kind/create-cluster
kind/create-cluster:
	kind create cluster --name ${KIND_CLUSTER_NAME} --config deploy/kind/kind-config.yaml
	docker exec -t ${KIND_CLUSTER_NAME}-control-plane update-ca-certificates

.PHONY: kind/delete-cluster
kind/delete-cluster:
	kind delete cluster --name ${KIND_CLUSTER_NAME}

.PHONY: kind/load-image
kind/load-image:
	kind load docker-image --name ${KIND_CLUSTER_NAME} ${IMAGE_REGISTRY}:latest

.PHONY: kustomize/manifests/dev
kustomize/manifests/dev:
	kubectl kustomize deploy/kustomize/overlays/dev

.PHONY: kustomize/manifests/prod
kustomize/manifests/prod:
	kubectl kustomize deploy/kustomize/overlays/prod

.PHONY: kustomize/lint/dev
kustomize/lint/dev:
	kubectl kustomize deploy/kustomize/overlays/dev | kube-linter lint -

.PHONY: kustomize/lint/prod
kustomize/lint/prod:
	kubectl kustomize deploy/kustomize/overlays/prod | kube-linter lint -

.PHONY: kustomize/apply/dev
kustomize/apply/dev:
	kubectl apply --kustomize deploy/kustomize/overlays/dev

.PHONY: kustomize/delete/dev
kustomize/delete/dev:
	kubectl delete --kustomize deploy/kustomize/overlays/dev

.PHONY: kustomize/bootstrap/dev
kustomize/bootstrap/dev: docker/build kind/load-image kustomize/apply/dev

.PHONY: sops/encrypt/secret
sops/encrypt/secret:
	sops --encrypt --input-type=dotenv \
	deploy/kustomize/overlays/prod/.secret.env > deploy/kustomize/overlays/prod/.secret.encrypted.env

.PHONY: sops/edit/secret
sops/edit/secret:
	SOPS_AGE_KEY_FILE=age.agekey \
	sops edit --input-type dotenv --output-type dotenv deploy/kustomize/overlays/prod/.secret.encrypted.env
