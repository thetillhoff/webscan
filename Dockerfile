FROM alpine:3.23.3@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659
ARG TARGETOS
ARG TARGETARCH

# Copy the pre-built binary directly from artifacts by name
COPY --chmod=755 artifacts/webscan_${TARGETOS}_${TARGETARCH} /usr/local/bin/webscan

WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/webscan"]
