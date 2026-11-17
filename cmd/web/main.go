package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	inventoryDB "github.com/raffy-io/inventory/internal/db"
	"github.com/raffy-io/inventory/internal/product"
)

func main() {
	ctx := context.Background()
	_ = godotenv.Load()

	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		slog.Error("DB_URL environment variable is required")
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, DB_URL)
	if err != nil {
		slog.Error("Failed to create connection pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	queries := inventoryDB.New(pool)

	r := chi.NewRouter()

	r.Mount("/products", product.ProductsRoutes(queries))

	port := os.Getenv("PORT")
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	slog.Info("Server is running", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
