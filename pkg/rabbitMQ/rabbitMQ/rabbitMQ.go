package rabbitMQ

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"library/pkg/config"
	"library/pkg/utils"
	"log"
)

type RabbitMQ struct {
	Conn *amqp091.Connection
	Ch   *amqp091.Channel
}

func NewConn(cfg config.GlobalEnv) (*RabbitMQ, error) {
	connection, err := amqp091.Dial(cfg.RabbitMQ)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		panic(err)
	}

	ch, err := connection.Channel()
	if err != nil {
		connection.Close()
		return nil, err
	}

	return &RabbitMQ{
		Conn: connection,
		Ch:   ch,
	}, nil
}

func (r *RabbitMQ) Producer(ctx context.Context, body string) {
	log := utils.GetLogger(ctx)

	queueName := "activation_queue"
	q, err := r.Ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	err = r.Ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Errorf("Failed to publish a message: %v", err)
	}

	log.Infof("Sent activation email for user")
}

func (r *RabbitMQ) Consumer() {
	queueName := "activation_queue"
	q, err := r.Ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := r.Ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
	}

	var forever chan struct{}

	go func() {
		for msg := range msgs {
			sendActivationEmail(msg.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (r *RabbitMQ) Close() {
	if r.Ch != nil {
		r.Ch.Close()
	}

	if r.Conn != nil {
		r.Conn.Close()
	}
}

func (r *RabbitMQ) GetConnection() *amqp091.Connection {
	return r.Conn
}

func (r *RabbitMQ) GetChannel() *amqp091.Channel {
	return r.Ch
}

func sendActivationEmail(body []byte) {
	fmt.Println(string(body))
	return
}
