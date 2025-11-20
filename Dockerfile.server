FROM golang:1.25.4 AS builder

WORKDIR /app
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o server .

FROM ubuntu:24.04

EXPOSE 6733

RUN apt-get update && apt-get install -y --no-install-recommends \
        bash \
        curl \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server
COPY deviceid.csv deviceid.csv

ENTRYPOINT ["/app/server"]