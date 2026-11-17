package product

import (
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
)

type ProductForm struct {
	Name   string
	Price  decimal.Decimal
	Weight decimal.Decimal
	Errors map[string]string
}

func newProductForm(name string, price decimal.Decimal, weight decimal.Decimal) *ProductForm {
	return &ProductForm{
		Name:   strings.TrimSpace(name),
		Price:  price,
		Weight: weight,
		Errors: make(map[string]string),
	}
}

func (f *ProductForm) addError(field, message string) {
	if f.Errors == nil {
		f.Errors = make(map[string]string)
	}
	f.Errors[field] = message
}

func (f *ProductForm) validate() bool {

	if utf8.RuneCountInString(f.Name) < 3 || utf8.RuneCountInString(f.Name) > 60 {
		f.addError("name", "Product name must be between 3 and 60 characters.")
	}

	// Validate Price (must be greater than zero)
	if f.Price.LessThanOrEqual(decimal.Zero) {
		f.addError("price", "Price must be greater than zero.")
	}

	// Validate Weight (must not be negative)
	if f.Weight.IsNegative() {
		f.addError("weight", "Weight cannot be negative.")
	}

	return len(f.Errors) == 0
}
