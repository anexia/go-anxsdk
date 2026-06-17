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

func MapTransportError(err error) error {
	if err == nil {
		return nil
	}

	var te *internal.TransportError
	if errors.As(err, &te) {
		return &APIError{
			StatusCode: te.StatusCode,
			Status:     te.Status,
			Body:       te.Body,
		}
	}

	return err
}
