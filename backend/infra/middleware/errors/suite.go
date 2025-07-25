// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package errors

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/server"
)

type option struct {
	kiteXSvrMWCompat bool
}

type Option func(ws *option)

func WithKiteXSvrMWCompat() Option {
	return func(o *option) {
		o.kiteXSvrMWCompat = true
	}
}

func NewServerSuite(opts ...Option) server.Suite {
	opt := &option{}
	for _, o := range opts {
		o(opt)
	}
	return &serverSuite{opt.kiteXSvrMWCompat}
}

type serverSuite struct {
	kiteXSvrMWCompat bool
}

func (suite *serverSuite) Options() []server.Option {
	opts := []server.Option{server.WithMiddleware(KiteXSvrErrorWrapMW())}
	if suite.kiteXSvrMWCompat {
		opts = append(opts, server.WithMiddleware(KiteXSvrCompatMW()))
	}
	return opts
}

func NewClientSuite() client.Suite {
	return &clientSuite{}
}

type clientSuite struct{}

func (suite *clientSuite) Options() []client.Option {
	return []client.Option{
		client.WithMiddleware(KiteXSvrErrorWrapMW()),
	}
}
