help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run unit tests
	go test -gcflags=-l -v -covermode=count ./...

.PHONY: cover
cover: ## Run unit tests with coverage report
	go test -gcflags=-l -coverprofile=coverage.txt ./...
	go tool cover -func=coverage.txt

.PHONY: all-tests
all-tests: test cover ## run all tests

.PHONY: lint
lint: ## Run linting
	@echo be sure golangci-lint is installed: https://golangci-lint.run/usage/install/
	golangci-lint run --config .golangci.yml

.PHONY: goimports
goimports: ## run goimports updating files in place
	@echo be sure goimports is installed
	goimports -w .

.PHONY: goimports-check
goimports-check: ## run goimports displaying diffs
	@echo be sure goimports is installed
	goimports -d . | (! grep .)

.PHONY: all-checks
all-checks: lint goimports ## run all formatters
	go mod tidy
	go vet ./...
