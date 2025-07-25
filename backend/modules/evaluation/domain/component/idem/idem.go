// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package idem

import (
	"context"
	"time"
)

//go:generate mockgen -destination=mocks/idempotent_service.go -package=mocks . IdempotentService
type IdempotentService interface {
	Set(ctx context.Context, key string, duration time.Duration) error
	SetNX(ctx context.Context, key string, duration time.Duration) (bool, error)
	Exist(ctx context.Context, key string) (bool, error)
	Del(ctx context.Context, key string) error
}
