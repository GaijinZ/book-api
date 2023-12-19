package gateway

import (
	"context"
	"library/pkg/config"
	"library/pkg/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Run(ctx context.Context, cfg config.GlobalEnv, port string) {
	log := utils.GetLogger(ctx)

	usersURL, err := url.Parse("http://localhost/v1" + ":" + cfg.UsersServerPort)
	if err != nil {
		log.Errorf("Failed to parse user server: %d", err)
	}
	usersProxy := httputil.NewSingleHostReverseProxy(usersURL)

	booksURL, err := url.Parse("http://localhost/v1" + ":" + cfg.BooksServerPort)
	if err != nil {
		log.Errorf("Failed to parse book server: %d", err)
	}
	booksProxy := httputil.NewSingleHostReverseProxy(booksURL)

	shopsURL, err := url.Parse("http://localhost/v1" + ":" + cfg.ShopsServerPort)
	if err != nil {
		log.Errorf("Failed to parse shops server: %d", err)
	}
	shopsProxy := httputil.NewSingleHostReverseProxy(shopsURL)

	transactionsURL, err := url.Parse("http://localhost/v1" + ":" + cfg.TransactionsServerPort)
	if err != nil {
		log.Errorf("Failed to parse transaction server: %d", err)
	}
	transactionsProxy := httputil.NewSingleHostReverseProxy(transactionsURL)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		usersProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		booksProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/shops", func(w http.ResponseWriter, r *http.Request) {
		shopsProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		transactionsProxy.ServeHTTP(w, r)
	})

	log.Infof("API Gateway server is starting on :%s", cfg.GatewayServerPort)

	if err = http.ListenAndServe(port, nil); err != nil {
		log.Errorf("Failed to start API Gateway: %v", err)
	}
}
