package main

import (
	gateway "library/internal/gateway/server"
	"os"
)

func main() {
	gateway.Run(":" + os.Getenv("GATEWAY_SERVER_PORT"))
}
