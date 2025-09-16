GOLANGCI_LINT_VERSION := 2.1.6
GO_TEST_CMD := go test -race
GO_TEST_API_CMD := go test -tags api_test

all:
	$(MAKE) clean
	$(MAKE) prepare
	$(MAKE) validate
	$(MAKE) build

prepare:
	go install ./...
	go fmt ./...
	$(MAKE) generate

validate:
	go vet ./...
	$(MAKE) lint
	$(MAKE) test

build:
	go build $(GO_BUILD_ARGS) -o bin/withoutmedianews ./cmd/withoutmedianews

generate:
	go generate ./...
	go mod tidy

lint:
	if [ ! -f ./bin/golangci-lint ]; then \
  		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/v${GOLANGCI_LINT_VERSION}/install.sh | sh -s -- -b "./bin" "v${GOLANGCI_LINT_VERSION}"; \
  	fi;
	./bin/golangci-lint run ./...

test:
	$(GO_TEST_CMD) ./...

test_api:
	$(GO_TEST_API_CMD) ./internal/api/httptests/...

clean:
	go clean

run_dev_deps:
	docker compose -f ./.deps/compose.yml up -d dev_database
