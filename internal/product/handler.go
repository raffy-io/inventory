package product

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"
)

type ProductHandler struct {
	service *ProductService
}

func newProductHandler(service *ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

const (
	// PgErrUniqueViolation is the PostgreSQL error code for unique key constraint violations.
	PgErrUniqueViolation = "23505"
)

func (h *ProductHandler) serveItemsPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	items, err := h.service.listItemsService(ctx)
	if err != nil {
		http.Error(w, "Failed to get items", http.StatusInternalServerError)
		return
	}

	// 1. Tell Postman/clients that we are sending JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 2. Encode the sqlc generated structs straight to the response body
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, "Failed to encode items", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) createItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.PostFormValue("name")
	priceStr := r.PostFormValue("price")
	weightStr := r.PostFormValue("weight")

	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	weight, err := decimal.NewFromString(weightStr)
	if err != nil {
		http.Error(w, "Invalid weight format", http.StatusBadRequest)
		return
	}

	form := newProductForm(name, price, weight)
	if !form.validate() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"message": "Validation error"}`))
		return
	}

	item, err := h.service.createItemService(ctx, name, price, weight)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == PgErrUniqueViolation {

			form.addError("name", "An item with this name already exists.")
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(`{"message": "Unique constraint error"}`))
			return
			
		}

		http.Error(w, "Failed to add item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, "Failed to encode items", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) updateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	priceStr := r.FormValue("price")
	weightStr := r.FormValue("weight")

	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	weight, err := decimal.NewFromString(weightStr)
	if err != nil {
		http.Error(w, "Invalid weight format", http.StatusBadRequest)
		return
	}

	form := newProductForm(name, price, weight)
	if !form.validate() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"message": "Validation error"}`))
		return
	}

	err = h.service.updateItemService(ctx, id, name, price, weight)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == PgErrUniqueViolation {
			
			form.addError("name", "An item with this name already exists.")
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(`{"message": "Unique constraint error"}`))
			return

		}

		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Item updated successfully"}`))

}

func (h *ProductHandler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	err = h.service.deleteItemService(ctx, id)
	if err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Item deleted successfully"}`))
}
