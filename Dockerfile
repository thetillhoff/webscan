# Dependency stage (cache-friendly)
FROM golang:1.26-alpine AS deps
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
RUN go build -o /artifacts/webscan ./cmd/webscan/

# Build stage for webscan-web
FROM src AS builder-web
WORKDIR /workspace
RUN go build -o /artifacts/webscan-web ./cmd/webscan-web/

# Final stage for CLI version
FROM alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11
ARG TARGETPLATFORM

# Copy the pre-built binary directly from artifacts by name
COPY --from=builder-cli --chmod=755 /artifacts/webscan /usr/local/bin/webscan

WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/webscan"]

# Final stage for web server (API + worker)
FROM alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11 AS web
COPY --from=builder-web --chmod=755 /artifacts/webscan-web /usr/local/bin/webscan-web

WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/webscan-web"]
