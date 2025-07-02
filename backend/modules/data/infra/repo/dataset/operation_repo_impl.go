// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"context"
	"strconv"

	"github.com/rs/xid"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

func (d *DatasetRepo) AddDatasetOperation(ctx context.Context, datasetID int64, op *entity.DatasetOperation) error {
	// generate id
	id, err := d.idGen.GenID(ctx)
	if err != nil {
		logs.CtxWarn(ctx, "gen id for item op failed, dataset_id=%d, err=%v", datasetID, err)
	}
	op.ID = strconv.FormatInt(id, 10)
	if id == 0 {
		op.ID = xid.New().String()
	}
	return d.optDAO.AddDatasetOperation(ctx, datasetID, op)
}

func (d *DatasetRepo) DelDatasetOperation(ctx context.Context, datasetID int64, opType entity.DatasetOpType, id string) error {
	return d.optDAO.DelDatasetOperation(ctx, datasetID, opType, id)
}

func (d *DatasetRepo) MGetDatasetOperations(ctx context.Context, datasetID int64, opTypes []entity.DatasetOpType) (map[entity.DatasetOpType][]*entity.DatasetOperation, error) {
	return d.optDAO.MGetDatasetOperations(ctx, datasetID, opTypes)
}
