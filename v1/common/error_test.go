package common_test

import (
	"net/http"
	"testing"

	"github.com/anexia/go-anxsdk/v1/common"
	"github.com/stretchr/testify/assert"
)

func TestIsNotFoundError(t *testing.T) {
	assert.False(t, common.IsNotFoundError(nil))
	assert.True(t, common.IsNotFoundError(&common.APIError{
		StatusCode: http.StatusNotFound,
	}))
	assert.False(t, common.IsNotFoundError(&common.APIError{
		StatusCode: http.StatusBadRequest,
	}))
}

func TestIsErrorWithStatusCode(t *testing.T) {
	err := &common.APIError{
		StatusCode: http.StatusNotFound,
	}

	assert.True(t, common.IsErrorWithStatusCode(err, http.StatusNotFound))
	assert.False(t, common.IsErrorWithStatusCode(err, http.StatusBadRequest))
	assert.False(t, common.IsErrorWithStatusCode(nil, http.StatusNotFound))
}
