BIN := "./bin/rotation"
DOCKER_IMG="rotation:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/rotation

run: build
	$(BIN) -config ./configs/config.toml

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
	# go test -race ./internal/... ./pkg/...
	go test -race -count 100 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	psql -p 5432 -h localhost -U postgres -c "CREATE USER someuser"
	psql -p 5432 -h localhost -U postgres -c "CREATE DATABASE rotationdb"
	goose -dir migrations postgres "user=someuser dbname=rotationdb sslmode=disable" up
unmigrate: 
	goose -dir migrations postgres "user=someuser dbname=rotationdb sslmode=disable" down
	psql -p 5432 -h localhost -U postgres -c "DROP DATABASE rotationdb"
	psql -p 5432 -h localhost -U postgres -c "DROP USER someuser"

.PHONY: build run build-img run-img version test lint migrate unmigrate
