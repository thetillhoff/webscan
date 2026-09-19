# Dependency stage (cache-friendly)
FROM golang:1.27-alpine@sha256:4cb7ac979db5fcc41cae44b2227ba5ab8a51e8807f40d9ba4dee20a0ad960b5b AS deps
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
FROM alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6
ARG TARGETPLATFORM

# Copy the pre-built binary directly from artifacts by name
COPY --from=builder-cli --chmod=755 /artifacts/webscan /usr/local/bin/webscan

RUN adduser -D -u 10001 app
USER app
WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/webscan"]

# Final stage for web server (API + worker)
FROM alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6 AS web
COPY --from=builder-web --chmod=755 /artifacts/webscan-web /usr/local/bin/webscan-web

RUN adduser -D -u 10001 app
USER app
WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/webscan-web"]
