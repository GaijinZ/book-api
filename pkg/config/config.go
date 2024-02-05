package config

type GlobalEnv struct {
	GatewayServerPort      string `envconfig:"gateway_server_port"`
	UsersServerPort        string `envconfig:"users_server_port"`
	BooksServerPort        string `envconfig:"books_server_port"`
	ShopsServerPort        string `envconfig:"shops_server_port"`
	TransactionsServerPort string `envconfig:"transactions_server_port"`
	RedisPort              string `envconfig:"redis_port"`
	PostgresBooks          string `envconfig:"postgres_books"`
	RabbitMQ               string `envconfig:"rabbitmq"`
	SecretKey              string `envconfig:"secret_key"`
	Network                string `envconfig:"network"`
	AutoSplitVar           string `split_words:"true"`
}
