package product_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jx-nin/sretail/services/catalog/internal/product"
)

func TestMemoryStoreGetByID(t *testing.T) {
	seed := []product.Product{
		{
			ID:         "prod-1",
			Name:       "Mechanical Keyboard",
			PriceCents: 12999,
		},
	}

	store := product.NewMemoryStore(seed)

	t.Run("returns an existing product", func(t *testing.T) {
		got, err := store.GetByID(context.Background(), "prod-1")
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if got != seed[0] {
			t.Errorf("GetByID() = %#v, want %#v", got, seed[0])
		}
	})

	t.Run("returns ErrNotFound for an unknown ID", func(t *testing.T) {
		_, err := store.GetByID(context.Background(), "missing")
		if !errors.Is(err, product.ErrNotFound) {
			t.Errorf("GetByID() error = %v, want ErrNotFound", err)
		}
	})
}
