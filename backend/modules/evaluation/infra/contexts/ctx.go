// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package contexts

import (
	"context"
)

type ctxWriteDBKey struct{}

type ctxWriteDBVal struct{}

func WithCtxWriteDB(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxWriteDBKey{}, ctxWriteDBVal{})
}

func CtxWriteDB(ctx context.Context) bool {
	return ctx.Value(ctxWriteDBKey{}) != nil
}

type userIDKeyType struct{}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKeyType{}, userID)
}

func GetUserID(ctx context.Context) string {
	userID, ok := ctx.Value(userIDKeyType{}).(string)
	if !ok {
		return ""
	}
	return userID
}
