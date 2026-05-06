#!/usr/bin/env bash

set -euo pipefail

REDIS_CONTAINER_NAME="webscan-test-redis"
REDIS_PORT="63790"

cleanup() {
  docker rm -f "${REDIS_CONTAINER_NAME}" >/dev/null 2>&1 || true
}

trap cleanup EXIT INT TERM

docker rm -f "${REDIS_CONTAINER_NAME}" >/dev/null 2>&1 || true
docker run --name "${REDIS_CONTAINER_NAME}" -p "${REDIS_PORT}:6379" -d redis:7-alpine >/dev/null

# Let redis initialize before webscan starts.
sleep 1

go run ../cmd/webscan-web --port 4173 --redis-addr "127.0.0.1:${REDIS_PORT}"
