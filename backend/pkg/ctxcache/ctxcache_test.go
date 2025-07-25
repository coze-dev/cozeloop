// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package ctxcache

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCtxCache(t *testing.T) {
	type testCase struct {
		name        string
		setup       func(context.Context) context.Context
		key         interface{}
		value       interface{}
		expectFound bool
		expectValue interface{}
	}

	tests := []testCase{
		{
			name: "should not find value when context is not initialized",
			setup: func(ctx context.Context) context.Context {
				return ctx
			},
			key:         "test1",
			value:       "1",
			expectFound: false,
		},
		{
			name: "should not find non-existent key in initialized context",
			setup: func(ctx context.Context) context.Context {
				return Init(ctx)
			},
			key:         "test",
			expectFound: false,
		},
		{
			name: "should store and retrieve string value",
			setup: func(ctx context.Context) context.Context {
				return Init(ctx)
			},
			key:         "key1",
			value:       "abc",
			expectFound: true,
			expectValue: "abc",
		},
		{
			name: "should store and retrieve int64 value with struct key",
			setup: func(ctx context.Context) context.Context {
				return Init(ctx)
			},
			key:         struct{}{},
			value:       int64(1),
			expectFound: true,
			expectValue: int64(1),
		},
		{
			name: "should store and retrieve complex struct value",
			setup: func(ctx context.Context) context.Context {
				return Init(ctx)
			},
			key: "temp",
			value: struct {
				a string
				b string
				c int64
				d []int64
			}{
				a: "1",
				b: "2",
				c: 3,
				d: []int64{123, 1232, 232},
			},
			expectFound: true,
			expectValue: struct {
				a string
				b string
				c int64
				d []int64
			}{
				a: "1",
				b: "2",
				c: 3,
				d: []int64{123, 1232, 232},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = tt.setup(ctx)

			if tt.value != nil {
				Store(ctx, tt.key, tt.value)
			}

			got, ok := Get[any](ctx, tt.key)
			assert.Equal(t, tt.expectFound, ok)
			if tt.expectFound {
				assert.True(t, reflect.DeepEqual(got, tt.expectValue))
			}
		})
	}
}
