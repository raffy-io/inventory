package auth

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	inventoryDB "github.com/raffy-io/inventory/internal/db"
	"github.com/raffy-io/inventory/internal/web"
)

func AuthRoutes(queries *inventoryDB.Queries, session *scs.SessionManager) http.Handler {
	r := chi.NewRouter()
	render := web.NewRenderer(session)
	s := newAuthService(queries, session)
	h := newAuthHandler(s, session, render)

	r.Get("/login", h.serveLogin)
	r.Post("/login", h.processLogin)
	r.Delete("/logout", h.processLogout)

	return r
}
