package main

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"library/pkg/config"
	"library/pkg/logger"
	"library/pkg/rabbitMQ/rabbitMQ"
	"library/pkg/utils"
)

func main() {
	var cfg config.GlobalEnv
	var ctx context.Context

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = context.WithValue(ctx, "logger", logger.NewLogger(2))
	log := utils.GetLogger(ctx)

	if err := envconfig.Process("bookapi", &cfg); err != nil {
		log.Fatalf(err.Error())
	}

	rmq, err := rabbitMQ.NewConn(cfg)
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
