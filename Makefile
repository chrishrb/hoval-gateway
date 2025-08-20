GOCMD=go

.PHONY: all
all: generate build format lint gosec integration ## Format, lint, build and test

.PHONY: generate
generate: ## Generate
	${GOCMD} generate ./...
	${GOCMD} mod tidy

.PHONY: test
test: ## Test
	${GOCMD} test ./...

.PHONY: build
build: ## Build
	${GOCMD} build -o bin/ ./...

.PHONY: integration
integration: ## Run unit and integration tests
	${GOCMD} test -tags=integration ./...

.PHONY: format
format: ## Format code
	${GOCMD} fmt ./...

.PHONY: lint
lint: ## Run linter
	golangci-lint run ./...

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
