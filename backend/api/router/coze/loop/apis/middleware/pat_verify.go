// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/cozeloop/backend/api/handler/coze/loop/apis"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/authn"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
)

func PatTokenVerifyMW(handler *apis.APIHandler) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authHeader := c.GetHeader("Authorization")

		if len(authHeader) == 0 {
			_ = c.Error(errorx.New("authorization header is empty"))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(string(authHeader), "Bearer ")
		verifyRes, err := handler.VerifyToken(ctx, &authn.VerifyTokenRequest{
			Token: token,
		})
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		if verifyRes.Valid == nil || !*verifyRes.Valid {
			_ = c.Error(errorx.New("invalid pat token"))
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}
