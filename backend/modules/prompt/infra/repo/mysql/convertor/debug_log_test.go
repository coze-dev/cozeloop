// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/infra/repo/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
)

func TestDebugLogsPO2DO(t *testing.T) {
	tests := []struct {
		name     string
		pos      []*model.PromptDebugLog
		expected []*entity.DebugLog
	}{
		{
			name:     "nil input",
			pos:      nil,
			expected: nil,
		},
		{
			name:     "empty slice",
			pos:      []*model.PromptDebugLog{},
			expected: []*entity.DebugLog{},
		},
		{
			name: "slice with nil element",
			pos: []*model.PromptDebugLog{
				nil,
				{
					ID:           1,
					PromptID:     100,
					SpaceID:      200,
					PromptKey:    "test_key",
					Version:      "1.0.0",
					InputTokens:  10,
					OutputTokens: 20,
					StartedAt:    ptr.Of(int64(1000)),
					EndedAt:      ptr.Of(int64(2000)),
					CostMs:       ptr.Of(int64(1000)),
					StatusCode:   ptr.Of(int32(200)),
					DebuggedBy:   ptr.Of("test_user"),
					DebugID:      1,
					DebugStep:    1,
				},
			},
			expected: []*entity.DebugLog{
				{
					ID:           1,
					PromptID:     100,
					SpaceID:      200,
					PromptKey:    "test_key",
					Version:      "1.0.0",
					InputTokens:  10,
					OutputTokens: 20,
					StartedAt:    time.UnixMilli(1000),
					EndedAt:      time.UnixMilli(2000),
					CostMS:       1000,
					StatusCode:   200,
					DebuggedBy:   "test_user",
					DebugID:      1,
					DebugStep:    1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DebugLogsPO2DO(tt.pos)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestDebugLogPO2DO(t *testing.T) {
	tests := []struct {
		name     string
		po       *model.PromptDebugLog
		expected *entity.DebugLog
	}{
		{
			name:     "nil input",
			po:       nil,
			expected: nil,
		},
		{
			name: "complete debug log",
			po: &model.PromptDebugLog{
				ID:           1,
				PromptID:     100,
				SpaceID:      200,
				PromptKey:    "test_key",
				Version:      "1.0.0",
				InputTokens:  10,
				OutputTokens: 20,
				StartedAt:    ptr.Of(int64(1000)),
				EndedAt:      ptr.Of(int64(2000)),
				CostMs:       ptr.Of(int64(1000)),
				StatusCode:   ptr.Of(int32(200)),
				DebuggedBy:   ptr.Of("test_user"),
				DebugID:      1,
				DebugStep:    1,
			},
			expected: &entity.DebugLog{
				ID:           1,
				PromptID:     100,
				SpaceID:      200,
				PromptKey:    "test_key",
				Version:      "1.0.0",
				InputTokens:  10,
				OutputTokens: 20,
				StartedAt:    time.UnixMilli(1000),
				EndedAt:      time.UnixMilli(2000),
				CostMS:       1000,
				StatusCode:   200,
				DebuggedBy:   "test_user",
				DebugID:      1,
				DebugStep:    1,
			},
		},
		{
			name: "debug log with nil optional fields",
			po: &model.PromptDebugLog{
				ID:           1,
				PromptID:     100,
				SpaceID:      200,
				PromptKey:    "test_key",
				Version:      "1.0.0",
				InputTokens:  10,
				OutputTokens: 20,
				StartedAt:    nil,
				EndedAt:      nil,
				CostMs:       nil,
				StatusCode:   nil,
				DebuggedBy:   nil,
				DebugID:      1,
				DebugStep:    1,
			},
			expected: &entity.DebugLog{
				ID:           1,
				PromptID:     100,
				SpaceID:      200,
				PromptKey:    "test_key",
				Version:      "1.0.0",
				InputTokens:  10,
				OutputTokens: 20,
				StartedAt:    time.UnixMilli(0),
				EndedAt:      time.UnixMilli(0),
				CostMS:       0,
				StatusCode:   0,
				DebuggedBy:   "",
				DebugID:      1,
				DebugStep:    1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DebugLogPO2DO(tt.po)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestDebugLogDO2PO(t *testing.T) {
	tests := []struct {
		name     string
		do       *entity.DebugLog
		expected *model.PromptDebugLog
	}{
		{
			name:     "nil input",
			do:       nil,
			expected: nil,
		},
		{
			name: "complete debug log",
			do: &entity.DebugLog{
				ID:           1,
				PromptID:     100,
				SpaceID:      200,
				PromptKey:    "test_key",
				Version:      "1.0.0",
				InputTokens:  10,
				OutputTokens: 20,
				StartedAt:    time.UnixMilli(1000),
				EndedAt:      time.UnixMilli(2000),
				CostMS:       1000,
				StatusCode:   200,
				DebuggedBy:   "test_user",
				DebugID:      1,
				DebugStep:    1,
			},
			expected: &model.PromptDebugLog{
				ID:           1,
				PromptID:     100,
				SpaceID:      200,
				PromptKey:    "test_key",
				Version:      "1.0.0",
				InputTokens:  10,
				OutputTokens: 20,
				StartedAt:    ptr.Of(int64(1000)),
				EndedAt:      ptr.Of(int64(2000)),
				CostMs:       ptr.Of(int64(1000)),
				StatusCode:   ptr.Of(int32(200)),
				DebuggedBy:   ptr.Of("test_user"),
				DebugID:      1,
				DebugStep:    1,
			},
		},
		{
			name: "debug log with zero values",
			do: &entity.DebugLog{
				ID:           1,
				PromptID:     100,
				SpaceID:      200,
				PromptKey:    "test_key",
				Version:      "1.0.0",
				InputTokens:  10,
				OutputTokens: 20,
				StartedAt:    time.UnixMilli(0),
				EndedAt:      time.UnixMilli(0),
				CostMS:       0,
				StatusCode:   0,
				DebuggedBy:   "",
				DebugID:      1,
				DebugStep:    1,
			},
			expected: &model.PromptDebugLog{
				ID:           1,
				PromptID:     100,
				SpaceID:      200,
				PromptKey:    "test_key",
				Version:      "1.0.0",
				InputTokens:  10,
				OutputTokens: 20,
				StartedAt:    ptr.Of(int64(0)),
				EndedAt:      ptr.Of(int64(0)),
				CostMs:       ptr.Of(int64(0)),
				StatusCode:   ptr.Of(int32(0)),
				DebuggedBy:   ptr.Of(""),
				DebugID:      1,
				DebugStep:    1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DebugLogDO2PO(tt.do)
			assert.Equal(t, tt.expected, got)
		})
	}
}
