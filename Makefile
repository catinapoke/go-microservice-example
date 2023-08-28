CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
PACKAGE=/cmd/app

all: build test

build: bindir
	go build -o ${BINDIR}/app ${PACKAGE}

test:
	go test ./...

run:
	go run ${PACKAGE}

LOCAL_BIN:=$(CURDIR)/bin
bin:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

GOOSE = PATH="$$PATH:$(LOCAL_BIN)" goose

db-status: bin
	cd migrations && $(GOOSE) postgres "postgresql://user:password@localhost:5432/example?sslmode=disable" status

db-create: bin
	cd migrations && $(GOOSE) postgres "postgresql://user:password@localhost:5432/example?sslmode=disable" up

db-down: bin
	cd migrations && $(GOOSE) postgres "postgresql://user:password@localhost:5432/example?sslmode=disable" down 0
