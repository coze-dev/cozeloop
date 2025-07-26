// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
)

type ItemDataPO struct {
	Data         []*entity.FieldData `json:"data,omitempty"`          // 数据内容
	RepeatedData []*entity.ItemData  `json:"repeated_data,omitempty"` // 多轮数据内容，与 Data 互斥
}
