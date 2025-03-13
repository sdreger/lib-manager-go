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
	docker build -t ${IMAGE_REGISTRY}:${BUILD_REF} -f deploy/docker/Dockerfile .

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
	kind load docker-image --name ${KIND_CLUSTER_NAME} ${IMAGE_REGISTRY}:${BUILD_REF}
