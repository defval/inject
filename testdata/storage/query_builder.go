package storage

// QueryBuilder
type QueryBuilder interface {
	UUID(uuid string)
}

// Option
type Option func(qb QueryBuilder)
