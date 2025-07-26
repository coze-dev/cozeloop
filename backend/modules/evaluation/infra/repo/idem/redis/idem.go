// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/coze-dev/coze-loop/backend/infra/redis"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
)

const (
	idemKeyValue = "1"
)

func NewIdemDAO(cmdable redis.Cmdable) IIdemDAO {
	const table = "idempotent"
	return &idemDAOImpl{cmdable: cmdable, table: table}
}

type idemDAOImpl struct {
	cmdable redis.Cmdable
	table   string
}

type IIdemDAO interface {
	Set(ctx context.Context, key string, duration time.Duration) error
	SetNX(ctx context.Context, key string, duration time.Duration) (bool, error)
	Exist(ctx context.Context, key string) (bool, error)
	Del(ctx context.Context, key string) error
}

func (i *idemDAOImpl) Set(ctx context.Context, key string, duration time.Duration) error {
	ckey := i.makeIdemKey(key)
	err := i.cmdable.Set(ctx, ckey, idemKeyValue, duration).Err()
	if err != nil {
		return errorx.Wrapf(err, "Set fail, key: %v", ckey)
	}
	return nil
}

func (i *idemDAOImpl) SetNX(ctx context.Context, key string, duration time.Duration) (bool, error) {
	ckey := i.makeIdemKey(key)
	ok, err := i.cmdable.SetNX(ctx, ckey, idemKeyValue, duration).Result()
	if err != nil {
		return false, errorx.Wrapf(err, "Set fail, key: %v", ckey)
	}
	return ok, nil
}

func (i *idemDAOImpl) Exist(ctx context.Context, key string) (bool, error) {
	ckey := i.makeIdemKey(key)
	res, err := i.cmdable.Exists(ctx, ckey).Result()
	if err != nil {
		return false, errorx.Wrapf(err, "Exist fail, key: %s", ckey)
	}
	return res == 1, nil
}

func (i *idemDAOImpl) Del(ctx context.Context, key string) error {
	ckey := i.makeIdemKey(key)
	_, err := i.cmdable.Del(ctx, ckey).Result()
	if err != nil {
		return errorx.Wrapf(err, "Del fail, key: %s", ckey)
	}
	return nil
}

func (i *idemDAOImpl) makeIdemKey(key string) string {
	return fmt.Sprintf("[%v]idem_key:%v", i.table, key)
}
