// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convert

import (
	"github.com/bytedance/gg/gptr"

	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/infra/repo/experiment/mysql/gorm_gen/model"
)

func NewExptStatsConverter() ExptStatsConverter {
	return ExptStatsConverter{}
}

type ExptStatsConverter struct{}

func (ExptStatsConverter) DO2PO(stats *entity.ExptStats) *model.ExptStats {
	return &model.ExptStats{
		ID:              stats.ID,
		SpaceID:         stats.SpaceID,
		ExptID:          stats.ExptID,
		PendingCnt:      stats.PendingTurnCnt,
		SuccessCnt:      stats.SuccessTurnCnt,
		FailCnt:         stats.FailTurnCnt,
		TerminatedCnt:   stats.TerminatedTurnCnt,
		ProcessingCnt:   stats.ProcessingTurnCnt,
		CreditCost:      stats.CreditCost,
		InputTokenCost:  gptr.Of(stats.InputTokenCost),
		OutputTokenCost: gptr.Of(stats.OutputTokenCost),
	}
}

func (ExptStatsConverter) PO2DO(stats *model.ExptStats) *entity.ExptStats {
	return &entity.ExptStats{
		ID:                stats.ID,
		SpaceID:           stats.SpaceID,
		ExptID:            stats.ExptID,
		PendingTurnCnt:    stats.PendingCnt,
		SuccessTurnCnt:    stats.SuccessCnt,
		FailTurnCnt:       stats.FailCnt,
		TerminatedTurnCnt: stats.TerminatedCnt,
		ProcessingTurnCnt: stats.ProcessingCnt,
		CreditCost:        stats.CreditCost,
		InputTokenCost:    gptr.Indirect(stats.InputTokenCost),
		OutputTokenCost:   gptr.Indirect(stats.OutputTokenCost),
	}
}
