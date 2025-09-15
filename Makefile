GOLANGCI_LINT_VERSION := 2.1.6
GO_TEST_CMD := go test -race

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
	$(GO_TEST_CMD) -coverprofile=coverage_out ./...
	go tool cover -func=coverage_out
	go tool cover -html=coverage_out -o coverage.html
	rm -f coverage_out

clean:
	go clean

run_dev_deps:
	docker compose -f ./.devdocker/compose.yml up -d dev_postgres
