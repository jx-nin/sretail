package product

import (
	"context"
	"errors"
)

var (
	ErrNotFound    = errors.New("product not found")
	ErrDuplicateID = errors.New("duplicate product ID")
)

type Product struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PriceCents int64  `json:"price_cents"`
}

type Store interface {
	List(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, id string) (Product, error)
}
