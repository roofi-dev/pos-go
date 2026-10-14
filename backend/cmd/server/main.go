package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"pos/backend/internal/httpapi"
	"pos/backend/internal/store"
)

func main() {
	databaseURL, err := databaseURLFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	productStore, err := store.OpenPostgres(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer productStore.Close()

	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           httpapi.NewHandler(productStore, os.Getenv("ADMIN_USERNAME"), os.Getenv("ADMIN_PASSWORD")),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("POS API berjalan di http://%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func databaseURLFromEnv() (string, error) {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return "", errors.New("DATABASE_URL wajib diatur sebelum menjalankan POS API")
	}
	return databaseURL, nil
}
