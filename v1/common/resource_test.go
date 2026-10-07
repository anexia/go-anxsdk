package common_test

import (
	"testing"

	"github.com/anexia/go-anxsdk/v1/common"
	"github.com/stretchr/testify/assert"
)

func TestIsEngineIdentifier(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid lowercase hex",
			input: "0123456789abcdef0123456789abcdef",
			want:  true,
		},
		{
			name:  "too short",
			input: "0123456789abcdef0123456789abcde",
			want:  false,
		},
		{
			name:  "too long",
			input: "0123456789abcdef0123456789abcdef0",
			want:  false,
		},
		{
			name:  "uppercase hex",
			input: "0123456789ABCDEF0123456789abcdef",
			want:  true,
		},
		{
			name:  "invalid hex character",
			input: "0123456789abcdef0123456789abcdeg",
			want:  false,
		},
		{
			name:  "special character",
			input: "0123456789abcdef0123456789abcde!",
			want:  false,
		},
		{
			name:  "space",
			input: "0123456789abcdef0123456789ab12 0",
			want:  false,
		},
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := common.IsEngineIdentifier(tt.input)
			assert.Equal(t, got, tt.want)
		})
	}
}
