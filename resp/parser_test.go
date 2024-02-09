package resp_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/aelnahas/sider/resp"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name  string
		input string

		expectedError error
		expectedAST   *resp.RawCommand
	}{
		{
			name:          "get",
			input:         "*2\r\n$3\r\nget\r\n$3\r\nfoo\r\n",
			expectedError: nil,
			expectedAST: &resp.RawCommand{
				Name:    "GET",
				Args:    []string{"foo"},
				Options: map[string]any{},
			},
		},
		{
			name:          "set",
			input:         "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			expectedError: nil,
			expectedAST: &resp.RawCommand{
				Name:    "SET",
				Args:    []string{"foo", "bar"},
				Options: map[string]any{},
			},
		},
		{
			name:          "set with option",
			input:         "*5\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$2\r\nEX\r\n$3\r\n100\r\n",
			expectedError: nil,
			expectedAST: &resp.RawCommand{
				Name: "SET",
				Args: []string{"foo", "bar"},
				Options: map[string]any{
					"EX": time.Duration(100000000000),
				},
				IsPubSubCMD: false,
			},
		},
		{
			name:          "ping without echo",
			input:         "*1\r\n$4\r\nPING\r\n",
			expectedError: nil,
			expectedAST: &resp.RawCommand{
				Name:    "PING",
				Args:    []string{},
				Options: map[string]any{},
			},
		},
		{
			name:          "ping with echo",
			input:         "*2\r\n$4\r\nPING\r\n$3\r\nfoo\r\n",
			expectedError: nil,
			expectedAST: &resp.RawCommand{
				Name:    "PING",
				Args:    []string{"foo"},
				Options: map[string]any{},
			},
		},
		{
			name:          "ping with echo",
			input:         "*2\r\n$4\r\nPING\r\n$3\r\nfoo\r\n",
			expectedError: nil,
			expectedAST: &resp.RawCommand{
				Name:    "PING",
				Args:    []string{"foo"},
				Options: map[string]any{},
			},
		},
		{
			name:          "invalid command",
			input:         "*1\r\n$3\r\nLOL\r\n",
			expectedError: resp.ErrUnknownCommand{Name: "LOL"},
			expectedAST:   nil,
		},
		{
			name:          "get given without options",
			input:         "*1\r\n$3\r\nGET\r\n",
			expectedError: errors.New("syntax err command GET is missing required args"),
			expectedAST:   nil,
		},
		{
			name:          "set with missing option value",
			input:         "*4\r\n$3\r\nset\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$2\r\nex\r\n",
			expectedError: errors.New("syntax err option EX is missing a value"),
			expectedAST:   nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			input := bytes.NewBufferString(tc.input)

			ast, err := resp.Parse(input)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedAST, ast)
		})
	}
}
