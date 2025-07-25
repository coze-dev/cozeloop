// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/observability/trace"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/loop_span"
)

func AdvanceInfoDO2DTO(info *loop_span.TraceAdvanceInfo) *trace.TraceAdvanceInfo {
	return &trace.TraceAdvanceInfo{
		TraceID: info.TraceId,
		Tokens: &trace.TokenCost{
			Input:  info.InputCost,
			Output: info.OutputCost,
		},
	}
}

func BatchAdvanceInfoDO2DTO(infos []*loop_span.TraceAdvanceInfo) []*trace.TraceAdvanceInfo {
	ret := make([]*trace.TraceAdvanceInfo, len(infos))
	for i, info := range infos {
		ret[i] = AdvanceInfoDO2DTO(info)
	}
	return ret
}

func FileMetaDO2DTO() {
}
