package main

import (
	"library/shops/server"
	"os"
)

func main() {
	server.Run(":" + os.Getenv("SHOPS_SERVER_PORT"))
}
