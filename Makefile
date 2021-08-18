BIN := "./bin/rotation"
DOCKER_IMG="rotation:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/rotation

run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f Dockerfile .

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

localon:
	docker run -d --name rb -p 15672:15672 -p 5672:5672 rabbitmq:3-management
	pg_ctl start -D '/Users/vladimirastrakhantsev/Library/Application Support/Postgres/var-13'
	psql -p 5432 -U postgres -c "CREATE USER someuser"
	psql -p 5432 -U postgres -c "CREATE DATABASE rotationdb"
	goose -dir migrations postgres "host=localhost user=someuser dbname=rotationdb sslmode=disable" up

localoff:
	goose -dir migrations postgres "user=someuser dbname=rotationdb sslmode=disable" down
	psql -p 5432 -U postgres -c "DROP DATABASE rotationdb"
	psql -p 5432 -U postgres -c "DROP USER someuser"
	pg_ctl stop -D '/Users/vladimirastrakhantsev/Library/Application Support/Postgres/var-13'
	docker stop rb
	docker rm rb

migrate:
	PGPASSWORD=1234 psql -h db -p 5432 -U postgres -c "CREATE USER someuser WITH PASSWORD '1234'"
	PGPASSWORD=1234 psql -h db -p 5432 -U postgres -c "CREATE DATABASE rotationdb"
	goose -dir migrations postgres "host=db user=someuser password=1234 dbname=rotationdb sslmode=disable" up

generate:
	rm -rf internal/pb
	mkdir -p internal/pb
	protoc --proto_path=api/ --go_out=internal/pb --go-grpc_out=internal/pb api/*.proto

evans:
	evans --proto api/rotation.proto repl

up:
	docker compose up -d --build

down:
	docker compose down

.PHONY: build run build-img run-img version test lint migrate generate evans db
