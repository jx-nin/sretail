package product

import "context"

type MemoryStore struct {
	products []Product
	byID     map[string]Product
}

var _ Store = (*MemoryStore)(nil)

func NewMemoryStore(products []Product) *MemoryStore {
	productCopy := make([]Product, len(products))
	copy(productCopy, products)

	byID := make(map[string]Product, len(productCopy))
	for _, product := range productCopy {
		byID[product.ID] = product
	}

	return &MemoryStore{
		products: productCopy,
		byID:     byID,
	}
}

func (s *MemoryStore) List(_ context.Context) ([]Product, error) {
	products := make([]Product, len(s.products))
	copy(products, s.products)

	return products, nil
}

func (s *MemoryStore) GetByID(_ context.Context, id string) (Product, error) {
	product, ok := s.byID[id]
	if !ok {
		return Product{}, ErrNotFound
	}

	return product, nil
}
