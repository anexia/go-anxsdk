package internal

// AllAttributesWrapper is a wrapper to state all attributes on an object should be loaded.
type AllAttributesWrapper struct {
	Params any
}

// NewAllAttributesWrapper creates a new AllAttributesWrapper.
func NewAllAttributesWrapper(params any) AllAttributesWrapper {
	return AllAttributesWrapper{Params: params}
}
