// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
)

//go:generate mockgen -destination=mocks/evaluation_set_schema.go -package=mocks . EvaluationSetSchemaService
type EvaluationSetSchemaService interface {
	UpdateEvaluationSetSchema(ctx context.Context, spaceID, evaluationSetID int64, fieldSchema []*entity.FieldSchema) (err error)
}
