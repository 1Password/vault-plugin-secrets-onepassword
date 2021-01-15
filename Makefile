GOARCH = amd64

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

.DEFAULT_GOAL := help

GIT_BRANCH := $(shell git symbolic-ref --short HEAD)
WORKTREE_CLEAN := $(shell git status --porcelain 1>/dev/null 2>&1; echo $$?)
SCRIPTS_DIR := $(CURDIR)/scripts

PLUGIN_NAME = op-connect
versionFile = $(CURDIR)/.VERSION
curVersion := $(shell cat $(versionFile) | sed 's/^v//')

all: fmt build start

test:	## Run test suite
	go test ./...

test/coverage:	## Run test suite with coverage report
	go test -v ./... -cover


build:	## Build 1Password secrets plugin for Vault
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o vault/plugins/$(PLUGIN_NAME) .

start:	## Start Vault in dev mode and load local plugins
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins

enable:	## Toggle 1Password secrets backend
	vault secrets enable $(PLUGIN_NAME)

clean:	## Remove compiled binary
	rm -f ./vault/plugins/$(PLUGIN_NAME)

fmt:	## Run gofmt on entire project
	go fmt $$(go list ./...)

help:	## Prints this help message
	@grep -E '^[\/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


## Release functions =====================

release/prepare: .check_git_clean	## Bumps version and creates release branch (call with 'release/prepare version=<new_version_number>')

	@test $(version) || (echo "[ERROR] version argument not set."; exit 1)
	@git fetch --quiet origin $(MAIN_BRANCH)

	@echo $(version) | tr -d '\n' | tee $(versionFile) &>/dev/null

	@NEW_VERSION=$(version) $(SCRIPTS_DIR)/prepare-release.sh

release/tag: .check_git_clean	## Creates git tag using version from package.json
	@git pull --ff-only
	@echo "Applying tag 'v$(curVersion)' to HEAD..."
	@git tag --sign "v$(curVersion)" -m "Release v$(curVersion)"
	@echo "[OK] Success!"
	@echo "Remember to call 'git push --tags' to persist the tag."

## Helper functions =====================

.check_git_clean:
ifneq ($(GIT_BRANCH), $(MAIN_BRANCH))
	@echo "[ERROR] Please checkout default branch '$(MAIN_BRANCH)' and re-run this command."; exit 1;
endif
ifneq ($(WORKTREE_CLEAN), 0)
	@echo "[ERROR] Uncommitted changes found in worktree. Address them and try again."; exit 1;
endif


.PHONY: test test/coverage build clean fmt start help enable release/prepare release/tag
