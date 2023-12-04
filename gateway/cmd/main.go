package main

import (
	"library/gateway/server"
	"os"
)

func main() {
	// TODO: 1. add config package, with global config state, using library such as https://github.com/kelseyhightower/envconfig would be enough
	// TODO: 2. add global logger
	// TODO: 3. pass context
	gateway.Run(":" + os.Getenv("GATEWAY_SERVER_PORT"))
}
