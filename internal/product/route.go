package product

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	inventoryDB "github.com/raffy-io/inventory/internal/db"
	"github.com/raffy-io/inventory/internal/web"
)

func ProductRoutes(queries *inventoryDB.Queries, session *scs.SessionManager) http.Handler {
	r := chi.NewRouter()
	render := web.NewRenderer(session)
	s := newProductService(queries)
	h := newProductHandler(s, render)

	r.Get("/", h.getAll)
	r.Get("/{id}", h.getByID)
	r.Post("/", h.Create)
	r.Patch("/{id}", h.update)
	r.Delete("/{id}", h.delete)

	return r
}
