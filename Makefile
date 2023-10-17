BUILD_ID := $(shell git rev-parse --short HEAD 2>/dev/null || echo no-commit-id)
WORKSPACE := $(shell pwd)
PKG := $(shell go list ./... | grep -v e2e | grep -v static | grep -v mocks | grep -v testing)
PKG_COMMAS := $(shell go list ./... | grep -v e2e | grep -v static | grep -v mocks | grep -v testing | tr '\n' ',')

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

.PHONY: sqlc-generate
sqlc-generate: ## Generates database code from sqlc.yaml config.
	sqlc generate

.PHONY: go-generate
go-generate: ## Run go generate in all packages
	./scripts/go-generate.sh

.PHONY: regen-db-and-mocks
regen-db-and-mocks: sqlc-generate go-generate	## Delete and regenerate all db and mock code

.PHONY: test
test: ## Run tests
	@go test -v -short $(PKG)

.PHONY: test-all
test-all: ## Run tests including integration
	@go test -v $(PKG)

.PHONY: test-coverage
test-coverage: ## Show test coverage
	@mkdir -p .cover
	@go test -v -covermode atomic -coverprofile .cover/cover.out $(PKG)

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
