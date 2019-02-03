package memory

// NewQueryBuilder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

// QueryBuilder
type QueryBuilder struct {
	uuid string
}

// UUID
func (b *QueryBuilder) UUID(uuid string) {
	b.uuid = uuid
}
