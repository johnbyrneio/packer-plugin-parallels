NAME=parallels
BINARY=packer-plugin-${NAME}
OS=$(uname go env GOOS)
PLUGIN_FQN=$(shell grep -E '^module' <go.mod | sed -E 's/module *//')

COUNT?=1
TEST?=$(shell go list ./...)
HASHICORP_PACKER_PLUGIN_SDK_VERSION?=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f2)
PLUGIN_PATH=$(shell echo "${PLUGIN_FQN}" | sed 's/packer-plugin-//')
PACKER_FOLDER=$(shell dirname $(shell which packer))
VERSION=$(shell grep -E '^\tVersion = ' <./version/version.go | sed -E 's/\tVersion = *//')

.PHONY: dev

build:
	@go build -o ${BINARY}

dev:
	@./scripts/build_debug.sh

test:
	@go test -race -count $(COUNT) $(TEST) -timeout=3m

install-packer-sdc: ## Install packer sofware development command
	go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@${HASHICORP_PACKER_PLUGIN_SDK_VERSION}

plugin-check: install-packer-sdc build
	@packer-sdc plugin-check ${BINARY}

testacc: dev
	@PACKER_ACC=1 go test -count $(COUNT) -v $(TEST) -timeout=120m

generate: install-packer-sdc
	@go generate ./...
	@if [ -d ".docs" ]; then rm -r ".docs"; fi
	packer-sdc renderdocs -src "docs" -partials docs-partials/ -dst ".docs/"
	@./.web-docs/scripts/compile-to-webdocs.sh "." ".docs" ".web-docs" "Parallels"
	@rm -r ".docs"
	# checkout the .docs folder for a preview of the docs
