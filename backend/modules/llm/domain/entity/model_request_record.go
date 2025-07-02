// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package entity

import "time"

type ModelRequestRecord struct {
	ID                  int64     `json:"id"`
	SpaceID             int64     `json:"space_id"`
	UserID              string    `json:"user_id"`
	UsageScene          Scenario  `json:"usage_scene"`
	UsageSceneEntityID  string    `json:"usage_scene_entity_id"`
	Frame               Frame     `json:"frame"`
	Protocol            Protocol  `json:"protocol"`
	ModelIdentification string    `json:"model_identification"`
	ModelAk             string    `json:"model_ak"`
	ModelID             string    `json:"model_id"`
	ModelName           string    `json:"model_name"`
	InputToken          int64     `json:"input_token"`
	OutputToken         int64     `json:"output_token"`
	Logid               string    `json:"logid"`
	ErrorCode           string    `json:"error_code"`
	ErrorMsg            *string   `json:"error_msg"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
