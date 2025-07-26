// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/prompt/domain/prompt"
	"github.com/coze-dev/coze-loop/backend/modules/prompt/domain/entity"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
)

func BatchDebugLogDO2DTO(dos []*entity.DebugLog) []*prompt.DebugLog {
	if dos == nil {
		return nil
	}
	result := make([]*prompt.DebugLog, 0, len(dos))
	for _, log := range dos {
		if log == nil {
			continue
		}
		result = append(result, DebugLogDO2DTO(log))
	}
	return result
}

func DebugLogDO2DTO(do *entity.DebugLog) *prompt.DebugLog {
	if do == nil {
		return nil
	}
	return &prompt.DebugLog{
		ID:           ptr.Of(do.ID),
		PromptID:     ptr.Of(do.PromptID),
		WorkspaceID:  ptr.Of(do.SpaceID),
		PromptKey:    ptr.Of(do.PromptKey),
		Version:      ptr.Of(do.Version),
		InputTokens:  ptr.Of(do.InputTokens),
		OutputTokens: ptr.Of(do.OutputTokens),
		CostMs:       ptr.Of(do.CostMS),
		StatusCode:   ptr.Of(do.StatusCode),
		DebuggedBy:   ptr.Of(do.DebuggedBy),
		DebugID:      ptr.Of(do.DebugID),
		DebugStep:    ptr.Of(do.DebugStep),
		StartedAt:    ptr.Of(do.StartedAt.UnixMilli()),
		EndedAt:      ptr.Of(do.EndedAt.UnixMilli()),
	}
}
