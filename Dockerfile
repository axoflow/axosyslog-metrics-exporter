FROM golang:1.26.5-bookworm@sha256:1ecb7edf62a0408027bd5729dfd6b1b8766e578e8df93995b225dfd0944eb651 AS builder

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

# Use distroless as minimal base image to package the Go binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/base-debian12@sha256:348dac1808083ccc3366399d6db835875b4eaf7c9b694783f5a3f353c4b58a28 AS prod

WORKDIR /
COPY --from=builder /LICENSE /
COPY --from=builder /NOTICE /
COPY --from=builder /app/dist/axosyslog-metrics-exporter .
ENTRYPOINT ["/axosyslog-metrics-exporter"]
