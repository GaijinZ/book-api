package main

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"library/pkg/config"
	"library/pkg/logger"
	"library/shops/server"
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

	go server.Run(&ctx, ":"+cfg.ShopsServerPort)

	select {
	case sig := <-interrupt:
		fmt.Printf("Received signal: %v\n", sig)
	case <-ctx.Done():
		fmt.Println("Context done")
	}
}
