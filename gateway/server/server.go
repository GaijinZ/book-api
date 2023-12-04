package gateway

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func Run(port string) {
	// TODO: 1. move all env parsing to global config
	// TODO: 2. error handling
	usersURL, _ := url.Parse("http://localhost" + ":" + os.Getenv("USERS_SERVER_PORT"))
	usersProxy := httputil.NewSingleHostReverseProxy(usersURL)

	// TODO: same here
	booksURL, _ := url.Parse("http://localhost" + ":" + os.Getenv("BOOKS_SERVER_PORT"))
	booksProxy := httputil.NewSingleHostReverseProxy(booksURL)

	// TODO: same here
	shopsURL, _ := url.Parse("http://localhost" + ":" + os.Getenv("SHOPS_SERVER_PORT"))
	shopsProxy := httputil.NewSingleHostReverseProxy(shopsURL)

	// TODO: same here
	transactionsURL, _ := url.Parse("http://localhost" + ":" + os.Getenv("TRANSACTIONS_SERVER_PORT"))
	transactionsProxy := httputil.NewSingleHostReverseProxy(transactionsURL)

	// TODO: 1. add middleware, for parsing traceID, user authentication/authorization
	// TODO: 2. move api version, to upper handler
	http.HandleFunc("/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		usersProxy.ServeHTTP(w, r)
	})

	// TODO: same here
	http.HandleFunc("/v1/books/", func(w http.ResponseWriter, r *http.Request) {
		booksProxy.ServeHTTP(w, r)
	})

	// TODO: same here
	http.HandleFunc("/v1/shops/", func(w http.ResponseWriter, r *http.Request) {
		shopsProxy.ServeHTTP(w, r)
	})

	// TODO: same here
	http.HandleFunc("/v1/transactions/", func(w http.ResponseWriter, r *http.Request) {
		transactionsProxy.ServeHTTP(w, r)
	})

	// TODO: don't hardcode port
	log.Println("API Gateway server is starting on :8080")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
