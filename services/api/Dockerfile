# syntax=docker/dockerfile:1

# Build from the repository root:
#   docker build -f services/api/Dockerfile -t asl-core-api .

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY pkg ./pkg
COPY services/api ./services/api

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux GOARCH=amd64 go build -o /out/api ./services/api/cmd

FROM alpine:3.19

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /out/api /usr/local/bin/api

USER app

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/api"]
