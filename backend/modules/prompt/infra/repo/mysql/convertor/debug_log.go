// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"time"

	"github.com/coze-dev/coze-loop/backend/modules/prompt/domain/entity"
	"github.com/coze-dev/coze-loop/backend/modules/prompt/infra/repo/mysql/gorm_gen/model"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
)

func DebugLogsPO2DO(pos []*model.PromptDebugLog) []*entity.DebugLog {
	if pos == nil {
		return nil
	}
	debugLogs := make([]*entity.DebugLog, 0, len(pos))
	for _, po := range pos {
		if po == nil {
			continue
		}
		debugLogs = append(debugLogs, DebugLogPO2DO(po))
	}
	return debugLogs
}

func DebugLogPO2DO(po *model.PromptDebugLog) *entity.DebugLog {
	if po == nil {
		return nil
	}
	return &entity.DebugLog{
		ID:           po.ID,
		PromptID:     po.PromptID,
		SpaceID:      po.SpaceID,
		PromptKey:    po.PromptKey,
		Version:      po.Version,
		InputTokens:  po.InputTokens,
		OutputTokens: po.OutputTokens,
		StartedAt:    time.UnixMilli(ptr.From(po.StartedAt)),
		EndedAt:      time.UnixMilli(ptr.From(po.EndedAt)),
		CostMS:       ptr.From(po.CostMs),
		StatusCode:   ptr.From(po.StatusCode),
		DebuggedBy:   ptr.From(po.DebuggedBy),
		DebugID:      po.DebugID,
		DebugStep:    po.DebugStep,
	}
}

//============================================================

func DebugLogDO2PO(do *entity.DebugLog) *model.PromptDebugLog {
	if do == nil {
		return nil
	}
	return &model.PromptDebugLog{
		ID:           do.ID,
		PromptID:     do.PromptID,
		SpaceID:      do.SpaceID,
		PromptKey:    do.PromptKey,
		Version:      do.Version,
		InputTokens:  do.InputTokens,
		OutputTokens: do.OutputTokens,
		StartedAt:    ptr.Of(do.StartedAt.UnixMilli()),
		EndedAt:      ptr.Of(do.EndedAt.UnixMilli()),
		CostMs:       ptr.Of(do.CostMS),
		StatusCode:   ptr.Of(do.StatusCode),
		DebuggedBy:   ptr.Of(do.DebuggedBy),
		DebugID:      do.DebugID,
		DebugStep:    do.DebugStep,
	}
}
