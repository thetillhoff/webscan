# Dependency stage (cache-friendly)
FROM golang:1.27-alpine@sha256:4c9fe60190a2a3350ddc51de80d0224b8a6698d12bdfc999fee45ea9d6c46dbc AS deps
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download

# Shared source stage
FROM deps AS src
WORKDIR /workspace
COPY . .

# Build stage for webscan CLI
FROM src AS builder-cli
WORKDIR /workspace
ARG VERSION=dev
RUN go build -ldflags="-X main.version=${VERSION}" -o /artifacts/webscan ./cmd/webscan/

# Build stage for webscan-web
FROM src AS builder-web
WORKDIR /workspace
ARG VERSION=dev
RUN go build -ldflags="-X main.version=${VERSION}" -o /artifacts/webscan-web ./cmd/webscan-web/

# Final stage for CLI version
FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b
ARG TARGETPLATFORM

# Copy the pre-built binary directly from artifacts by name
COPY --from=builder-cli --chmod=755 /artifacts/webscan /usr/local/bin/webscan

RUN adduser -D -u 10001 app
USER app
WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/webscan"]

# Final stage for web server (API + worker)
FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b AS web
COPY --from=builder-web --chmod=755 /artifacts/webscan-web /usr/local/bin/webscan-web

RUN adduser -D -u 10001 app
USER app
WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/webscan-web"]
