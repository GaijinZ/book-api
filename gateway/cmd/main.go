package main

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"library/gateway/server"
	"library/pkg/config"
	"library/pkg/logger"
	"library/pkg/utils"
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

	go gateway.Run(ctx, cfg, ":"+cfg.GatewayServerPort)

	for {
		select {
		case sig := <-interrupt:
			log.Infof("Received signal: %v\n", sig)
		case <-ctx.Done():
			log.Infof("Context done")
		}
	}
}
