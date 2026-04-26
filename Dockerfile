# syntax=docker/dockerfile:1.7

FROM golang:1.25.3-alpine AS builder
WORKDIR /src

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/fluffy-mouton-api ./cmd/api

FROM alpine:3.21
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata wget && \
    addgroup -S app && adduser -S -G app app

COPY --from=builder /out/fluffy-mouton-api /app/fluffy-mouton-api

USER app
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
    CMD wget -qO- "http://127.0.0.1:${PORT:-8080}/api/v1/health/check" || exit 1

CMD ["/app/fluffy-mouton-api"]
