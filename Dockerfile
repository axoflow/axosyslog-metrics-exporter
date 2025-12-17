FROM golang:1.24.11-bookworm@sha256:fc58bb98c4b7ebc8211c94df9dee40489e48363c69071bceca91aa59023b0dee AS builder

RUN apt-get update \
  && apt-get install --yes entr \
  && rm -rf /var/lib/apt/lists/*

COPY LICENSE /
COPY NOTICE /

WORKDIR /app

# Copy the Go Modules manifests
COPY go.* ./
COPY pkg/ ./pkg/

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . .

# Cache intermediate build artifacts
# RUN go build ./...
RUN make build

# Use distroless as minimal base image to package the axo-controller binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/base-debian12:latest AS prod

WORKDIR /
COPY --from=builder /LICENSE /
COPY --from=builder /NOTICE /
COPY --from=builder /app/dist/axosyslog-metrics-exporter .
ENTRYPOINT ["/axosyslog-metrics-exporter"]
