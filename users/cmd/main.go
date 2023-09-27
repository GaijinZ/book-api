package main

import (
	"library/users/server"
	"os"
)

func main() {
	server.Run(":" + os.Getenv("USERS_SERVER_PORT"))
}
