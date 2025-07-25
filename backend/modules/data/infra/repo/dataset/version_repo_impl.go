// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"context"

	"github.com/bytedance/gg/gslice"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/convertor"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/pagination"
)

func (d *DatasetRepo) CreateVersion(ctx context.Context, version *entity.DatasetVersion, opt ...repo.Option) error {
	MaybeGenID(ctx, d.idGen, version)
	versionPO, err := convertor.VersionDO2PO(version)
	if err != nil {
		return err
	}
	err = d.versionDao.CreateVersion(ctx, versionPO, Opt2DBOpt(opt...)...)
	if err != nil {
		return err
	}
	version.ID = versionPO.ID
	version.CreatedAt = versionPO.CreatedAt
	return nil
}

func (d *DatasetRepo) GetVersion(ctx context.Context, spaceID, versionID int64, opt ...repo.Option) (*entity.DatasetVersion, error) {
	versionPO, err := d.versionDao.GetVersion(ctx, spaceID, versionID, Opt2DBOpt(opt...)...)
	if err != nil {
		return nil, err
	}
	version, err := convertor.ConvertVersionPOToDO(versionPO)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (d *DatasetRepo) MGetVersions(ctx context.Context, spaceID int64, ids []int64, opt ...repo.Option) ([]*entity.DatasetVersion, error) {
	versionPOs, err := d.versionDao.MGetVersions(ctx, spaceID, ids, Opt2DBOpt(opt...)...)
	if err != nil {
		return nil, err
	}
	return gslice.TryMap(versionPOs, convertor.ConvertVersionPOToDO).Get()
}

func (d *DatasetRepo) GetItemCountOfVersion(ctx context.Context, versionID int64) (*int64, error) {
	return d.versionRedisDao.GetItemCountOfVersion(ctx, versionID)
}

func (d *DatasetRepo) SetItemCountOfVersion(ctx context.Context, datasetID int64, n int64) error {
	return d.versionRedisDao.SetItemCountOfVersion(ctx, datasetID, n)
}

func (d *DatasetRepo) ListVersions(ctx context.Context, params *repo.ListDatasetVersionsParams, opt ...repo.Option) ([]*entity.DatasetVersion, *pagination.PageResult, error) {
	daoParam := &mysql.ListDatasetVersionsParams{
		Paginator:   params.Paginator,
		SpaceID:     params.SpaceID,
		DatasetID:   params.DatasetID,
		IDs:         params.IDs,
		Versions:    params.Versions,
		VersionNums: params.VersionNums,
		VersionLike: params.VersionLike,
	}
	pos, p, err := d.versionDao.ListVersions(ctx, daoParam, Opt2DBOpt(opt...)...)
	if err != nil {
		return nil, nil, err
	}
	dos, err := gslice.TryMap(pos, convertor.ConvertVersionPOToDO).Get()
	if err != nil {
		return nil, nil, err
	}
	return dos, p, nil
}

func (d *DatasetRepo) CountVersions(ctx context.Context, params *repo.ListDatasetVersionsParams, opt ...repo.Option) (int64, error) {
	daoParam := &mysql.ListDatasetVersionsParams{
		Paginator:   params.Paginator,
		SpaceID:     params.SpaceID,
		DatasetID:   params.DatasetID,
		IDs:         params.IDs,
		Versions:    params.Versions,
		VersionNums: params.VersionNums,
		VersionLike: params.VersionLike,
	}
	return d.versionDao.CountVersions(ctx, daoParam, Opt2DBOpt(opt...)...)
}

func (d *DatasetRepo) PatchVersion(ctx context.Context, patch, where *entity.DatasetVersion, opt ...repo.Option) error {
	po, err := convertor.VersionDO2PO(patch)
	if err != nil {
		return err
	}
	wherePO, err := convertor.VersionDO2PO(where)
	if err != nil {
		return err
	}
	return d.versionDao.PatchVersion(ctx, po, wherePO, Opt2DBOpt(opt...)...)
}
