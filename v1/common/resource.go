package common

// Resource is a combination of id and name.
type Resource struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// GetID returns the Identifier of the [Resource].
func (r Resource) GetID() string {
	return r.Identifier
}

// IsEngineIdentifier checks if the given string s looks like an Anexia Engine Identifier.
// Does not very if it is a valid resource.
func IsEngineIdentifier(s string) bool {
	const identifierLength = 32

	if len(s) != identifierLength {
		return false
	}

	for _, char := range s {
		if (char < '0' || char > '9') && (char < 'a' || char > 'f') && (char < 'A' || char > 'F') {
			return false
		}
	}

	return true
}
