package v1

// State represents a resources current state.
type State struct {
	Text  string       `json:"text"`
	Title string       `json:"title"`
	ID    ClusterState `json:"id"`
	Type  int          `json:"type"`
}
