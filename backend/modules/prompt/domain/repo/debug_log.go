// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/modules/prompt/domain/entity"
)

//go:generate mockgen -destination=mocks/debug_log_repo.go -package=mocks . IDebugLogRepo
type IDebugLogRepo interface {
	SaveDebugLog(ctx context.Context, debugLog *entity.DebugLog) (err error)
	ListDebugHistory(ctx context.Context, param ListDebugHistoryParam) (r *ListDebugHistoryResult, err error)
}

type ListDebugHistoryParam struct {
	PromptID  int64
	UserID    string
	DaysLimit int32
	PageSize  int32
	PageToken *int64
}

type ListDebugHistoryResult struct {
	DebugHistory  []*entity.DebugLog
	NextPageToken int64
	HasMore       bool
}
