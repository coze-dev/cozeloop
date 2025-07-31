// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"sync"

	"github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/component/rpc"
	"github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/entity"
)

var (
	evaluationSetSchemaServiceOnce = sync.Once{}
	evaluationSetSchemaServiceImpl EvaluationSetSchemaService
)

type EvaluationSetSchemaServiceImpl struct {
	datasetRPCAdapter rpc.IDatasetRPCAdapter
}

func NewEvaluationSetSchemaServiceImpl(datasetRPCAdapter rpc.IDatasetRPCAdapter) EvaluationSetSchemaService {
	evaluationSetSchemaServiceOnce.Do(func() {
		evaluationSetSchemaServiceImpl = &EvaluationSetSchemaServiceImpl{
			datasetRPCAdapter: datasetRPCAdapter,
		}
	})
	return evaluationSetSchemaServiceImpl
}

func (d *EvaluationSetSchemaServiceImpl) UpdateEvaluationSetSchema(ctx context.Context, spaceID, evaluationSetID int64, fieldSchema []*entity.FieldSchema) (err error) {
	// 依赖数据集服务
	return d.datasetRPCAdapter.UpdateDatasetSchema(ctx, spaceID, evaluationSetID, fieldSchema)
}
