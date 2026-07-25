package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bruno1186/url-shortener/internal/server"
	"github.com/bruno1186/url-shortener/internal/shortener"
)

func main() {
	addr := getenv("ADDR", ":8080")
	baseURL := getenv("BASE_URL", "http://localhost"+addr)

	svc := shortener.New(shortener.NewMemoryStore())
	srv := server.New(svc, baseURL)

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      srv.Routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("url-shortener listening on %s (base url: %s)", addr, baseURL)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
