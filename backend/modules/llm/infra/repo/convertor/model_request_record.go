// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"github.com/coze-dev/cozeloop/backend/modules/llm/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/llm/infra/repo/gorm_gen/model"
)

func ModelReqRecordDO2PO(record *entity.ModelRequestRecord) *model.ModelRequestRecord {
	if record == nil {
		return nil
	}
	return &model.ModelRequestRecord{
		ID:                  record.ID,
		SpaceID:             record.SpaceID,
		UserID:              record.UserID,
		UsageScene:          string(record.UsageScene),
		UsageSceneEntityID:  record.UsageSceneEntityID,
		Frame:               string(record.Frame),
		Protocol:            string(record.Protocol),
		ModelIdentification: record.ModelIdentification,
		ModelAk:             record.ModelAk,
		ModelID:             record.ModelID,
		ModelName:           record.ModelName,
		InputToken:          record.InputToken,
		OutputToken:         record.OutputToken,
		Logid:               record.Logid,
		ErrorCode:           record.ErrorCode,
		ErrorMsg:            record.ErrorMsg,
		CreatedAt:           record.CreatedAt,
		UpdatedAt:           record.UpdatedAt,
	}
}
