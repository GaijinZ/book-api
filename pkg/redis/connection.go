package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"library/pkg/config"
)

type Client struct {
	Client *redis.Client
}

func NewRedis(cfg config.GlobalEnv) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + cfg.RedisPort,
		Password: "",
		DB:       0,
	})

	return &Client{
		Client: client,
	}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.Client.Ping(ctx).Result()

	return err
}
