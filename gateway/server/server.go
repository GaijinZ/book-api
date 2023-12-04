package gateway

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func Run(port string) {
	usersURL, _ := url.Parse("http://localhost" + ":" + os.Getenv("USERS_SERVER_PORT"))
	usersProxy := httputil.NewSingleHostReverseProxy(usersURL)

	booksURL, _ := url.Parse("http://localhost" + ":" + os.Getenv("BOOKS_SERVER_PORT"))
	booksProxy := httputil.NewSingleHostReverseProxy(booksURL)

	shopsURL, _ := url.Parse("http://localhost" + ":" + os.Getenv("SHOPS_SERVER_PORT"))
	shopsProxy := httputil.NewSingleHostReverseProxy(shopsURL)

	transactionsURL, _ := url.Parse("http://localhost" + ":" + os.Getenv("TRANSACTIONS_SERVER_PORT"))
	transactionsProxy := httputil.NewSingleHostReverseProxy(transactionsURL)

	http.HandleFunc("/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		usersProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/v1/books/", func(w http.ResponseWriter, r *http.Request) {
		booksProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/v1/shops/", func(w http.ResponseWriter, r *http.Request) {
		shopsProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/v1/transactions/", func(w http.ResponseWriter, r *http.Request) {
		transactionsProxy.ServeHTTP(w, r)
	})

	log.Println("API Gateway server is starting on :8080")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
