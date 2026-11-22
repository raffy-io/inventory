package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type authHandler struct {
	service *authService
	session *scs.SessionManager
}

func newAuthHandler(service *authService, session *scs.SessionManager) *authHandler {
	return &authHandler{
		service: service,
		session: session,
	}
}

func (h *authHandler) serveLogin(w http.ResponseWriter, r *http.Request) {}

func (h *authHandler) processLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	form := newAuthForm(email, password)
	if !form.validate() {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": "validation error",
		})
		return
	}

	user, err := h.service.processLogin(ctx, form.email, form.password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"message": "Invalid email or password",
			})
			return
		}
		slog.Error("login failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})
		return
	}

	if err := h.session.RenewToken(ctx); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"message": "failed to renew token",
		})
		return
	}
	// 4. Put user_id into session
	h.session.Put(ctx, "user_id", user.ID)
	h.session.Put(ctx, "role", user.Role)

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Login successful",
	})
}

func (h *authHandler) processLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Destroy the session and clear the cookie
	if err := h.session.Destroy(ctx); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "failed to destroy session",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "successfully logged out",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
