package auth

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	inventoryDB "github.com/raffy-io/inventory/internal/db"
)

func AuthRoutes(queries *inventoryDB.Queries, session *scs.SessionManager) http.Handler {
	r := chi.NewRouter()
	s := newAuthService(queries, session)
	h := newAuthHandler(s, session)

	r.Get("/login", h.serveLogin)
	r.Post("/login", h.processLogin)
	r.Delete("/logout", h.processLogout)

	return r
}
