// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package session

import (
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/pkg/errors"

	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
)

// NewSessionMD creates a session middleware for Kitex.
func NewSessionMD(ap AuthProvider, ss ISessionService) func(next endpoint.Endpoint) endpoint.Endpoint {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp any) error {
			// 获取会话信息
			sessionID, ok := ctx.Value(SessionKey).(string)
			if !ok || sessionID == "" {
				return errorx.New("session not found")
			}

			// 验证会话
			session, err := ss.ValidateSession(ctx, sessionID)
			if err != nil {
				return errors.Wrap(err, "validate session failed")
			}

			// 获取用户信息
			user, err := ap.GetLoginUser(ctx, req)
			if err != nil {
				return errors.Wrap(err, "get login user failed")
			}

			// 验证用户ID是否匹配
			if user.ID != session.UserID {
				return errorx.New("user id mismatch")
			}

			// 注入用户信息到上下文
			ctx = WithCtxUser(ctx, user)

			return next(ctx, req, resp)
		}
	}
}

type AuthProvider interface {
	GetLoginUser(ctx context.Context, req any) (*User, error)
}

// AuthProviderFn 是 AuthProvider 的函数类型实现
type AuthProviderFn func(ctx context.Context, req any) (*User, error)

func (fn AuthProviderFn) GetLoginUser(ctx context.Context, req any) (*User, error) {
	return fn(ctx, req)
}
