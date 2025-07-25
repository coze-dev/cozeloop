// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package viper

import (
	"context"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"

	"github.com/coze-dev/cozeloop/backend/pkg/conf"
)

type fileConfLoaderFactoryOpt struct {
	path string
}

type FileConfLoaderFactoryOpt func(*fileConfLoaderFactoryOpt)

func WithFactoryConfigPath(path string) FileConfLoaderFactoryOpt {
	return func(c *fileConfLoaderFactoryOpt) {
		c.path = path
	}
}

type FileConfLoaderOpt func(*fileConfLoaderOpt)

func WithConfigPath(path string) FileConfLoaderOpt {
	return func(c *fileConfLoaderOpt) {
		c.path = path
	}
}

func WithSearchPathDir(search bool) FileConfLoaderOpt {
	return func(c *fileConfLoaderOpt) {
		c.search = search
	}
}

type fileConfLoaderOpt struct {
	path   string
	search bool
}

type fileConfLoader struct {
	v *viper.Viper
}

func (v *fileConfLoader) Get(ctx context.Context, key string) any {
	if v == nil || v.v == nil {
		return nil
	}
	return v.v.Get(key)
}

func (v *fileConfLoader) Unmarshal(ctx context.Context, value any, opts ...conf.DecodeOptionFn) error {
	if v == nil || v.v == nil {
		return nil
	}
	return v.v.Unmarshal(value, toViperDecodeOption(opts)...)
}

func (v *fileConfLoader) UnmarshalKey(ctx context.Context, key string, value any, opts ...conf.DecodeOptionFn) error {
	if v == nil || v.v == nil {
		return nil
	}
	return v.v.UnmarshalKey(key, value, toViperDecodeOption(opts)...)
}

type fileConfLoaderFactory struct {
	configPath string
}

func (f *fileConfLoaderFactory) NewConfigLoader(dsn string) (conf.IConfigLoader, error) {
	return NewFileConfLoader(dsn, WithConfigPath(f.configPath), WithSearchPathDir(true))
}

func toViperDecodeOption(opts []conf.DecodeOptionFn) []viper.DecoderConfigOption {
	opt := &conf.DecodeOption{}
	for _, fn := range opts {
		fn(opt)
	}

	var vopts []viper.DecoderConfigOption
	if len(opt.TagName) > 0 {
		vopts = append(vopts, func(config *mapstructure.DecoderConfig) {
			config.TagName = opt.TagName
		})
	}
	return vopts
}
