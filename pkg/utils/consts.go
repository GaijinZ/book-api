package utils

import (
	"context"
	"library/pkg/logger"
)

func GetLogger(ctx context.Context) logger.Logger {
	return ctx.Value("logger").(logger.Logger)
}
