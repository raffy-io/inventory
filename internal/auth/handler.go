package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/raffy-io/inventory/internal/web"
)

type authHandler struct {
	service *authService
	session *scs.SessionManager
	render  *web.Renderer
}

func newAuthHandler(service *authService, session *scs.SessionManager, render *web.Renderer) *authHandler {
	return &authHandler{
		service: service,
		session: session,
		render:  render,
	}
}

func (h *authHandler) serveLogin(w http.ResponseWriter, r *http.Request) {
	form := newAuthForm("", "")
	component := login(form)

	if err := h.render.Render(w, r, "Login", component); err != nil {
		slog.Error("failed to render login template", "err", err, "path", r.URL.Path)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *authHandler) processLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	form := newAuthForm(email, password)
	if !form.validate() {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := loginForm(form).Render(ctx, w); err != nil {
			slog.Error("failed to render login form", "err", err, "path", r.URL.Path)
		}
		return
	}

	user, err := h.service.processLogin(ctx, form.email, form.password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			form.addError("email", "invalid email or password")
			form.addError("password", "invalid email or password")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := loginForm(form).Render(ctx, w); err != nil {
				slog.Error("failed to render login form", "err", err, "path", r.URL.Path)
			}
			return
		}
		slog.Error("login failed", "err", err, "path", r.URL.Path)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := h.session.RenewToken(ctx); err != nil {
		slog.Error("failed to renew session token", "err", err, "path", r.URL.Path)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.session.Put(ctx, "user_id", user.ID)
	h.session.Put(ctx, "role", user.Role)

	w.Header().Set("HX-Location", "/products")
	w.WriteHeader(http.StatusOK)
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
