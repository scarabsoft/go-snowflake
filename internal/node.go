package internal

//NodeIDProvider provides an ID of a given node
type NodeIDProvider interface {
	// Returns the ID of the given node, where the ID gets generated. The implementation must provide unique IDs for
	// each instance, otherwise it can not be guaranteed that generated IDs are unique
	// Only invoked once
	ID() (uint8, error)
}

type fixedNodeIdProviderImpl struct {
	id uint8
}

func (f fixedNodeIdProviderImpl) ID() (uint8, error) {
	return f.id, nil
}

type Result struct {
	ID    uint64
	Error error
}

func NewFixedNodeIdProvider(id uint8) *fixedNodeIdProviderImpl {
	return &fixedNodeIdProviderImpl{id}
}
