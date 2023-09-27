# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder

WORKDIR /userapi

COPY go.mod go.sum ./
RUN go mod download

COPY users ./users
COPY pkg/middleware ./pkg/middleware
COPY pkg/password-hash.go ./pkg/password-hash.go
COPY pkg/sql.go ./pkg/sqp.go

RUN CGO_ENABLED=0 GOOS=linux go build -o /userapi-build ./users/cmd

FROM alpine:3.18

WORKDIR /

COPY --from=builder /userapi-build ./

EXPOSE 5000

CMD ["./userapi-build"]
