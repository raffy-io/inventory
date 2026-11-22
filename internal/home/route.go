package home

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/raffy-io/inventory/internal/web"
)

func RootRoute(session *scs.SessionManager) http.Handler{
	r := chi.NewRouter()
	renderer := web.NewRenderer(session)
	h := newHomeHandler(renderer)
	
	r.Get("/",h.serveHomepage)
	return r
}