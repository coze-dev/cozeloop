// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"

	"github.com/bytedance/gg/gptr"

	"github.com/coze-dev/coze-loop/backend/infra/redis"
	"github.com/coze-dev/coze-loop/backend/modules/data/infra/rediskey"
	"github.com/coze-dev/coze-loop/backend/modules/data/pkg/errno"
)

type VersionDAO interface {
	GetItemCountOfVersion(ctx context.Context, versionID int64) (*int64, error)
	SetItemCountOfVersion(ctx context.Context, versionID int64, n int64) error
}

func NewVersionDAO(redisCli redis.Cmdable) VersionDAO {
	return &VersionDAOImpl{
		redis: redisCli,
	}
}

type VersionDAOImpl struct {
	redis redis.Cmdable
}

func (r *VersionDAOImpl) GetItemCountOfVersion(ctx context.Context, versionID int64) (*int64, error) {
	key := rediskey.FormatDatasetVersionItemCountKey(versionID)
	cli := r.redis
	count, err := cli.Get(ctx, key).Int64()
	if err != nil {
		if redis.IsNilError(err) {
			return nil, nil
		}
		return nil, errno.RedisErr(err, "key=%s", key)
	}
	return gptr.Of(count), nil
}

func (r *VersionDAOImpl) SetItemCountOfVersion(ctx context.Context, versionID int64, n int64) error {
	key := rediskey.FormatDatasetVersionItemCountKey(versionID)
	cli := r.redis
	if err := cli.Set(ctx, key, n, 0).Err(); err != nil {
		return errno.RedisErr(err, "key=%s", key)
	}
	return nil
}
