package memory

import (
	"fmt"

	"github.com/defval/injector/testdata/order"
	"github.com/defval/injector/testdata/storage"
)

// NewOrderRepository
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

// OrderRepository
type OrderRepository struct {
	uuid map[string]order.Order
}

func (r *OrderRepository) Save(order *order.Order) (err error) {
	r.uuid[order.UUID] = *order
	return nil
}

func (r *OrderRepository) FindOne(options ...storage.Option) (_ *order.Order, err error) {
	var qb = NewQueryBuilder()

	for _, opt := range options {
		opt(qb)
	}

	o, exist := r.uuid[qb.uuid]
	if !exist {
		return nil, fmt.Errorf("order %s not found", qb.uuid)
	}

	return &o, nil
}
