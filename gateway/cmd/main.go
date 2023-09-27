package main

import (
	"library/gateway/server"
	"os"
)

func main() {
	gateway.Run(":" + os.Getenv("GATEWAY_SERVER_PORT"))
}
