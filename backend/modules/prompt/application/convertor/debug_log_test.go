// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/prompt/domain/prompt"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/entity"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
)

func TestBatchDebugLogDO2DTO(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		dos  []*entity.DebugLog
		dtos []*prompt.DebugLog
	}{
		{
			name: "nil input",
			dos:  nil,
			dtos: nil,
		},
		{
			name: "empty array",
			dos:  []*entity.DebugLog{},
			dtos: []*prompt.DebugLog{},
		},
		{
			name: "array with nil element",
			dos:  []*entity.DebugLog{nil},
			dtos: []*prompt.DebugLog{},
		},
		{
			name: "single debug log",
			dos: []*entity.DebugLog{
				{
					ID:           123,
					PromptID:     456,
					SpaceID:      789,
					PromptKey:    "test_prompt",
					Version:      "1.0.0",
					InputTokens:  100,
					OutputTokens: 200,
					StartedAt:    now,
					EndedAt:      now.Add(time.Second * 2),
					CostMS:       2000,
					StatusCode:   0,
					DebuggedBy:   "test_user",
					DebugID:      1001,
					DebugStep:    1,
				},
			},
			dtos: []*prompt.DebugLog{
				{
					ID:           ptr.Of(int64(123)),
					PromptID:     ptr.Of(int64(456)),
					WorkspaceID:  ptr.Of(int64(789)),
					PromptKey:    ptr.Of("test_prompt"),
					Version:      ptr.Of("1.0.0"),
					InputTokens:  ptr.Of(int64(100)),
					OutputTokens: ptr.Of(int64(200)),
					StartedAt:    ptr.Of(now.UnixMilli()),
					EndedAt:      ptr.Of(now.Add(time.Second * 2).UnixMilli()),
					CostMs:       ptr.Of(int64(2000)),
					StatusCode:   ptr.Of(int32(0)),
					DebuggedBy:   ptr.Of("test_user"),
					DebugID:      ptr.Of(int64(1001)),
					DebugStep:    ptr.Of(int32(1)),
				},
			},
		},
		{
			name: "multiple debug logs",
			dos: []*entity.DebugLog{
				{
					ID:           123,
					PromptID:     456,
					SpaceID:      789,
					PromptKey:    "test_prompt",
					Version:      "1.0.0",
					InputTokens:  100,
					OutputTokens: 200,
					StartedAt:    now,
					EndedAt:      now.Add(time.Second * 2),
					CostMS:       2000,
					StatusCode:   0,
					DebuggedBy:   "test_user",
					DebugID:      1001,
					DebugStep:    1,
				},
				{
					ID:           124,
					PromptID:     456,
					SpaceID:      789,
					PromptKey:    "test_prompt",
					Version:      "1.0.0",
					InputTokens:  150,
					OutputTokens: 250,
					StartedAt:    now.Add(time.Second * 3),
					EndedAt:      now.Add(time.Second * 6),
					CostMS:       3000,
					StatusCode:   0,
					DebuggedBy:   "test_user",
					DebugID:      1002,
					DebugStep:    2,
				},
			},
			dtos: []*prompt.DebugLog{
				{
					ID:           ptr.Of(int64(123)),
					PromptID:     ptr.Of(int64(456)),
					WorkspaceID:  ptr.Of(int64(789)),
					PromptKey:    ptr.Of("test_prompt"),
					Version:      ptr.Of("1.0.0"),
					InputTokens:  ptr.Of(int64(100)),
					OutputTokens: ptr.Of(int64(200)),
					StartedAt:    ptr.Of(now.UnixMilli()),
					EndedAt:      ptr.Of(now.Add(time.Second * 2).UnixMilli()),
					CostMs:       ptr.Of(int64(2000)),
					StatusCode:   ptr.Of(int32(0)),
					DebuggedBy:   ptr.Of("test_user"),
					DebugID:      ptr.Of(int64(1001)),
					DebugStep:    ptr.Of(int32(1)),
				},
				{
					ID:           ptr.Of(int64(124)),
					PromptID:     ptr.Of(int64(456)),
					WorkspaceID:  ptr.Of(int64(789)),
					PromptKey:    ptr.Of("test_prompt"),
					Version:      ptr.Of("1.0.0"),
					InputTokens:  ptr.Of(int64(150)),
					OutputTokens: ptr.Of(int64(250)),
					StartedAt:    ptr.Of(now.Add(time.Second * 3).UnixMilli()),
					EndedAt:      ptr.Of(now.Add(time.Second * 6).UnixMilli()),
					CostMs:       ptr.Of(int64(3000)),
					StatusCode:   ptr.Of(int32(0)),
					DebuggedBy:   ptr.Of("test_user"),
					DebugID:      ptr.Of(int64(1002)),
					DebugStep:    ptr.Of(int32(2)),
				},
			},
		},
		{
			name: "mixed valid and nil elements",
			dos: []*entity.DebugLog{
				nil,
				{
					ID:           123,
					PromptID:     456,
					SpaceID:      789,
					PromptKey:    "test_prompt",
					Version:      "1.0.0",
					InputTokens:  100,
					OutputTokens: 200,
					StartedAt:    now,
					EndedAt:      now.Add(time.Second * 2),
					CostMS:       2000,
					StatusCode:   0,
					DebuggedBy:   "test_user",
					DebugID:      1001,
					DebugStep:    1,
				},
				nil,
			},
			dtos: []*prompt.DebugLog{
				{
					ID:           ptr.Of(int64(123)),
					PromptID:     ptr.Of(int64(456)),
					WorkspaceID:  ptr.Of(int64(789)),
					PromptKey:    ptr.Of("test_prompt"),
					Version:      ptr.Of("1.0.0"),
					InputTokens:  ptr.Of(int64(100)),
					OutputTokens: ptr.Of(int64(200)),
					StartedAt:    ptr.Of(now.UnixMilli()),
					EndedAt:      ptr.Of(now.Add(time.Second * 2).UnixMilli()),
					CostMs:       ptr.Of(int64(2000)),
					StatusCode:   ptr.Of(int32(0)),
					DebuggedBy:   ptr.Of("test_user"),
					DebugID:      ptr.Of(int64(1001)),
					DebugStep:    ptr.Of(int32(1)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := BatchDebugLogDO2DTO(tt.dos)
			assert.Equal(t, tt.dtos, result)
		})
	}
}
