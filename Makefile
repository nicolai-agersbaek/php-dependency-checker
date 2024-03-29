APP_NAME := "php-dependency-checker"
PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

.PHONY: dep build lint test install

## Run the gofmt tool on all package `.go` files
fmt:
	gofmt -w $(GOFMT_FILES)

test: ## Run unit tests
	@go test -short ${PKG_LIST}

lint: ## Perform code linting
	@go get -u golang.org/x/lint/golint
	@golint ${PKG_LIST}

dep: ## Get the dependencies
	@go get -v -d ./...

build: dep ## Build the binary file
	@go build -i -v -o ./bin/${APP_NAME} ./cmd/${APP_NAME}

install: dep ## Install the binary file
	@go install -i -v ./cmd/${APP_NAME}
