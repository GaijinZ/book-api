package config

type GlobalEnv struct { // cleanup unused variables
	GatewayServerPort        string `envconfig:"gateway_server_port"`
	UsersServerPort          string `envconfig:"users_server_port"`
	BooksServerPort          string `envconfig:"books_server_port"`
	ShopsServerPort          string `envconfig:"shops_server_port"`
	TransactionsServerPort   string `envconfig:"transactions_server_port"`
	PostgresBooks            string `envconfig:"postgres_books"`
	PostgresUser             string `envconfig:"postgres_user"`
	PostgresPassword         string `envconfig:"postgres_password"`
	PostgresBooksDB          string `envconfig:"postgres_books_db"`
	PostgresBooksPort        string `envconfig:"postgres_books_port"`
	PostgresBooksContainerIp string `envconfig:"postgres_books_container_ip"`
	SecretKey                string `envconfig:"secret_key"`
	Network                  string `envconfig:"network"`
	AutoSplitVar             string `split_words:"true"`
}
