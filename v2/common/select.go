package common

// IDTitleTuple is an ID to Title tuple, which is used for selects.
type IDTitleTuple[T comparable] struct {
	ID    T      `json:"id"`
	Title string `json:"title"`
}
