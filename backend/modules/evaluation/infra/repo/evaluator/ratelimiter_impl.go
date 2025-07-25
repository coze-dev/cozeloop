// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package evaluator

import (
	"context"

	"github.com/coze-dev/cozeloop/backend/infra/limiter"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/consts"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/repo"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/pkg/conf"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

type RateLimiterImpl struct {
	limiter limiter.IRateLimiter
}

func NewRateLimiterImpl(ctx context.Context, limiterFactory limiter.IRateLimiterFactory, evaluatorConfiger conf.IConfiger) repo.RateLimiter {
	return &RateLimiterImpl{
		limiter: limiterFactory.NewRateLimiter(limiter.WithRules(evaluatorConfiger.GetRateLimiterConf(ctx)...)),
	}
}

func (s *RateLimiterImpl) AllowInvoke(ctx context.Context, spaceID int64) bool {
	tags := []limiter.Tag{
		{K: "space_id", V: spaceID},
	}
	res, err := s.limiter.AllowN(ctx, consts.RateLimitBizKeyEvaluator, 1, limiter.WithTags(tags...))
	if err != nil {
		logs.CtxError(ctx, "allow invoke failed, err=%v", err)
		return true
	}
	if res.Allowed {
		logs.CtxInfo(ctx, "[AllowInvoke] allow invoke")
		return true
	}
	logs.CtxInfo(ctx, "[AllowInvoke] not allow invoke")
	return false
}
