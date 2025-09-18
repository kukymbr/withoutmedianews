GOLANGCI_LINT_VERSION := 2.1.6
GO_TEST_CMD := go test -race
GO_TEST_API_CMD := go test -tags api_test
DB_PORT := 5432
DB_DSN := postgres://postgres:postgres@localhost:$(DB_PORT)/postgres?sslmode=disable

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
	GENNA_DATABASE_DSN=$(DB_DSN) go generate ./...
	go mod tidy

generate_genna:
	go tool genna model -c $(DB_DSN) -o internal/db/models.gen.go -t public.* -f

lint:
	if [ ! -f ./bin/golangci-lint ]; then \
  		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/v${GOLANGCI_LINT_VERSION}/install.sh | sh -s -- -b "./bin" "v${GOLANGCI_LINT_VERSION}"; \
  	fi;
	./bin/golangci-lint run ./...

test:
	$(GO_TEST_CMD) ./...

test_api:
	DB_PORT=$(DB_PORT) $(GO_TEST_API_CMD) ./internal/api/httptests/...

clean:
	go clean

run_dev_deps:
	docker compose -f ./.deps/compose.yml up -d dev_database

mfd-xml:
	go tool mfd-generator xml -c "$(DB_DSN)" -m ./db/model/withoutmedianews.mfd -n withoutmedianews

mfd-model:
	go tool mfd-generator model -m ./db/model/withoutmedianews.mfd -p db -o ./internal/db

mfd-repo:
	go tool mfd-generator repo -m ./db/model/withoutmedianews.mfd -p db -o ./internal/db