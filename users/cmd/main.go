package main

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"library/pkg/config"
	"library/pkg/logger"
	"library/pkg/utils"
	"library/users/server"
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
	ctx = context.WithValue(ctx, "logger", logger.NewLogger(2))
	log := utils.GetLogger(ctx)

	if err := envconfig.Process("bookapi", &cfg); err != nil {
		log.Fatalf(err.Error())
	}

	go server.Run(&ctx, cfg)

	select {
	case <-interrupt:
		log.Infof("Received a shutdown signal...")
		close(interrupt)
	case <-ctx.Done():
		log.Infof("Context done")
		close(interrupt)
	}
}
