package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/raffy-io/inventory"
	"github.com/raffy-io/inventory/internal/auth"
	inventoryDB "github.com/raffy-io/inventory/internal/db"
	"github.com/raffy-io/inventory/internal/middleware"
	"github.com/raffy-io/inventory/internal/product"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	_ = godotenv.Load()

	// Database
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return errors.New("DB_URL environment variable is required")
	}

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	defer pool.Close()

	// Run migrations
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	if err := inventory.RunMigrations(ctx, sqlDB); err != nil {
		return fmt.Errorf("migration failed: %v", err)
	}
	slog.Info("Migrations applied successfully!")

	queries := inventoryDB.New(pool)

	// Sessions
	isProduction := os.Getenv("APP_ENV") == "production"
	store := pgxstore.NewWithCleanupInterval(pool, 5*time.Minute)
	defer store.StopCleanup()

	sessionManager := scs.New()
	sessionManager.Store = store
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = isProduction

	// Routes
	r := chi.NewRouter()

	// -----static assets CSS JS IMG----- //
	staticFS, err := fs.Sub(inventory.EmbeddedAssets, "static")
	if err != nil {
		return fmt.Errorf("Cant find static assets: %w", err)
	}
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Use(sessionManager.LoadAndSave)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/products", http.StatusFound)
		})
		r.Mount("/auth", auth.AuthRoutes(queries, sessionManager))

	})

	// Protected routes
	authMW := &middleware.Auth{
		Sessions:  sessionManager,
		LoginPath: "/auth/login",
		HomePath:  "/products",
	}
	r.Group(func(r chi.Router) {
		r.Use(sessionManager.LoadAndSave)
		r.Use(authMW.RequireAuth)
		r.Mount("/products", product.ProductRoutes(queries, sessionManager))
	})

	// Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		<-ctx.Done()
		slog.Info("Shutting down server gracefully")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server shutdown error", "error", err)
		}
	}()

	slog.Info("Server is running", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}
