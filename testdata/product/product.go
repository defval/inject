package product

// Product
type Product struct {
	ID    int64
	UUID  string
	Name  string
	Price float64
}

func (p *Product) Save(repository Repository) (err error) {
	return repository.Save(p)
}

// Repository
type Repository interface {
	Save(product *Product) (err error)
}
