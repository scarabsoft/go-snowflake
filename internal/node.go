package internal

type NodeProvider interface {
	// Returns the ID of the given host, where the ID gets generated. The implementation must provide unique IDs for
	// each instance, otherwise it can not be guaranteed that generated IDs are unique
	// Only invoked once
	ID() (uint8, error)
}

type fixedNodeProviderImpl struct {
	id uint8
}

func (f fixedNodeProviderImpl) ID() (uint8, error) {
	return f.id, nil
}

type Result struct {
	ID    uint64
	Error error
}

func NewFixedNodeProvider(id uint8) *fixedNodeProviderImpl {
	return &fixedNodeProviderImpl{id}
}
