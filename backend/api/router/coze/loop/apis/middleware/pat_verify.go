// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"
	"strings"

	"github.com/bytedance/gg/gptr"
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-loop/backend/api/handler/coze/loop/apis"
	"github.com/coze-dev/coze-loop/backend/infra/middleware/session"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/foundation/authn"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/foundation/user"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
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
		if verifyRes.Valid == nil || !*verifyRes.Valid || len(verifyRes.GetUserID()) == 0 {
			_ = c.Error(errorx.New("invalid pat token"))
			c.Abort()
			return
		}

		userID := verifyRes.GetUserID()
		resp, err := handler.GetUserInfo(ctx, &user.GetUserInfoRequest{
			UserID: gptr.Of(userID),
		})
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		if resp.GetUserInfo() == nil {
			_ = c.Error(errorx.New("user not found"))
			c.Abort()
			return
		}

		ctx = session.WithCtxUser(ctx, &session.User{
			ID:    userID,
			Name:  resp.GetUserInfo().GetName(),
			Email: resp.GetUserInfo().GetEmail(),
		})

		c.Next(ctx)
	}
}
