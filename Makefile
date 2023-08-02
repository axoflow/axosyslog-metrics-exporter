GOARCH    := $(shell go env GOARCH)
GOOS      := $(shell go env GOOS)
GOVERSION := $(shell go env GOVERSION)

DIST ?= ${PWD}/dist
GIT_TAG_COMMAND = git describe --always --dirty --tags 2>/dev/null
GIT_TAG ?= $(shell ${GIT_TAG_COMMAND} || cat .version-tag 2>/dev/null || echo dev)
GO_LDFLAGS := -X 'main.Version=${GIT_TAG}'

# build tools
BIN ?= ${PWD}/bin/${GOOS}/${GOARCH}

LICENSEI := ${BIN}/licensei
LICENSEI_VERSION = v0.8.0
GOLANGCI_LINT := ${BIN}/golangci-lint
GOLANGCI_LINT_VERSION := v1.51.2

.PHONY: fmt
fmt: ## format Go sources
	go fmt ./...

.PHONY: tidy
tidy: ## ensures go.mod dependecies
	find . -iname "go.mod" | sort -r | xargs -L1 sh -c 'set -x; cd $$(dirname $$0); go mod tidy'

.PHONY: build
build: ## build
	go build -o $(DIST)/axosyslog-metrics-exporter -ldflags="$(GO_LDFLAGS)" $(BUILDFLAGS) ./

.PHONY: run
run: ## runs project locally with go run
	go run $(BUILDFLAGS) ./ $(ARGS)

.PHONY: docker-build
docker-build: ## builds docker container locally
	docker build . -t "axoflow.local.dev/axosyslog-metrics-exporter:latest" -f Dockerfile

.PHONY: test
test: ## runs unit tests
	go test ./...

.PHONY: lint
lint: ${GOLANGCI_LINT} ## check coding style
	${GOLANGCI_LINT} run ${LINTER_FLAGS}

## =========================
## ==  Tool dependencies  ==
## =========================

${BIN}:
	mkdir -p $@

${LICENSEI}: ${LICENSEI}_${LICENSEI_VERSION}_${GOVERSION} | ${BIN}

.PHONY: check-license
check-license: ${LICENSEI} .licensei.cache  ## check license + copyright headers
	${LICENSEI} check
	${LICENSEI} header

${LICENSEI}: ${LICENSEI}_${LICENSEI_VERSION}_${GOVERSION}
	ln -sf $(notdir $<) $@

${LICENSEI}_${LICENSEI_VERSION}_${GOVERSION}: IMPORT_PATH := github.com/goph/licensei/cmd/licensei
${LICENSEI}_${LICENSEI_VERSION}_${GOVERSION}: VERSION := ${LICENSEI_VERSION}
${LICENSEI}_${LICENSEI_VERSION}_${GOVERSION}: | ${BIN}
	${go_install_binary}

.licensei.cache: ${LICENSEI}
ifndef GITHUB_TOKEN
	@>&2 echo "WARNING: building licensei cache without Github token, rate limiting might occur."
	@>&2 echo "(Hint: If too many licenses are missing, try specifying a Github token via the environment variable GITHUB_TOKEN.)"
endif
	${LICENSEI} cache

${GOLANGCI_LINT}: ${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}_${GOVERSION} | ${BIN}
	ln -sf $(notdir $<) $@

${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}_${GOVERSION}: IMPORT_PATH := github.com/golangci/golangci-lint/cmd/golangci-lint
${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}_${GOVERSION}: VERSION := ${GOLANGCI_LINT_VERSION}
${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}_${GOVERSION}: | ${BIN}
	${go_install_binary}

define go_install_binary
find ${BIN} -name '$(notdir ${IMPORT_PATH})_*' -exec rm {} +
GOBIN=${BIN} go install ${IMPORT_PATH}@${VERSION}
mv ${BIN}/$(notdir ${IMPORT_PATH}) $@
endef

# Self-documenting Makefile
.DEFAULT_GOAL = help
.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
