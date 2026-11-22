package auth

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/alexedwards/scs/v2"
	inventoryDB "github.com/raffy-io/inventory/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type authService struct {
	queries *inventoryDB.Queries
	session *scs.SessionManager
}

func newAuthService(queries *inventoryDB.Queries, session *scs.SessionManager) *authService {
	return &authService{
		queries: queries,
		session: session,
	}
}

func (a *authService) serveLogin() {

}
func (a *authService) processLogin(ctx context.Context, email, password string) (inventoryDB.User, error) {
	user, err := a.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return inventoryDB.User{}, ErrInvalidCredentials
		}
		return inventoryDB.User{}, err
	}

	if !CheckPassword(password, user.PasswordHash) {
		return inventoryDB.User{}, ErrInvalidCredentials
	}

	return user, nil
}

// HashPassword takes a plain-text password and returns its bcrypt hash.
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword compares a plain-text password against a hashed password.
// Returns true if they match, false otherwise.
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		slog.Info("Password check failed!", "error", err)
	} else {
		slog.Info("Password matches perfectly!")
	}
	return err == nil
}
