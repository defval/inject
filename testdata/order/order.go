package order

import (
	"github.com/defval/injector/testdata/storage"
)

// Order
type Order struct {
	ID       int64
	UUID     string
	Products []interface{}
	Amount   float64
}

// Save
func (o *Order) Save(repository Repository) (err error) {
	return repository.Save(o)
}

// Repository
type Repository interface {
	Save(order *Order) (err error)
	FindOne(options ...storage.Option) (_ *Order, err error)
}

// ByUUID
func ByUUID(uuid string) storage.Option {
	return func(qb storage.QueryBuilder) {
		qb.UUID(uuid)
	}
}
