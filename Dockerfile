FROM golang:1.23-bullseye@sha256:462521f1b7cbf410002a8cc4d91bc897f35cd430854d7240596282f9441fe4a7 as builder

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
FROM gcr.io/distroless/base-debian11:latest as prod

WORKDIR /
COPY --from=builder /LICENSE /
COPY --from=builder /NOTICE /
COPY --from=builder /app/dist/axosyslog-metrics-exporter .
ENTRYPOINT ["/axosyslog-metrics-exporter"]
