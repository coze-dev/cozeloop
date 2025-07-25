// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mcache

import (
	"time"
)

type IByteCache interface {
	Get(key []byte) ([]byte, error)
	Set(key []byte, value []byte, expiration time.Duration) error
}
