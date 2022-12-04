BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

docker-up:
	docker-compose -f ./deployments/docker-compose.yaml up -d
	while ! docker-compose -f ./deployments/docker-compose.yaml exec --user postgres db psql -c "select 'db ready!'"; do sleep 1; done;

docker-stop:
	docker-compose -f ./deployments/docker-compose.yaml down

build:
	docker-compose -f ./deployments/docker-compose.yaml build
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" .

run: build docker-up migrate
	$(BIN) serve-http --config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

install-migrator:
	(which goose > /dev/null) || go install github.com/pressly/goose/v3/cmd/goose@latest

migrate: install-migrator
	goose --dir ./migrations postgres "postgres://postgres:password@localhost:15432/postgres?sslmode=disable" up

.PHONY: build run build-img run-img version test lint