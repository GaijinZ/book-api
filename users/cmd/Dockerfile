# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder

WORKDIR /userapi

COPY go.mod go.sum ./
RUN go mod download

COPY users ./users
COPY books/models ./books/models
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -o /userapi-build ./users/cmd

FROM alpine:3.18

WORKDIR /

COPY --from=builder /userapi-build ./

EXPOSE 5000

CMD ["./userapi-build"]
