CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
PACKAGE=${CURDIR}/cmd/app

all: build test

build: bindir
	go build -o ${BINDIR}/app ${PACKAGE}

test:
	go test ./...

run-package:
	go run ${PACKAGE}
run: build
	docker compose up --force-recreate --build --remove-orphans

bindir:
	mkdir -p ${BINDIR}

LOCAL_BIN:=$(CURDIR)/bin
bin:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

GOOSE = PATH="$$PATH:$(LOCAL_BIN)" goose

db-status: bin
	cd migrations/postgres && $(GOOSE) postgres "postgresql://user:password@localhost:5432/example?sslmode=disable" status
	cd migrations/clickhouse && $(GOOSE) clickhouse "tcp://clickuser:password1@localhost:9000/clickdb" status

db-create: bin
	cd migrations/postgres && $(GOOSE) postgres "postgresql://user:password@localhost:5432/example?sslmode=disable" up
	cd migrations/clickhouse && $(GOOSE) clickhouse "tcp://clickuser:password1@localhost:9000/clickdb" up

db-down: bin
	cd migrations/postgres && $(GOOSE) postgres "postgresql://user:password@localhost:5432/example?sslmode=disable" down
	cd migrations/clickhouse && $(GOOSE) clickhouse "tcp://clickuser:password1@localhost:9000/clickdb" down
