package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/godIsAProgrammer/url-shortener-go/shortener"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8789"
	}

	store := shortener.NewStore()
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           shortener.NewMux(store),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("url-shortener listening on %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
