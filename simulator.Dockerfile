FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
      bash curl \
      ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ENV S3_URL="https://sy-fleet-interview-assets.s3.us-east-2.amazonaws.com/device-simulator-linux-arm64"
ENV BIN_PATH=/usr/local/bin/app

RUN curl -fSL "$S3_URL" -o "$BIN_PATH" && chmod +x "$BIN_PATH"

ENTRYPOINT ["/usr/local/bin/app"]