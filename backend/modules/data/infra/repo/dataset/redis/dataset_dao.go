// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"reflect"
	"strconv"

	"github.com/bytedance/gg/gslice"

	"github.com/coze-dev/coze-loop/backend/infra/redis"
	"github.com/coze-dev/coze-loop/backend/modules/data/infra/rediskey"
	"github.com/coze-dev/coze-loop/backend/modules/data/pkg/errno"
)

//go:generate mockgen -destination=mocks/dataset_dao.go -package=mocks . DatasetDAO
type DatasetDAO interface {
	GetItemCount(ctx context.Context, datasetID int64) (int64, error)
	MGetItemCount(ctx context.Context, datasetIDs ...int64) (map[int64]int64, error)
	IncrItemCount(ctx context.Context, datasetID int64, n int64) (int64, error)
	SetItemCount(ctx context.Context, datasetID int64, n int64) error
}

func NewDatasetDAO(redisCli redis.Cmdable) DatasetDAO {
	return &DatasetDAOImpl{
		redis: redisCli,
	}
}

type DatasetDAOImpl struct {
	redis redis.Cmdable
}

func (r *DatasetDAOImpl) SetItemCount(ctx context.Context, datasetID int64, n int64) error {
	key := rediskey.FormatDatasetItemCountKey(datasetID)
	cli := r.redis
	if err := cli.Set(ctx, key, n, 0).Err(); err != nil {
		return errno.RedisErr(err, "key=%s", key)
	}
	return nil
}

func (r *DatasetDAOImpl) IncrItemCount(ctx context.Context, datasetID int64, n int64) (int64, error) {
	key := rediskey.FormatDatasetItemCountKey(datasetID)
	cli := r.redis
	total, err := cli.IncrBy(ctx, key, n).Result()
	if err != nil {
		return 0, errno.RedisErr(err, "key=%s", key)
	}
	return total, nil
}

func (r *DatasetDAOImpl) GetItemCount(ctx context.Context, datasetID int64) (int64, error) {
	key := rediskey.FormatDatasetItemCountKey(datasetID)
	cli := r.redis
	n, err := cli.Get(ctx, key).Int64()
	if err != nil {
		if redis.IsNilError(err) {
			return 0, nil
		}
		return 0, errno.RedisErr(err, "key=%s", key)
	}
	return n, nil
}

func (r *DatasetDAOImpl) MGetItemCount(ctx context.Context, datasetIDs ...int64) (map[int64]int64, error) {
	if len(datasetIDs) == 0 {
		return nil, nil
	}

	keys := gslice.Map(datasetIDs, rediskey.FormatDatasetItemCountKey)
	cli := r.redis
	values, err := cli.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, errno.RedisErr(err, "keys=%v", keys)
	}

	m := make(map[int64]int64, len(keys))
	for i, v := range values {
		if i >= len(datasetIDs) {
			break
		}
		id := datasetIDs[i]
		if v == nil {
			m[id] = 0
			continue
		}
		vstr, ok := v.(string)
		if !ok {
			return nil, errno.InternalErr(err, "not support value type=%s, value=%v", reflect.TypeOf(v).String(), v)
		}
		n, err := strconv.ParseInt(vstr, 10, 64)
		if err != nil {
			return nil, errno.InternalErr(err, "convert %s to int64", v)
		}
		m[id] = n
	}
	return m, nil
}
