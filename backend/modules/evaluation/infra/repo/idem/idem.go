// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package idem

import (
	"context"
	"time"

	"github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/component/idem"
	"github.com/coze-dev/coze-loop/backend/modules/evaluation/infra/repo/idem/redis"
)

func NewIdempotentService(idemDao redis.IIdemDAO) idem.IdempotentService {
	return &idemRepoImpl{
		idemDao: idemDao,
	}
}

type idemRepoImpl struct {
	idemDao redis.IIdemDAO
}

func (i *idemRepoImpl) Set(ctx context.Context, key string, duration time.Duration) error {
	return i.idemDao.Set(ctx, key, duration)
}

func (i *idemRepoImpl) SetNX(ctx context.Context, key string, duration time.Duration) (bool, error) {
	return i.idemDao.SetNX(ctx, key, duration)
}

func (i *idemRepoImpl) Exist(ctx context.Context, key string) (bool, error) {
	return i.idemDao.Exist(ctx, key)
}

func (i *idemRepoImpl) Del(ctx context.Context, key string) error {
	return i.idemDao.Del(ctx, key)
}
