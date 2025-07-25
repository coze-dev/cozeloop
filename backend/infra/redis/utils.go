// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func IsNilError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, redis.Nil)
}
