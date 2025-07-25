// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

func (s *DatasetServiceImpl) UpdateSchema(ctx context.Context, ds *entity.Dataset, fields []*entity.FieldSchema, updatedBy string) error {
	// 拼装 + 校验
	preSchema, err := s.repo.GetSchema(ctx, ds.SpaceID, ds.SchemaID)
	if err != nil {
		return err
	}
	compatible := schemaCompatible(preSchema.Fields, fields)
	if !compatible {
		if err := s.ensureEmptyDataset(ctx, ds); err != nil {
			return err
		}
	}
	allFields, err := mergeSchema(ds, preSchema, fields)
	if err != nil {
		return err
	}
	postCheck := func() error {
		if compatible {
			return nil
		}
		return s.ensureEmptyDataset(ctx, ds) // 仅空数据集允许不兼容变更 :)
	}
	// 加锁
	release, err := s.withUpdateSchemaBarrier(ctx, ds.ID)
	if err != nil {
		return err
	}
	defer release()

	// 原地更新
	if !preSchema.Immutable {
		preSchema.Fields = allFields
		preSchema.UpdatedBy = updatedBy
		return s.updateSchema(ctx, preSchema, postCheck)
	}

	// 使用新的 schema
	logs.CtxInfo(ctx, "rotate schema to new id, dataset_id=%d, pre_schema_id=%d", ds.ID, preSchema.ID)
	if err := s.rotateSchema(ctx, ds, allFields, updatedBy, postCheck); err != nil {
		return err
	}
	return nil
}

func (s *DatasetServiceImpl) ensureEmptyDataset(ctx context.Context, ds *entity.Dataset) error {
	itemCount, err := s.repo.GetItemCount(ctx, ds.ID)
	if err != nil {
		return err
	}
	if itemCount > 0 {
		return errno.Errorf(errno.ImcompatibleDatasetSchemaCode, `imcompatible schema change on non-empty dataset`)
	}
	return nil
}

func (s *DatasetServiceImpl) rotateSchema(ctx context.Context, ds *entity.Dataset, fields []*entity.FieldSchema, updatedBy string, postCheck func() error) error {
	schema := newSchemaOfDataset(ds, fields)
	schema.UpdatedBy = updatedBy
	schema.ID = 0

	return s.txDB.Transaction(ctx, func(tx *gorm.DB) error {
		opt := repo.WithTransaction(tx)

		if err := s.repo.CreateSchema(ctx, schema, opt); err != nil {
			return err
		}

		patch := &entity.Dataset{
			SchemaID:      schema.ID,
			LastOperation: entity.DatasetOpTypeUpdateSchema,
			UpdatedBy:     updatedBy,
		}

		if err := s.repo.PatchDataset(ctx, patch, repo.NewDatasetWhere(ds.SpaceID, ds.ID), opt); err != nil {
			return err
		}

		return postCheck()
	})
}

func (s *DatasetServiceImpl) updateSchema(ctx context.Context, schema *entity.DatasetSchema, postCheck func() error) error {
	return s.txDB.Transaction(ctx, func(tx *gorm.DB) error {
		opt := repo.WithTransaction(tx)

		v := schema.UpdateVersion
		schema.UpdateVersion = v + 1

		if err := s.repo.UpdateSchema(ctx, v, schema, opt); err != nil {
			return err
		}

		patch := &entity.Dataset{LastOperation: entity.DatasetOpTypeUpdateSchema, UpdatedBy: schema.UpdatedBy}
		if err := s.repo.PatchDataset(ctx, patch, repo.NewDatasetWhere(schema.SpaceID, schema.DatasetID), opt); err != nil {
			return err
		}

		return postCheck()
	})
}
