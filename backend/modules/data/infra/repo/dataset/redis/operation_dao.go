// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"time"

	"github.com/bytedance/gg/gslice"
	"github.com/bytedance/sonic"
	goredis "github.com/redis/go-redis/v9"

	"github.com/coze-dev/coze-loop/backend/infra/redis"
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/coze-loop/backend/modules/data/infra/rediskey"
	"github.com/coze-dev/coze-loop/backend/modules/data/pkg/errno"
	"github.com/coze-dev/coze-loop/backend/pkg/logs"
)

type OperationDAO interface {
	AddDatasetOperation(ctx context.Context, datasetID int64, op *entity.DatasetOperation) error
	DelDatasetOperation(ctx context.Context, datasetID int64, opType entity.DatasetOpType, id string) error
	MGetDatasetOperations(ctx context.Context, datasetID int64, opTypes []entity.DatasetOpType) (map[entity.DatasetOpType][]*entity.DatasetOperation, error)
}

func NewOperationDAO(redisCli redis.Cmdable) OperationDAO {
	return &OperationDAOImpl{
		redis: redisCli,
	}
}

type OperationDAOImpl struct {
	redis redis.Cmdable
}

func (r *OperationDAOImpl) AddDatasetOperation(ctx context.Context, datasetID int64, op *entity.DatasetOperation) error {
	data, err := sonic.Marshal(op)
	if err != nil {
		return errno.InternalErr(err, "marshal dataset operation")
	}

	key := rediskey.FormatDatasetOperationKey(datasetID, string(op.Type))
	cli := r.redis
	if err := cli.HSet(ctx, key, op.ID, data).Err(); err != nil {
		return errno.RedisErr(err, "key=%s", key)
	}

	return nil
}

func (r *OperationDAOImpl) DelDatasetOperation(ctx context.Context, datasetID int64, opType entity.DatasetOpType, id string) error {
	key := rediskey.FormatDatasetOperationKey(datasetID, string(opType))
	cli := r.redis
	if err := cli.HDel(ctx, key, id).Err(); err != nil {
		return errno.RedisErr(err, "key=%s", key)
	}
	return nil
}

func (r *OperationDAOImpl) MGetDatasetOperations(ctx context.Context, datasetID int64, opTypes []entity.DatasetOpType) (map[entity.DatasetOpType][]*entity.DatasetOperation, error) {
	if len(opTypes) == 0 {
		return nil, nil
	}

	keys := gslice.Map(opTypes, func(opType entity.DatasetOpType) string {
		return rediskey.FormatDatasetOperationKey(datasetID, string(opType))
	})

	cli := r.redis
	pipe := cli.Pipeline()

	gCmds := make([]*goredis.MapStringStringCmd, 0, len(keys))
	for _, key := range keys {
		gCmds = append(gCmds, pipe.HGetAll(ctx, key))
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return nil, errno.RedisErr(err, "keys=%v", keys)
	}

	var (
		m     = make(map[entity.DatasetOpType][]*entity.DatasetOperation)
		now   = time.Now()
		dCmds []*goredis.IntCmd
	)

	for i, cmd := range gCmds {
		if i >= len(opTypes) {
			break
		}

		opType := opTypes[i]
		for id, val := range cmd.Val() {
			op := &entity.DatasetOperation{}
			if err := sonic.UnmarshalString(val, op); err != nil {
				return nil, errno.InternalErr(err, "unmarshal dataset operation, key=%s[%s]", keys[i], id)
			}
			op.ID = id
			op.Type = opType
			if op.TS.Add(op.TTL).Before(now) {
				key := rediskey.FormatDatasetOperationKey(datasetID, string(opType))
				dCmds = append(dCmds, pipe.HDel(ctx, key, id))
				continue
			}
			m[opType] = append(m[opType], op)
		}
	}

	if len(dCmds) > 0 {
		cmds := gslice.Map(dCmds, (*goredis.IntCmd).String)
		logs.CtxInfo(ctx, "del %d expired dataset operations, commands=%v", len(dCmds), cmds)
		if _, err := pipe.Exec(ctx); err != nil {
			logs.CtxError(ctx, "del expired dataset operations failed, commands=%v, err=%v", cmds, err)
		}
	}
	return m, nil
}
