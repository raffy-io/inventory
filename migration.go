package inventory

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/pressly/goose/v3"
)

//go:embed database/migration/*.sql
var migrationFiles embed.FS

func RunMigrations(ctx context.Context, db *sql.DB) error {
	goose.SetBaseFS(migrationFiles)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.UpContext(ctx, db, "database/migration"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err := seedInitialOwner(ctx, db); err != nil {
		return fmt.Errorf("failed to seed owner credentials: %w", err)
	}

	return nil
}

func seedInitialOwner(ctx context.Context, db *sql.DB) error {
	ownerPasswordHash := strings.TrimSpace(os.Getenv("OWNER_PASSWORD_HASH"))
	if ownerPasswordHash == "" {
		slog.Warn("OWNER_PASSWORD_HASH not set; skipping initial owner seed")
		return nil
	}

	ownerUsername := strings.TrimSpace(os.Getenv("OWNER_USERNAME"))
	if ownerUsername == "" {
		ownerUsername = "owner"
	}

	ownerEmail := strings.TrimSpace(os.Getenv("OWNER_EMAIL"))
	if ownerEmail == "" {
		ownerEmail = "admin@inventory.local"
	}

	query := `
		INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, 'owner')
		ON CONFLICT (email) DO NOTHING
	`

	if _, err := db.ExecContext(ctx, query, ownerUsername, ownerEmail, ownerPasswordHash); err != nil {
		return err
	}

	slog.Info("ensured initial owner user exists", "email", ownerEmail)
	return nil
}
