FROM golang:1.27 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY internal ./internal
COPY services/catalog ./services/catalog

RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/catalog \
    ./services/catalog/cmd/catalog

FROM alpine:edge

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S appuser \
    && adduser -S -G appuser -H -s /sbin/nologin appuser

COPY --from=builder --chown=appuser:appuser /out/catalog /app/catalog

USER appuser:appuser

ENV GIN_MODE=release

EXPOSE 8081
STOPSIGNAL SIGTERM

ENTRYPOINT ["/app/catalog"]
