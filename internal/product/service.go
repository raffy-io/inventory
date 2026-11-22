package product

import (
	"context"

	inventoryDB "github.com/raffy-io/inventory/internal/db"
	"github.com/shopspring/decimal"
)

type productService struct {
	queries *inventoryDB.Queries
}

func newProductService(queries *inventoryDB.Queries) *productService {
	return &productService{
		queries: queries,
	}
}

func (p *productService) getAll(ctx context.Context) ([]inventoryDB.Product, error) {
	return p.queries.ListProducts(ctx)
}

func (p *productService) getByID(ctx context.Context, id int64) (inventoryDB.Product, error) {
	return p.queries.GetProductByID(ctx, id)
}

func (p *productService) create(ctx context.Context, name string, price decimal.Decimal, quantity decimal.Decimal) (inventoryDB.Product, error) {
	return p.queries.CreateProduct(ctx, inventoryDB.CreateProductParams{
		Name:     name,
		Price:    price,
		Quantity: quantity,
	})
}

func (p *productService) update(ctx context.Context, id int64, name string, price decimal.Decimal, quantity decimal.Decimal) error {
	return p.queries.UpdateProduct(ctx, inventoryDB.UpdateProductParams{
		ID:       id,
		Name:     name,
		Price:    price,
		Quantity: quantity,
	})
}

func (p *productService) delete(ctx context.Context, id int64) error {
	return p.queries.DeleteProduct(ctx, id)
}
