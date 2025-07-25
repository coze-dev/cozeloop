// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package errorx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewByCode(t *testing.T) {
	type testCase struct {
		name     string
		code     int32
		options  []Option
		expected int32
	}

	tests := []testCase{
		{
			name:     "should create error with code",
			code:     1001,
			options:  []Option{WithExtraMsg("test")},
			expected: 1001,
		},
		{
			name:     "should create error without options",
			code:     1002,
			options:  nil,
			expected: 1002,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewByCode(tt.code, tt.options...)
			assert.NotNil(t, err)

			statusErr, ok := FromStatusError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.expected, statusErr.Code())
		})
	}
}

func TestNew(t *testing.T) {
	type testCase struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}

	tests := []testCase{
		{
			name:     "should create error with format",
			format:   "test error: %s",
			args:     []interface{}{"message"},
			expected: "test error: message",
		},
		{
			name:     "should create error without args",
			format:   "test error",
			args:     nil,
			expected: "test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.format, tt.args...)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), tt.expected)
		})
	}
}

func TestWrapByCode(t *testing.T) {
	type testCase struct {
		name     string
		err      error
		code     int32
		options  []Option
		expected int32
	}

	tests := []testCase{
		{
			name:     "should wrap error with code",
			err:      errors.New("test error"),
			code:     1001,
			options:  []Option{WithExtraMsg("wrapped")},
			expected: 1001,
		},
		{
			name:     "should return nil for nil error",
			err:      nil,
			code:     1001,
			options:  nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WrapByCode(tt.err, tt.code, tt.options...)
			if tt.err == nil {
				assert.Nil(t, err)
				return
			}

			assert.NotNil(t, err)
			statusErr, ok := FromStatusError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.expected, statusErr.Code())
		})
	}
}

func TestWrapf(t *testing.T) {
	type testCase struct {
		name     string
		err      error
		format   string
		args     []interface{}
		expected string
	}

	tests := []testCase{
		{
			name:     "should wrap error with format",
			err:      errors.New("test error"),
			format:   "wrapped: %s",
			args:     []interface{}{"message"},
			expected: "wrapped: message",
		},
		{
			name:     "should return nil for nil error",
			err:      nil,
			format:   "wrapped: %s",
			args:     []interface{}{"message"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Wrapf(tt.err, tt.format, tt.args...)
			if tt.err == nil {
				assert.Nil(t, err)
				return
			}

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), tt.expected)
		})
	}
}

func TestFromStatusError(t *testing.T) {
	type testCase struct {
		name     string
		err      error
		expected bool
	}

	tests := []testCase{
		{
			name:     "should convert StatusError",
			err:      NewByCode(1001),
			expected: true,
		},
		{
			name:     "should not convert standard error",
			err:      errors.New("test error"),
			expected: false,
		},
		{
			name:     "should handle nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusErr, ok := FromStatusError(tt.err)
			assert.Equal(t, tt.expected, ok)
			if tt.expected {
				assert.NotNil(t, statusErr)
				assert.Implements(t, (*StatusError)(nil), statusErr)
			}
		})
	}
}

func TestErrorWithoutStack(t *testing.T) {
	type testCase struct {
		name     string
		err      error
		expected string
	}

	tests := []testCase{
		{
			name:     "should remove stack trace",
			err:      NewByCode(1001),
			expected: "",
		},
		{
			name:     "should handle nil error",
			err:      nil,
			expected: "",
		},
		{
			name:     "should handle standard error",
			err:      errors.New("test error"),
			expected: "test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ErrorWithoutStack(tt.err)
			if tt.err == nil {
				assert.Empty(t, result)
				return
			}

			if tt.expected != "" {
				assert.Equal(t, tt.expected, result)
			} else {
				assert.NotEmpty(t, result)
				assert.NotContains(t, result, "stack=")
			}
		})
	}
}
