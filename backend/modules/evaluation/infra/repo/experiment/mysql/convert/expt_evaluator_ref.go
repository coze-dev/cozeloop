// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convert

import (
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/infra/repo/experiment/mysql/gorm_gen/model"
)

func NewExptEvaluatorRefConverter() *ExptEvaluatorRefConverter {
	return &ExptEvaluatorRefConverter{}
}

type ExptEvaluatorRefConverter struct{}

func (ExptEvaluatorRefConverter) DO2PO(refs []*entity.ExptEvaluatorRef) []*model.ExptEvaluatorRef {
	models := make([]*model.ExptEvaluatorRef, 0, len(refs))
	for _, ref := range refs {
		models = append(models, &model.ExptEvaluatorRef{
			ID:                 ref.ID,
			SpaceID:            ref.SpaceID,
			ExptID:             ref.ExptID,
			EvaluatorID:        ref.EvaluatorID,
			EvaluatorVersionID: ref.EvaluatorVersionID,
		})
	}
	return models
}

func (ExptEvaluatorRefConverter) PO2DO(refs []*model.ExptEvaluatorRef) []*entity.ExptEvaluatorRef {
	entities := make([]*entity.ExptEvaluatorRef, 0, len(refs))
	for _, ref := range refs {
		entities = append(entities, &entity.ExptEvaluatorRef{
			SpaceID:            ref.SpaceID,
			ExptID:             ref.ExptID,
			EvaluatorID:        ref.EvaluatorID,
			EvaluatorVersionID: ref.EvaluatorVersionID,
		})
	}
	return entities
}
