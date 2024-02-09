package resp_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aelnahas/sider/resp"
	"github.com/stretchr/testify/assert"
)

func TestEncoder(t *testing.T) {
	tests := []struct {
		name string
		data any

		expected []byte
	}{
		{
			name:     "string",
			data:     "ok",
			expected: []byte("+ok\r\n"),
		},
		{
			name:     "int",
			data:     100,
			expected: []byte(":100\r\n"),
		},
		{
			name:     "error",
			data:     errors.New("my error"),
			expected: []byte("-my error\r\n"),
		},
		{
			name:     "nil",
			data:     nil,
			expected: []byte("$-1\r\n"),
		},
		{
			name:     "unexpected type",
			data:     true,
			expected: []byte("-unknown response type bool\r\n"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			res := resp.Encode(tc.data)
			assert.Equal(t, tc.expected, res, fmt.Sprintf("received %s", string(res)))

		})
	}
}
