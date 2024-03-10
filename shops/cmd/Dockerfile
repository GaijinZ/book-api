# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder

WORKDIR /shops

COPY go.mod go.sum ./
RUN go mod download

COPY shops ./shops
COPY pkg/config ./pkg/config
COPY pkg/logger ./pkg/logger
COPY pkg/utils ./pkg/utils
COPY pkg/postgres ./pkg/postgres
COPY pkg/redis ./pkg/redis

RUN CGO_ENABLED=0 GOOS=linux go build -o /shops-build ./shops/cmd

FROM alpine:3.18

WORKDIR /

COPY --from=builder /shops-build ./

EXPOSE 5002

CMD ["./shops-build"]
