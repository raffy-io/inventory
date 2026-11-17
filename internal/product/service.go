package product

import (
	"context"

	inventoryDB "github.com/raffy-io/inventory/internal/db"
	"github.com/shopspring/decimal"
)

type ProductService struct {
	queries *inventoryDB.Queries
}

func newProductService(queries *inventoryDB.Queries) *ProductService {
	return &ProductService{
		queries: queries,
	}
}

func (p *ProductService) listItemsService(ctx context.Context) ([]inventoryDB.Item, error) {
	return p.queries.ListItems(ctx)
}

func (p *ProductService) createItemService(ctx context.Context, name string, price decimal.Decimal, weight decimal.Decimal) (inventoryDB.Item, error) {

	return p.queries.CreateItem(ctx, inventoryDB.CreateItemParams{
		Name:   name,
		Price:  price,
		Weight: weight,
	})

}

func (p *ProductService) updateItemService(ctx context.Context, id int64, name string, price decimal.Decimal, weight decimal.Decimal) error {

	return p.queries.UpdateItem(ctx, inventoryDB.UpdateItemParams{
		ID:     id,
		Name:   name,
		Price:  price,
		Weight: weight,
	})
}

func (p *ProductService) deleteItemService(ctx context.Context, id int64) error {

	return p.queries.DeleteItem(ctx, id)

}
