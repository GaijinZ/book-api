package main

import (
	"library/internal/users/server"
	"os"
)

func main() {
	server.Run(":" + os.Getenv("USERS_SERVER_PORT"))
}
