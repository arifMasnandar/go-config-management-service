    # Stage 1: Builder
    FROM golang:1.24-alpine AS builder

    WORKDIR /app

    COPY go.mod go.sum ./
    RUN go mod download

    COPY . .
    RUN go build -o /app/gocms ./cmd/http/main.go

    # Stage 2: Runner
    FROM alpine:latest

    WORKDIR /app

    COPY --from=builder /app/gocms .
    COPY .env .env

    EXPOSE 8080

    CMD ["./gocms"]