// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package session

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSessionMD(t *testing.T) {
	for _, tt := range []struct {
		name     string
		fn       AuthProviderFn
		wantUser *User
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fn: func(ctx context.Context, req any) (*User, error) {
				return &User{ID: "123", Name: "John Doe", Email: "john.doe@example.com", AppID: 1}, nil
			},
			wantUser: &User{ID: "123", Name: "John Doe", Email: "john.doe@example.com", AppID: 1},
		},
		{
			name: "error",
			fn: func(ctx context.Context, req any) (*User, error) {
				return nil, errors.New("error")
			},
			wantErr: assert.Error,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			md := NewSessionMD(tt.fn, NewSessionService())
			md(func(ctx context.Context, req, resp any) error {
				user, ok := UserInCtx(ctx)
				assert.True(t, ok)
				assert.Equal(t, tt.wantUser, user)
				return nil
			})
		})
	}

	t.Run("context utils", func(t *testing.T) {
		ctx := context.TODO()
		user := &User{ID: "123", Name: "John Doe", Email: "john.doe@example.com", AppID: 1}
		ctx = WithCtxUser(ctx, user)
		assert.Equal(t, user.ID, UserIDInCtxOrEmpty(ctx))
		assert.Equal(t, user.AppID, AppIDInCtxOrEmpty(ctx))
	})
}
