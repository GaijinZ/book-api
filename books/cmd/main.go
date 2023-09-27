package main

import (
	"library/books/server"
	"os"
)

func main() {
	server.Run(":" + os.Getenv("BOOKS_SERVER_PORT"))
}
