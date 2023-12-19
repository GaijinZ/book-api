package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"library/pkg/config"
	"library/pkg/postgres"
	"library/pkg/utils"
	"library/shops/handler"
	"library/shops/repository"
	"time"
)

func Run(ctx *context.Context, cfg config.GlobalEnv) {
	log := utils.GetLogger(*ctx)

	configDB := postgres.DBConfig{
		DriverName:      "postgres",
		DataSourceName:  cfg.PostgresBooks,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	}

	db, err := postgres.NewDB(*ctx, configDB)
	if err != nil {
		log.Errorf("Failed to configure db connection: %v", err)
	}
	defer db.Close()

	shopRepository := repository.NewShopRepository(*ctx, *db)
	shopHandler := handler.NewShopHandler(*ctx, shopRepository)

	router := gin.Default()

	v1 := router.Group("/v1/shops")

	v1.POST("/load-books", shopHandler.LoadBooks)

	router.Run(":" + cfg.ShopsServerPort)
}
