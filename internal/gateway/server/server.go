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

	http.HandleFunc("/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		usersProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/v1/books/", func(w http.ResponseWriter, r *http.Request) {
		booksProxy.ServeHTTP(w, r)
	})

	log.Println("API Gateway server is starting on :8080")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
