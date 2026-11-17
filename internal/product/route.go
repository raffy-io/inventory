package product

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	inventoryDB "github.com/raffy-io/inventory/internal/db"
)

func ProductsRoutes(queries *inventoryDB.Queries) http.Handler {

	r := chi.NewRouter()
	service := newProductService(queries)
	h := newProductHandler(service)

	r.Get("/", h.serveItemsPage)
	r.Post("/", h.createItem)
	r.Patch("/{id}", h.updateProduct)
	r.Delete("/{id}", h.deleteProduct)

	return r
}
