// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package byted

import (
	"time"

	"github.com/coocood/freecache"

	"github.com/coze-dev/cozeloop/backend/pkg/mcache"
)

// NewLRUCache size is in bytes.
// WARN: k/v must not exceed cache_size/1024. For example, in a 1GB cache,
// a single kv pair cannot exceed 1MB.
func NewLRUCache(size int) mcache.IByteCache {
	return &lruCache{
		c: freecache.NewCache(size),
	}
}

type lruCache struct {
	c *freecache.Cache
}

func (l *lruCache) Get(key []byte) ([]byte, error) {
	return l.c.Get(key)
}

func (l *lruCache) Set(key []byte, value []byte, expiration time.Duration) error {
	return l.c.Set(key, value, int(expiration.Seconds()))
}
