# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder

WORKDIR /gateway

COPY go.mod go.sum ./
RUN go mod download

COPY gateway ./gateway

RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway-build ./gateway/cmd

FROM alpine:3.18

WORKDIR /

COPY --from=builder /gateway-build ./

EXPOSE 8080

CMD ["./gateway-build"]