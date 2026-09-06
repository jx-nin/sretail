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

func TestMemoryStoreList(t *testing.T) {
	seed := []product.Product{
		{
			ID:         "prod-1",
			Name:       "Mechanical Keyboard",
			PriceCents: 12999,
		},
		{
			ID:         "prod-2",
			Name:       "Ergonomic Mouse",
			PriceCents: 10999,
		},
	}

	store := product.NewMemoryStore(seed)

	seed[0].Name = "Modified Keyboard"

	got, err := store.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("List() returned %d products, want 2", len(got))
	}

	want := []product.Product{
		{
			ID:         "prod-1",
			Name:       "Mechanical Keyboard",
			PriceCents: 12999,
		},
		{
			ID:         "prod-2",
			Name:       "Ergonomic Mouse",
			PriceCents: 10999,
		},
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("List()[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}

	got[0].Name = "Changed Through Result"

	gotAgain, err := store.List(context.Background())
	if err != nil {
		t.Fatalf("second List() error = %v", err)
	}

	if gotAgain[0] != want[0] {
		t.Errorf("second List()[0] = %#v, want %#v", gotAgain[0], want[0])
	}
}

func TestMemoryStoreListEmpty(t *testing.T) {
	store := product.NewMemoryStore(nil)

	got, err := store.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if got == nil {
		t.Fatal("List() returned nil, want an empty slice")
	}

	if len(got) != 0 {
		t.Errorf("List() returned %d products, want 0", len(got))
	}
}
