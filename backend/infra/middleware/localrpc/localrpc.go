// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package localrpc

import (
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
)

type localCallFlagCtxKeyType struct{}

var localCallFlagCtxKey = localCallFlagCtxKeyType{}

func LocalCallFlagMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp any) (err error) {
		ctx = context.WithValue(ctx, localCallFlagCtxKey, "1")
		return next(ctx, req, resp)
	}
}

func IsLocalCall(ctx context.Context) bool {
	flag, ok := ctx.Value(localCallFlagCtxKey).(string)
	if !ok {
		return false
	}
	return flag == "1"
}
