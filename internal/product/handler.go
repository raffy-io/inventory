package product

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
	inventoryDB "github.com/raffy-io/inventory/internal/db"
	"github.com/raffy-io/inventory/internal/web"
	"github.com/shopspring/decimal"
)

type productHandler struct {
	service *productService
	render  *web.Renderer
}

func newProductHandler(service *productService, render *web.Renderer) *productHandler {
	return &productHandler{
		service: service,
		render:  render,
	}
}

const (
	// PgErrUniqueViolation is the PostgreSQL error code for unique key constraint violations.
	PgErrUniqueViolation = "23505"
)

func (h *productHandler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, err := h.service.getAll(ctx)

	if err != nil {
		http.Error(w, "Failed to get products data..", http.StatusInternalServerError)
		return
	}
	if products == nil {
		products = []inventoryDB.Product{}
	}

	component := productsPage()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.render.Render(w, r, "Welcome", component); err != nil {
		slog.Error("failed to render login template", "err", err, "path", r.URL.Path)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

func (h *productHandler) getByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid product ID",
		})
		return
	}

	product, err := h.service.getByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"message": "Product not found",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})
		return
	}

	writeJSON(w, http.StatusOK, product)
}

func (h *productHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.PostFormValue("name")
	priceStr := r.PostFormValue("price")
	quantityStr := r.PostFormValue("quantity")

	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid price format",
		})
		return
	}
	quantity, err := decimal.NewFromString(quantityStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid quantity format",
		})
		return
	}

	// validate
	form := newProductForm(name, price, quantity)

	if !form.validate() {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Validation error",
		})
		return
	}

	product, err := h.service.create(ctx, form.name, form.price, form.quantity)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == PgErrUniqueViolation {
			form.addError("name", "An item with this name already exists.")
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"message": "Unique constraint error",
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Failed to add product",
		})
		return
	}

	writeJSON(w, http.StatusCreated, product)
}

func (h *productHandler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid product ID",
		})
		return
	}

	if err := r.ParseForm(); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Failed to parse form",
		})
		return
	}

	name := r.PostFormValue("name")
	priceStr := r.PostFormValue("price")
	quantityStr := r.PostFormValue("quantity")

	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid price format",
		})
		return
	}
	quantity, err := decimal.NewFromString(quantityStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid quantity format",
		})
		return
	}

	// validate
	form := newProductForm(name, price, quantity)

	if !form.validate() {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Validation error",
		})
		return
	}

	err = h.service.update(ctx, id, form.name, form.price, form.quantity)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Failed to update product",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Successfully updated product",
	})
}

func (h *productHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Invalid product ID",
		})
		return
	}

	err = h.service.delete(ctx, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Failed to delete product",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Successfully deleted product",
	})

}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
