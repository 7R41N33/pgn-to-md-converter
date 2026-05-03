FROM golang:1.25-alpine AS base

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

FROM base AS dev

COPY . .

EXPOSE 8080

CMD ["/bin/sh"]

FROM base AS builder

COPY . .

RUN go build -o /app/bin/app ./cmd/app

FROM alpine:3.20 AS production

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/app .

RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -D -s /bin/sh appuser && \
    chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["/app/app"]
