// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/coze-dev/coze-loop/backend/infra/redis"
	"github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/entity"
	"github.com/coze-dev/coze-loop/backend/modules/evaluation/infra/repo/experiment/redis/convert"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/conv"
)

//go:generate mockgen -source=quota.go -destination=mocks/quota.go -package=mocks -mock_names=IQuotaDAO=MockIQuotaDAO
type IQuotaDAO interface {
	GetQuotaSpaceExpt(ctx context.Context, spaceID int64) (*entity.QuotaSpaceExpt, error)
	SetQuotaSpaceExpt(ctx context.Context, spaceID int64, qse *entity.QuotaSpaceExpt) error
}

func NewQuotaDAO(cmdable redis.Cmdable) IQuotaDAO {
	const table = "experiment"
	return &quotaDAOImpl{cmdable: cmdable, table: table}
}

type quotaDAOImpl struct {
	cmdable redis.Cmdable
	table   string
}

func (q *quotaDAOImpl) GetQuotaSpaceExpt(ctx context.Context, spaceID int64) (*entity.QuotaSpaceExpt, error) {
	key := q.makeQuotaSpaceExptKey(spaceID)
	got, err := q.cmdable.Get(ctx, key).Result()
	if err != nil && !redis.IsNilError(err) {
		return nil, errorx.Wrapf(err, "redis get fail, key: %v", key)
	}
	return convert.NewQuotaSpaceExptConverter().ToDO(conv.UnsafeStringToBytes(got))
}

func (q *quotaDAOImpl) SetQuotaSpaceExpt(ctx context.Context, spaceID int64, qse *entity.QuotaSpaceExpt) error {
	bytes, err := convert.NewQuotaSpaceExptConverter().FromDO(qse)
	if err != nil {
		return err
	}

	key := q.makeQuotaSpaceExptKey(spaceID)
	if err := q.cmdable.Set(ctx, key, bytes, time.Hour*24*2).Err(); err != nil {
		return errorx.Wrapf(err, "redis set key: %v", key)
	}

	return nil
}

func (q *quotaDAOImpl) makeQuotaSpaceExptKey(spaceID int64) string {
	return fmt.Sprintf("[%s]quota_space_expt:%d", q.table, spaceID)
}
