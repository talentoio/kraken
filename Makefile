TOOLS_BIN = $(shell pwd)/tools/bin
GOLANGCI_LINT_VERSION = 1.56.2
CUR_GOLANGCI_LINT_VERSION = $(shell $(TOOLS_BIN)/golangci-lint --version | sed 's/^.*version //g'| sed 's/ .*//g')

## format-code      : go fmt, go mod tidy, download and run goimports
.PHONY: format-code
format-code:
	go fmt ./...
	go mod tidy
	@echo "goimports -w ."
	@if [ -f $(TOOLS_BIN)/goimports ]; then \
      	$(TOOLS_BIN)/goimports -w . ;\
      else \
      	GOBIN=$(TOOLS_BIN) go install golang.org/x/tools/cmd/goimports@v0.10.0; \
      	$(TOOLS_BIN)/goimports -w . ;\
    fi


## mocks         	  : generate mocks
.PHONY: mocks
mocks:
	go install github.com/golang/mock/mockgen@v1.6.0
	go generate ./...
	make format-code


## swagger          : generate swagger documentation
.PHONY: swagger
swagger:
	go install github.com/swaggo/swag/cmd/swag@v1.8.12
	swag fmt
	swag init --dir ./internal/resthttp/,./internal/services/domain --generalInfo /routes.go --outputTypes yaml


## test             : run unit tests
.PHONY: test
test:
	go test -v -race -cover -timeout 15s ./...


## linter           : download and run golangci-lint
.PHONY: linter
linter:
	@echo 'current linter version: ' $(CUR_GOLANGCI_LINT_VERSION)
	@if [ '$(CUR_GOLANGCI_LINT_VERSION)' != '$(GOLANGCI_LINT_VERSION)' ]; then\
	  echo 'install linter' $(GOLANGCI_LINT_VERSION);\
	  rm $(TOOLS_BIN)/golangci-lint;\
	  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v$(GOLANGCI_LINT_VERSION);\
	  mkdir -p $(TOOLS_BIN);\
	  mv ./bin/golangci-lint $(TOOLS_BIN)/;\
	  rmdir ./bin || echo;\
	 fi
	 @echo 'running linter'
	 $(TOOLS_BIN)/golangci-lint run

## run           : build and run app using docker-compose
.PHONY: run
run:
	docker-compose  up -d