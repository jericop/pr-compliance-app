BUILD_ID := $(shell git rev-parse --short HEAD 2>/dev/null || echo no-commit-id)
WORKSPACE := $(shell pwd)
PKG := $(shell go list ./... | grep -v e2e | grep -v static | grep -v mocks | grep -v testing)
PKG_COMMAS := $(shell go list ./... | grep -v e2e | grep -v static | grep -v mocks | grep -v testing | tr '\n' ',')
IMAGE_NAME := runatlantis/atlantis

.DEFAULT_GOAL := help

.PHONY: help
help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z\/_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: id
id: ## Output BUILD_ID being used
	@echo $(BUILD_ID)

.PHONY: debug
debug: ## Output internal make variables
	@echo BUILD_ID = $(BUILD_ID)
	@echo IMAGE_NAME = $(IMAGE_NAME)
	@echo WORKSPACE = $(WORKSPACE)
	@echo PKG = $(PKG)

.PHONY: build-service
build-service: ## Build the main Go service
	CGO_ENABLED=0 GOOS=linux go build -o build/pr-compliance-app .

.PHONY: build
build: build-service ## Runs make build-service

.PHONY: all
all: build-service ## Runs make build-service

.PHONY: clean
clean: ## Cleans compiled binary
	@rm -f build

.PHONY: go-generate
go-generate: ## Run go generate in all packages
	./scripts/go-generate.sh

.PHONY: regen-mocks
regen-mocks: ## Delete and regenerate all mocks
	sqlc-generate
	find . -type f | grep mocks | xargs rm
	@# not using $(PKG) here because that includes directories that have now
	@# been made empty, causing go generate to fail.
	./scripts/go-generate.sh

.PHONY: sqlc-generate
sqlc-generate:
	sqlc generate

.PHONY: test
test: ## Run tests
	@go test -short $(PKG)

.PHONY: test-all
test-all: ## Run tests including integration
	@go test  $(PKG)

.PHONY: test-coverage
test-coverage: ## Show test coverage
	@mkdir -p .cover
	@go test -covermode atomic -coverprofile .cover/cover.out $(PKG)

.PHONY: test-coverage-html
test-coverage-html: ## Show test coverage and output html
	@mkdir -p .cover
	@go test -covermode atomic -coverpkg $(PKG_COMMAS) -coverprofile .cover/cover.out $(PKG)
	go tool cover -html .cover/cover.out

.PHONY: lint
lint: ## Run linter locally
	golangci-lint run

.PHONY: check-lint
check-lint: ## Run linter in CI/CD. If running locally use 'lint'
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.49.0
	./bin/golangci-lint run -j 4 --timeout 5m

