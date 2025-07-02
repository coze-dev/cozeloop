// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"fmt"

	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gslice"
	"gorm.io/gorm"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/component/mq"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

func (s *DatasetServiceImpl) BatchGetVersionedDatasetsWithOpt(ctx context.Context, spaceID int64, versionIDs []int64, opt *GetOpt) ([]*VersionedDatasetWithSchema, error) {
	if len(versionIDs) == 0 {
		return nil, nil
	}
	versions, err := s.repo.MGetVersions(ctx, spaceID, versionIDs)
	if err != nil {
		return nil, err
	}
	datasetIDs := gslice.Map(versions, func(v *entity.DatasetVersion) int64 { return v.DatasetID })
	var opts []repo.Option
	if opt.WithDeleted {
		opts = append(opts, repo.WithDeleted())
	}
	datasets, err := s.repo.MGetDatasets(ctx, spaceID, datasetIDs, opts...)
	if err != nil {
		return nil, err
	}
	schemaIDs := gslice.Map(versions, func(v *entity.DatasetVersion) int64 { return v.SchemaID })
	schemas, err := s.repo.MGetSchema(ctx, spaceID, schemaIDs)
	if err != nil {
		return nil, err
	}
	datasetM := gslice.ToMap(datasets, func(d *entity.Dataset) (int64, *entity.Dataset) {
		return d.ID, d
	})
	schemaM := gslice.ToMap(schemas, func(s *entity.DatasetSchema) (int64, *entity.DatasetSchema) {
		return s.ID, s
	})
	return gslice.FilterMap(versions, func(v *entity.DatasetVersion) (*VersionedDatasetWithSchema, bool) {
		dataset, ok1 := datasetM[v.DatasetID]
		schema, ok2 := schemaM[v.SchemaID]
		return &VersionedDatasetWithSchema{Version: v, Dataset: dataset, Schema: schema}, ok1 && ok2
	}), nil
}

func (s *DatasetServiceImpl) GetOrSetItemCountOfVersion(ctx context.Context, version *entity.DatasetVersion) (int64, error) {
	if version.SnapshotStatus == entity.SnapshotStatusCompleted {
		return version.ItemCount, nil
	}
	c, err := s.repo.GetItemCountOfVersion(ctx, version.ID)
	if err != nil {
		return 0, err
	}
	if c != nil {
		return gptr.Indirect(c), nil
	}

	logs.CtxInfo(ctx, "counting item count of version, dataset_id=%d, version_id=%d, version=%s", version.DatasetID, version.ID, version.Version)
	query := NewListItemsParamsFromVersion(version)
	n, err := s.repo.CountItems(ctx, query, repo.WithMaster())
	if err != nil {
		return 0, err
	}
	logs.CtxInfo(ctx, "%d item count found, version=%d", n, version.ID)
	if err := s.repo.SetItemCountOfVersion(ctx, version.ID, n); err != nil {
		return 0, err
	}
	return n, nil
}

func (s *DatasetServiceImpl) GetVersionWithOpt(ctx context.Context, spaceID, versionID int64, opt *GetOpt) (*entity.DatasetVersion, *DatasetWithSchema, error) {
	var opts []repo.Option
	if opt.WithDeleted {
		opts = append(opts, repo.WithDeleted())
	}
	version, err := s.repo.GetVersion(ctx, spaceID, versionID)
	if err != nil {
		return nil, nil, err
	}
	if version == nil {
		return nil, nil, errno.NotFoundErrorf("version %d is not found", versionID)
	}
	if version.SnapshotStatus != entity.SnapshotStatusCompleted {
		count, err := s.GetOrSetItemCountOfVersion(ctx, version)
		if err != nil {
			return nil, nil, err
		}
		version.ItemCount = count
	}
	dataset, err := s.repo.GetDataset(ctx, spaceID, version.DatasetID, opts...)
	if err != nil {
		return nil, nil, err
	}
	if dataset == nil {
		return nil, nil, errno.InternalErrorf("dataset %d is not found, space_id=%d", version.DatasetID, spaceID)
	}
	schema, err := s.repo.GetSchema(ctx, spaceID, version.SchemaID)
	if err != nil {
		return nil, nil, err
	}
	return version, &DatasetWithSchema{Dataset: dataset, Schema: schema}, nil
}

func (s *DatasetServiceImpl) CreateVersion(ctx context.Context, ds *DatasetWithSchema, version *entity.DatasetVersion) error {
	if err := validateVersion(ds.LatestVersion, version.Version); err != nil {
		return err
	}
	patchVersionWithDataset(ds.Dataset, version)

	release, err := s.withCreateVersionBarrier(ctx, ds.ID)
	if err != nil {
		return err
	}
	defer release()

	if err := s.createVersion(ctx, ds, version); err != nil {
		return err
	}
	msg := &entity.JobRunMessage{
		Type:     entity.DatasetSnapshotJob,
		SpaceID:  ds.SpaceID,
		Extra:    map[string]string{"version_id": fmt.Sprintf("%d", version.ID)},
		Operator: ds.UpdatedBy,
	}
	err = s.producer.Send(ctx, msg, mq.WithKey(fmt.Sprintf("%d", version.ID)))
	if err != nil {
		return errno.InternalErr(err, "send create_snapshot msg failed, version_id:%d", version.ID)
	}
	return nil
}

func (s *DatasetServiceImpl) createVersion(ctx context.Context, ds *DatasetWithSchema, version *entity.DatasetVersion) error {
	if err := s.txDB.Transaction(ctx, func(tx *gorm.DB) error {
		opt := repo.WithTransaction(tx)
		// 插入 version
		if err := s.repo.CreateVersion(ctx, version, opt); err != nil {
			return err
		}

		// 更新 dataset 关联的版本信息
		dWhere := repo.NewDatasetWhere(ds.SpaceID, ds.ID, func(d *entity.Dataset) {
			d.NextVersionNum = ds.NextVersionNum // 防止版本并发
		})
		dPatch := &entity.Dataset{
			LatestVersion:  version.Version,
			NextVersionNum: ds.NextVersionNum + 1,
			LastOperation:  entity.DatasetOpTypeCreateVersion,
			UpdatedBy:      ds.UpdatedBy,
		}
		if err := s.repo.PatchDataset(ctx, dPatch, dWhere, opt); err != nil {
			return err
		}

		// 锁定 schema 不可变更
		if ds.Schema.Immutable {
			return nil
		}
		sPatch := &entity.DatasetSchema{
			ID:            ds.SchemaID,
			SpaceID:       ds.SpaceID,
			UpdateVersion: ds.Schema.UpdateVersion + 1,
			Immutable:     true,
		}
		return s.repo.UpdateSchema(ctx, ds.Schema.UpdateVersion, sPatch, opt)
	}); err != nil {
		return err
	}
	return nil
}

func NewListItemsParamsFromVersion(version *entity.DatasetVersion, taps ...func(*repo.ListItemsParams)) *repo.ListItemsParams {
	params := &repo.ListItemsParams{
		SpaceID:   version.SpaceID,
		DatasetID: version.DatasetID,
		DelVNGt:   version.VersionNum,
		AddVNLte:  version.VersionNum,
	}
	for _, tap := range taps {
		tap(params)
	}
	return params
}
