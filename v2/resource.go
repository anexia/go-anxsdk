package v2

// Resource is a combination of id and name.
type Resource struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// GetID returns the Identifier of the [Resource].
func (r Resource) GetID() string {
	return r.Identifier
}

// IDTitleTuple is an ID to Title tuple, which is used for selects.
type IDTitleTuple[T comparable] struct {
	ID    T      `json:"id"`
	Title string `json:"title"`
}
