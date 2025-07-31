// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"context"

	"github.com/bytedance/gg/gslice"

	"github.com/coze-dev/coze-loop/backend/infra/db"
	"github.com/coze-dev/coze-loop/backend/infra/idgen"
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/coze-loop/backend/pkg/logs"
)

type getSetID interface {
	GetID() int64
	canSetID
}

type canSetID interface {
	SetID(int64)
}

// MaybeGenID 在 ID 不为时，生成 ID
func MaybeGenID[T getSetID](ctx context.Context, cli idgen.IIDGenerator, s ...T) {
	noID := gslice.Filter(s, func(t T) bool { return t.GetID() <= 0 })
	n := len(noID)
	if n == 0 {
		return
	}

	ids, err := cli.GenMultiIDs(ctx, n)
	if err != nil {
		logs.CtxWarn(ctx, "generate ids failed, id will be auto-incremented, err=%v", err)
		return
	}
	if len(ids) != n {
		logs.CtxWarn(ctx, "generate ids got %d ids, want %d", len(ids), n)
		return
	}
	for i, ele := range noID {
		ele.SetID(ids[i])
	}
}

func Opt2DBOpt(opt ...repo.Option) []db.Option {
	repoOpt := &repo.Opt{}
	for _, fn := range opt {
		fn(repoOpt)
	}
	res := []db.Option{}
	if repoOpt.WithDeleted {
		res = append(res, db.WithDeleted())
	}
	if repoOpt.WithMaster {
		res = append(res, db.WithMaster())
	}
	if repoOpt.TX != nil {
		res = append(res, db.WithTransaction(repoOpt.TX))
	}
	return res
}
