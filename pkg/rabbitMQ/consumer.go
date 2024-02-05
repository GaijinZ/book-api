package main

import (
	"library/pkg/rabbitMQ/rabbitMQ"
	"log"
)

func main() {
	rmq, err := rabbitMQ.NewConn()
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ instance: %v", err)
	}
	defer rmq.Close()

	ch := rmq.GetChannel()
	conn := rmq.GetConnection()

	rabbitMR := &rabbitMQ.RabbitMQ{
		Conn: conn,
		Ch:   ch,
	}

	rabbitMR.Consumer()
}
