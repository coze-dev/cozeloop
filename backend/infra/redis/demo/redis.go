// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package demo

import (
	"github.com/redis/go-redis/v9"
)

type ClientFactory struct{}

func (c *ClientFactory) NewClient(opt *redis.Options) (*redis.Client, error) {
	return redis.NewClient(opt), nil
}
