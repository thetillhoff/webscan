FROM alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11
ARG TARGETOS
ARG TARGETARCH

# Copy the pre-built binary directly from artifacts by name
COPY --chmod=755 artifacts/webscan_${TARGETOS}_${TARGETARCH} /usr/local/bin/webscan

WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/webscan"]
