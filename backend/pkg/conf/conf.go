// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"context"
)

//go:generate mockgen -destination=mocks/conf.go -package=mocks . IConfigLoader
type IConfigLoader interface {
	Get(ctx context.Context, key string) any
	UnmarshalKey(ctx context.Context, key string, value any, opts ...DecodeOptionFn) error
	Unmarshal(ctx context.Context, value any, opts ...DecodeOptionFn) error
}

//go:generate mockgen -destination=mocks/conf_factory.go -package=mocks . IConfigLoaderFactory
type IConfigLoaderFactory interface {
	NewConfigLoader(dsn string) (IConfigLoader, error)
}

type DecodeOptionFn func(opt *DecodeOption)

func WithTagName(tagName string) DecodeOptionFn {
	return func(d *DecodeOption) {
		d.TagName = tagName
	}
}

type DecodeOption struct {
	TagName string
}
