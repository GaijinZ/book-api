# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder

WORKDIR /userapi

COPY go.mod go.sum ./
RUN go mod download

COPY internal/users ./internal/users
COPY cmd/users ./cmd/users
COPY database/postgres ./database/postgres
COPY utils ./utils

RUN CGO_ENABLED=0 GOOS=linux go build -o /userapi-build ./cmd/users

FROM alpine:3.18

WORKDIR /

COPY --from=builder /userapi-build ./

EXPOSE 5000

CMD ["./userapi-build"]