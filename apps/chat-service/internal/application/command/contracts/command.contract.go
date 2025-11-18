package contracts

// Command represents a command interface
type Command interface {
	Execute(ctx interface{}, req interface{}) (interface{}, error)
}

