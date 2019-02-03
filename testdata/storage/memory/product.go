package memory

import "github.com/defval/injector/testdata/product"

// NewProductRepository
func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		uuid: make(map[string]product.Product),
	}
}

// ProductRepository
type ProductRepository struct {
	uuid map[string]product.Product
}

// Save
func (r *ProductRepository) Save(p *product.Product) (err error) {
	r.uuid[p.UUID] = *p
	return nil
}
