package common

// State represents a resources current state.
type State[TState any] struct {
	Text  string `json:"text"`
	Title string `json:"title"`
	ID    TState `json:"id"`
	Type  int    `json:"type"`
}
