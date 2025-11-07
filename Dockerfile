FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" -o /bin/mpb ./cmd/main.go

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" -o /bin/migrate ./cmd/migrate.go

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" -o /bin/healthcheck ./cmd/healthcheck/main.go

FROM gcr.io/distroless/static-debian11

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /bin/mpb /
COPY --from=builder /bin/migrate /
COPY --from=builder /bin/healthcheck /

COPY --from=builder /app/migrations ./migrations

LABEL maintainer="Maks Kozlov <maks210306@yandex.by>"
LABEL version="1.0.0"

EXPOSE 8000

USER 65532:65532

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/healthcheck"]

CMD ["/mpb"]