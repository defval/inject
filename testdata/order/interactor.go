package order

// NewInteractor
func NewInteractor(orders Repository) *Interactor {
	return &Interactor{
		orders: orders,
	}
}

// Interactor
type Interactor struct {
	orders Repository
}
