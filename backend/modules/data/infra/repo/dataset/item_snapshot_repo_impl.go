// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"context"

	"github.com/bytedance/gg/gslice"

	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/coze-loop/backend/modules/data/infra/repo/dataset/mysql"
	"github.com/coze-dev/coze-loop/backend/modules/data/infra/repo/dataset/mysql/convertor"
	"github.com/coze-dev/coze-loop/backend/modules/data/pkg/pagination"
)

func (d *DatasetRepo) BatchUpsertItemSnapshots(ctx context.Context, snapshots []*entity.ItemSnapshot, opt ...repo.Option) (int64, error) {
	MaybeGenID(ctx, d.idGen, snapshots...)
	pos, err := gslice.TryMap(snapshots, convertor.ItemSnapshotDO2PO).Get()
	if err != nil {
		return 0, err
	}
	return d.itemSnapshotDAO.BatchUpsertItemSnapshots(ctx, pos, Opt2DBOpt(opt...)...)
}

func (d *DatasetRepo) ListItemSnapshots(ctx context.Context, params *repo.ListItemSnapshotsParams, opt ...repo.Option) ([]*entity.ItemSnapshot, *pagination.PageResult, error) {
	daoParam := &mysql.ListItemSnapshotsParams{
		Paginator: params.Paginator,
		SpaceID:   params.SpaceID,
		VersionID: params.VersionID,
	}
	pos, p, err := d.itemSnapshotDAO.ListItemSnapshots(ctx, daoParam, Opt2DBOpt(opt...)...)
	if err != nil {
		return nil, nil, err
	}
	dos, err := gslice.TryMap(pos, convertor.ConvertItemSnapshotPO2DO).Get()
	if err != nil {
		return nil, nil, err
	}
	return dos, p, nil
}

func (d *DatasetRepo) CountItemSnapshots(ctx context.Context, params *repo.ListItemSnapshotsParams, opt ...repo.Option) (int64, error) {
	daoParam := &mysql.ListItemSnapshotsParams{
		Paginator: params.Paginator,
		SpaceID:   params.SpaceID,
		VersionID: params.VersionID,
	}
	return d.itemSnapshotDAO.CountItemSnapshots(ctx, daoParam, Opt2DBOpt(opt...)...)
}
