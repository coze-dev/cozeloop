// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package confmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRetrieved(t *testing.T) {
	tests := []struct {
		name    string
		rawConf any
		err     bool
	}{
		{
			name:    "valid map config",
			rawConf: map[string]any{"key": "value"},
			err:     false,
		},
		{
			name:    "invalid config type",
			rawConf: make(chan int),
			err:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRetrieved(tt.rawConf)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDecodeConfig(t *testing.T) {
	tests := []struct {
		name   string
		input  any
		output any
		err    bool
	}{
		{
			name:   "valid config",
			input:  map[string]any{"field": "value"},
			output: &struct{ Field string }{},
			err:    false,
		},
		{
			name:   "invalid config",
			input:  "not a map",
			output: &struct{}{},
			err:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DecodeConfig(tt.input, tt.output)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
