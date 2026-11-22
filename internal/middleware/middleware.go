package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
)

// Session keys and roles: one place, no magic strings.
const (
	SessionUserID = "user_id"
	SessionRole   = "role"
	RoleAdmin     = "admin"
)

type contextKey string

const userIDKey contextKey = "user_id"

func GetUserID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}

// Auth holds what changes between projects.
type Auth struct {
	Sessions  *scs.SessionManager
	LoginPath string
	HomePath  string
}

func (a Auth) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.Sessions.Exists(r.Context(), SessionUserID) {
			slog.Info("unauthenticated request", "path", r.URL.Path)
			redirect(w, r, a.LoginPath)
			return
		}

		userID := a.Sessions.GetInt64(r.Context(), SessionUserID)
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin must run AFTER RequireAuth.
func (a Auth) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.Sessions.GetString(r.Context(), SessionRole) != RoleAdmin {
			slog.Warn("unauthorized admin request", "path", r.URL.Path)
			redirect(w, r, a.HomePath)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// redirect handles both HTMX and normal requests in one place.
func redirect(w http.ResponseWriter, r *http.Request, path string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", path)
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func NoSurf(secure bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := nosurf.New(next)
		h.SetBaseCookie(http.Cookie{
			Path:     "/",
			HttpOnly: true,
			Secure:   secure,
			SameSite: http.SameSiteLaxMode,
		})
		h.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Error("bad or missing CSRF token", "path", r.URL.Path, "method", r.Method)
			http.Error(w, "Bad or missing CSRF token", http.StatusForbidden)
		}))
		return h
	}
}
