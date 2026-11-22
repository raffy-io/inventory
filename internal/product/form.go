package product

import (
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
)

type productForm struct {
	name     string
	price    decimal.Decimal
	quantity decimal.Decimal
	Errors   map[string]string
}

func newProductForm(name string, price decimal.Decimal, quantity decimal.Decimal) *productForm {
	return &productForm{
		name:     strings.TrimSpace(name),
		price:    price,
		quantity: quantity,
	}
}

func (f *productForm) addError(field, message string) {
	if f.Errors == nil {
		f.Errors = make(map[string]string)
	}
	f.Errors[field] = message
}

func (f *productForm) validate() bool {
	if utf8.RuneCountInString(f.name) < 3 || utf8.RuneCountInString(f.name) > 60 {
		f.addError("name", "Product name must be between 3 and 60 characters.")
	}

	if f.price.LessThanOrEqual(decimal.Zero) {
		f.addError("price", "Price must be greater than zero.")
	}

	if f.price.LessThanOrEqual(decimal.Zero) {
		f.addError("quantity", "Quantity must be greater than zero.")
	}

	return len(f.Errors) == 0
}
