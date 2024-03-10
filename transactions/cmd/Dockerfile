# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder

WORKDIR /transactions

COPY go.mod go.sum ./
RUN go mod download

COPY transactions ./transactions
COPY users/models ./users/models
COPY books/models ./books/models
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -o /transactions-build ./transactions/cmd

FROM alpine:3.18

WORKDIR /

COPY --from=builder /transactions-build ./

EXPOSE 5003

CMD ["./transactions-build"]
