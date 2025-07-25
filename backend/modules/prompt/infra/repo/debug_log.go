// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"slices"
	"time"

	"golang.org/x/exp/maps"

	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/repo"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/infra/repo/mysql"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/infra/repo/mysql/convertor"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/infra/repo/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
	loopslices "github.com/coze-dev/cozeloop/backend/pkg/lang/slices"
)

type DebugLogRepoImpl struct {
	idgen       idgen.IIDGenerator
	debugLogDAO mysql.IDebugLogDAO
}

func NewDebugLogRepo(
	idgen idgen.IIDGenerator,
	debugLogDao mysql.IDebugLogDAO,
) repo.IDebugLogRepo {
	return &DebugLogRepoImpl{
		idgen:       idgen,
		debugLogDAO: debugLogDao,
	}
}

func (d *DebugLogRepoImpl) SaveDebugLog(ctx context.Context, debugLog *entity.DebugLog) (err error) {
	if debugLog == nil {
		return nil
	}
	id, err := d.idgen.GenID(ctx)
	if err != nil {
		return err
	}
	debugLogPO := convertor.DebugLogDO2PO(debugLog)
	debugLogPO.ID = id
	return d.debugLogDAO.Save(ctx, debugLogPO)
}

func (d *DebugLogRepoImpl) ListDebugHistory(ctx context.Context, param repo.ListDebugHistoryParam) (r *repo.ListDebugHistoryResult, err error) {
	pageSize := int(param.PageSize)
	debugLogs, err := d.debugLogDAO.List(ctx, mysql.ListParam{
		PromptID:    ptr.Of(param.PromptID),
		UserID:      ptr.Of(param.UserID),
		StartBefore: param.PageToken,
		StartAfter:  ptr.Of(time.Now().Add(-1 * time.Duration(param.DaysLimit) * 24 * time.Hour).UnixMilli()),
		// 支持function call多步调试记录串联后，一次完整的bug可能对应多条记录，为了准确分页，每次调试历史只取第一步记录
		DebugStep: ptr.Of(int32(1)),
		Limit:     ptr.Of(pageSize + 1),
	})
	if err != nil {
		return nil, err
	}
	if len(debugLogs) == 0 {
		return &repo.ListDebugHistoryResult{}, nil
	}
	var nextPageToken int64
	var hasMore bool
	if len(debugLogs) > pageSize {
		hasMore = true
		nextPageToken = ptr.From(debugLogs[pageSize].StartedAt)
		debugLogs = debugLogs[:pageSize]
	}
	debugLogMap := loopslices.ToMap(debugLogs, func(e *model.PromptDebugLog) (int64, *model.PromptDebugLog) {
		if e == nil {
			return 0, nil
		}
		return e.DebugID, e
	})
	debugIDs := loopslices.Transform(debugLogs, func(e *model.PromptDebugLog, idx int) int64 {
		if e == nil {
			return 0
		}
		return e.DebugID
	})
	// 查询可能存在的多步调试记录
	allLogs, err := d.debugLogDAO.List(ctx, mysql.ListParam{
		DebugIDs: debugIDs,
	})
	if err != nil {
		return nil, err
	}
	// 用最后一步的结束时间覆盖第一步的结束时间
	for _, stepLog := range allLogs {
		if stepLog == nil {
			continue
		}
		if log, ok := debugLogMap[stepLog.DebugID]; ok && ptr.From(stepLog.EndedAt) > ptr.From(log.EndedAt) {
			log.EndedAt = stepLog.EndedAt
			log.CostMs = ptr.Of(ptr.From(log.EndedAt) - ptr.From(log.StartedAt))
			log.InputTokens += stepLog.InputTokens
			log.OutputTokens += stepLog.OutputTokens
		}
	}
	debugLogs = maps.Values(debugLogMap)
	slices.SortFunc(debugLogs, func(a *model.PromptDebugLog, b *model.PromptDebugLog) int {
		if ptr.From(a.StartedAt) > ptr.From(b.StartedAt) {
			return -1
		}
		return 1
	})
	result := &repo.ListDebugHistoryResult{
		DebugHistory:  convertor.DebugLogsPO2DO(debugLogs),
		NextPageToken: nextPageToken,
		HasMore:       hasMore,
	}
	return result, nil
}
