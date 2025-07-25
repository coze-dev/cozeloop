// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import "time"

type DebugLog struct {
	ID           int64     `json:"id"`
	PromptID     int64     `json:"prompt_id"`
	SpaceID      int64     `json:"space_id"`
	PromptKey    string    `json:"prompt_key"`
	Version      string    `json:"version"`
	InputTokens  int64     `json:"input_tokens"`
	OutputTokens int64     `json:"output_tokens"`
	StartedAt    time.Time `json:"started_at"`
	EndedAt      time.Time `json:"ended_at"`
	CostMS       int64     `json:"cost_ms"`
	StatusCode   int32     `json:"status_code"`
	DebuggedBy   string    `json:"debugged_by"`
	DebugID      int64     `json:"debug_id"`
	DebugStep    int32     `json:"debug_step"`
}
