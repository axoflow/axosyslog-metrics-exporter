DIST ?= ${PWD}/dist
GIT_TAG_COMMAND = git describe --always --dirty --tags 2>/dev/null
GIT_TAG ?= $(shell ${GIT_TAG_COMMAND} || cat .version-tag 2>/dev/null || echo dev)
GO_LDFLAGS := -X 'main.Version=${GIT_TAG}'

.PHONY: fmt
fmt: ## format Go sources
	go fmt ./...

.PHONY: tidy
tidy: ## ensures go.mod dependecies
	find . -iname "go.mod" | sort -r | xargs -L1 sh -c 'set -x; cd $$(dirname $$0); go mod tidy'

.PHONY: build
build: ## build
	go build -o $(DIST)/metrics-exporter -ldflags="$(GO_LDFLAGS)" $(BUILDFLAGS) ./

.PHONY: run
run: ## runs project locally with go run
	go run $(BUILDFLAGS) ./ $(ARGS)

.PHONY: docker-build
docker-build: ## builds docker container locally
	docker build . -t "axoflow.local.dev/metrics-exporter:latest" -f Dockerfile

# Self-documenting Makefile
.DEFAULT_GOAL = help
.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
