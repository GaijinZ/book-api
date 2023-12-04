package main

import (
	"library/transactions/server"
	"os"
)

func main() {
	server.Run(":" + os.Getenv("TRANSACTIONS_SERVER_PORT"))
}
