package contracts

// Query represents a query interface
type Query interface {
	Execute(ctx interface{}, req interface{}) (interface{}, error)
}

