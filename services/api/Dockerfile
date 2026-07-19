# syntax=docker/dockerfile:1

# Build from the repository root:
#   docker build -f services/api/Dockerfile -t asl-core-api .

FROM golang:1.25.5-alpine AS builder

WORKDIR /src

ENV CGO_ENABLED=0 GOTOOLCHAIN=local

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY pkg ./pkg
COPY services/api ./services/api

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/api ./services/api/cmd

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/api /api

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/api"]
