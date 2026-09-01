package common

import (
	"errors"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
)

// APIError represents an api level error.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

// Error implements the error interface for APIError.
func (a *APIError) Error() string {
	return fmt.Sprintf("api error: StatusCode=%d, Status=%s, Body=%s", a.StatusCode, a.Status, a.Body)
}

// MapTransportError maps an internal transport error to an externally usable error.
func MapTransportError(err error) error {
	if err == nil {
		return nil
	}

	if te, ok := errors.AsType[*internal.TransportError](err); ok {
		return &APIError{
			StatusCode: te.StatusCode,
			Status:     te.Status,
			Body:       te.Body,
		}
	}

	return err
}
