package store

import (
	"context"
	"errors"

	"pos/backend/internal/model"
)

var (
	ErrProductNotFound    = errors.New("product not found")
	ErrInsufficientStock  = errors.New("stock adjustment would make stock negative")
	ErrInvalidStockChange = errors.New("stock change must not be zero")
)

type ProductStore interface {
	ListProducts(context.Context, bool) ([]model.Product, error)
	GetProduct(context.Context, string) (model.Product, error)
	CreateProduct(context.Context, model.Product) error
	UpdateProduct(context.Context, model.Product) error
	SetProductActive(context.Context, string, bool) error
	AdjustStock(context.Context, string, int, string, string) (model.Product, error)
}
