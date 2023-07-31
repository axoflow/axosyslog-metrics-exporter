FROM golang:1.20-bullseye as builder

RUN apt-get update \
  && apt-get install --yes entr \
  && rm -rf /var/lib/apt/lists/*

COPY LICENSE /

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

# add shell for debugging purposes
# COPY --from=busybox:1.35.0-uclibc /bin /bin

WORKDIR /
COPY --from=builder /LICENSE /
COPY --from=builder /app/dist/metrics-exporter .
ENTRYPOINT ["/metrics-exporter"]
