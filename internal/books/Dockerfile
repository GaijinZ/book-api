# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder

WORKDIR /bookapi

COPY go.mod go.sum ./
RUN go mod download

COPY internal/books ./internal/books
COPY cmd/books ./cmd/books
COPY database/postgres ./database/postgres
COPY internal/users/models ./internal/users/models
COPY utils ./utils

RUN CGO_ENABLED=0 GOOS=linux go build -o /bookapi-build ./cmd/books

FROM alpine:3.18

WORKDIR /

COPY --from=builder /bookapi-build ./

EXPOSE 5001

CMD ["./bookapi-build"]