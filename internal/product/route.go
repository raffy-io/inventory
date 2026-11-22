package product

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	inventoryDB "github.com/raffy-io/inventory/internal/db"
)

func ProductRoutes(queries *inventoryDB.Queries) http.Handler {
	r := chi.NewRouter()
	s := newProductService(queries)
	h := newProductHandler(s)

	r.Get("/", h.getAll)
	r.Get("/{id}", h.getByID)
	r.Post("/", h.Create)
	r.Patch("/{id}", h.update)
	r.Delete("/{id}", h.delete)

	return r
}
