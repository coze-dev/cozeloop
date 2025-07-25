// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package ctxcache

import (
	"context"
	"sync"
)

type ctxCacheKey struct{}

func Init(ctx context.Context) context.Context {
	if val, exist := ctx.Value(ctxCacheKey{}).(*sync.Map); exist && val != nil {
		return ctx
	}
	return context.WithValue(ctx, ctxCacheKey{}, new(sync.Map))
}

func Get[T any](ctx context.Context, key any) (value T, ok bool) {
	var zero T

	cacheMap, valid := ctx.Value(ctxCacheKey{}).(*sync.Map)
	if !valid {
		return zero, false
	}

	loadedValue, exists := cacheMap.Load(key)
	if !exists {
		return zero, false
	}

	if v, match := loadedValue.(T); match {
		return v, true
	}

	return zero, false
}

func Store(ctx context.Context, key any, obj any) {
	if cacheMap, ok := ctx.Value(ctxCacheKey{}).(*sync.Map); ok {
		cacheMap.Store(key, obj)
	}
}
