FROM golang:1.27 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY internal ./internal
COPY services/gateway ./services/gateway

RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/gateway \
    ./services/gateway/cmd/gateway

FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S appuser \
    && adduser -S -G appuser -H -s /sbin/nologin appuser

COPY --from=builder --chown=appuser:appuser /out/gateway /app/gateway

USER appuser:appuser

ENV GIN_MODE=release

EXPOSE 8080
STOPSIGNAL SIGTERM

ENTRYPOINT ["/app/gateway"]
