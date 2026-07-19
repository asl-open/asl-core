# Build from the repository root:
#   docker build -f services/api/Dockerfile -t asl-core-api .

FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY pkg ./pkg
COPY services/api ./services/api

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./services/api/cmd

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/api /api

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/api"]
