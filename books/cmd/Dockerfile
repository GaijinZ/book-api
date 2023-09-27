# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder

WORKDIR /bookapi

COPY go.mod go.sum ./
RUN go mod download

COPY books ./books
COPY users/models ./users/models
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -o /bookapi-build ./books/cmd

FROM alpine:3.18

WORKDIR /

COPY --from=builder /bookapi-build ./

EXPOSE 5001

CMD ["./bookapi-build"]
