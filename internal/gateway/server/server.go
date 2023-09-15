package gateway

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Run(port string) {
	// Create a reverse proxy for the microservice (port 5000)
	usersURL, _ := url.Parse("http://localhost:5000")
	usersProxy := httputil.NewSingleHostReverseProxy(usersURL)

	// Create a reverse proxy for the microservice (port 5001)
	booksURL, _ := url.Parse("http://localhost:5001")
	booksProxy := httputil.NewSingleHostReverseProxy(booksURL)

	http.HandleFunc("/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		// Modify the request path to forward to the appropriate microservice
		usersProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/v1/books/", func(w http.ResponseWriter, r *http.Request) {
		// Modify the request path to forward to the appropriate microservice
		booksProxy.ServeHTTP(w, r)
	})

	log.Println("API Gateway server is starting on :8080")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
