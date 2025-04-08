.PHONY: count-go-lines test build help check-build check-vet check-all install

VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
GIT_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GO_VERSION := $(shell go version | awk '{print $$3}')
GO_BIN := $(shell go env GOPATH)/bin
SHELL_NAME := $(shell basename $$SHELL)

help: ## Display this help message
	@echo "Makefile for SOA (State of the Art), the targets are listed below:\\n"
	@awk 'BEGIN {FS = ":.*?## "; OFS = " : ";} \
		/^[a-zA-Z_-]+:.*?##/ { \
			printf "\033[36m%-22s\033[0m%s\n", $$1, $$2; \
		}' $(MAKEFILE_LIST)
	@echo "\\n"

build: ## Compile the executable under build/mpat
	go build -ldflags "-X main.Version=$(VERSION) \
		-X main.GitCommit=$(GIT_COMMIT) \
		-X main.BuildDate=$(BUILD_DATE) \
		-X main.GoVersion=$(GO_VERSION)" \
		-o build/soa cmd/soa/main.go

install: build ## Copy the binary into $GOPATH/bin foler.
	install ./build/soa "$(GO_BIN)/soa"

	@echo "\nTo enable shell completions:"
	@printf "\033[36msource <(soa completion $(SHELL_NAME))\033[0m\n"
	@echo ""


test: ## Run go test for all of the packages
	go test ./...

check-build: ## Check for compilation issues
	go build ./...

vet: ## Check for linting issues
	go vet ./...

check-all: check-vet check-build ## Check for building and linting issues

clean: ## Run go clean command, also removes the mod cache
	go clean --modcache

count-go-lines: ## Return the number of lines for each go file
	@find . -name "*.go" -exec wc -l {} +

version: ## Get the version the binary will be compiled to
	@echo $(VERSION)
