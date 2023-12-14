package main

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"library/pkg/config"
	"library/pkg/logger"
	"library/users/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var cfg config.GlobalEnv
	var ctx context.Context

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = context.WithValue(ctx, "logger", logger.NewLogger())

	if err := envconfig.Process("bookapi", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	go server.Run(ctx, ":"+cfg.UsersServerPort)

	select {
	case <-interrupt:
		logrus.Infof("Received a shutdown signal...")
		close(interrupt)
	case <-ctx.Done():
		logrus.Infof("Context done")
		close(interrupt)
	}
}
