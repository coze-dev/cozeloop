// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package component

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockConfig struct{}

type mockConfigWithValidator struct{}

func (m *mockConfigWithValidator) Validate() error {
	return nil
}

type mockErrConfigWithValidator struct{}

func (m *mockErrConfigWithValidator) Validate() error {
	return fmt.Errorf("test error")
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: false,
		},
		{
			name:    "basic type",
			cfg:     "test",
			wantErr: false,
		},
		{
			name:    "struct without validator",
			cfg:     mockConfig{},
			wantErr: false,
		},
		{
			name:    "struct with validator",
			cfg:     &mockConfigWithValidator{},
			wantErr: false,
		},
		{
			name:    "slice of configs",
			cfg:     []Config{"test1", "test2"},
			wantErr: false,
		},
		{
			name:    "map of configs",
			cfg:     map[string]Config{"key1": "value1"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateConfig_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		errText string
	}{
		{
			name:    "validator returns error",
			cfg:     mockErrConfigWithValidator{},
			errText: "test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.cfg)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errText)
		})
	}
}

func TestCallValidateIfPossible(t *testing.T) {
	tests := []struct {
		name    string
		value   reflect.Value
		wantErr bool
	}{
		{
			name:    "non-validator type",
			value:   reflect.ValueOf("test"),
			wantErr: false,
		},
		{
			name:    "validator type",
			value:   reflect.ValueOf(&mockConfigWithValidator{}),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := callValidateIfPossible(tt.value)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
