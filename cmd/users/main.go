package main

import (
	"os"
	"userapi/internal/users/server"
)

func main() {
	server.Run(os.Getenv("USERS_SERVER_PORT"))
}
