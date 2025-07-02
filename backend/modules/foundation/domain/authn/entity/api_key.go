// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"time"
)

const (
	APIKeyStatusNormal  = 0
	APIKeyStatusDeleted = 1
)

type APIKey struct {
	ID         int64
	Key        string
	Name       string
	Status     int32
	UserID     int64
	ExpiredAt  int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  int64
	LastUsedAt int64
}
