GOARCH = amd64

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

.DEFAULT_GOAL := all

all: fmt build start

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o vault/plugins/op-connect .

start:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins

enable:
	vault secrets enable op-connect

clean:
	rm -f ./vault/plugins/op-connect

fmt:
	go fmt $$(go list ./...)

test:	## Run test suite
	go test ./...

.PHONY: build clean fmt start enable