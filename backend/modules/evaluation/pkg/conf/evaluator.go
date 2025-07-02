// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"context"

	"github.com/samber/lo"

	"github.com/coze-dev/cozeloop/backend/infra/limiter"
	evaluatordto "github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/evaluation/domain/evaluator"
	"github.com/coze-dev/cozeloop/backend/pkg/conf"
)

//go:generate mockgen -destination=mocks/evaluator_configer.go -package=mocks . IConfiger
type IConfiger interface {
	GetEvaluatorTemplateConf(ctx context.Context) (etf map[string]map[string]*evaluatordto.EvaluatorContent)
	GetEvaluatorToolConf(ctx context.Context) (etf map[string]*evaluatordto.Tool) // tool_key -> tool
	GetRateLimiterConf(ctx context.Context) (rlc []limiter.Rule)
	GetEvaluatorToolMapping(ctx context.Context) (etf map[string]string)            // prompt_template_key -> tool_key
	GetEvaluatorPromptSuffix(ctx context.Context) (suffix map[string]string)        // suffix_key -> suffix
	GetEvaluatorPromptSuffixMapping(ctx context.Context) (suffix map[string]string) // model_id -> suffix_key
}

func NewEvaluatorConfiger(configFactory conf.IConfigLoaderFactory) IConfiger {
	loader, err := configFactory.NewConfigLoader("evaluation.yaml")
	if err != nil {
		return nil
	}
	return &configer{
		loader: loader,
	}
}

func (c *configer) GetEvaluatorTemplateConf(ctx context.Context) (etf map[string]map[string]*evaluatordto.EvaluatorContent) {
	const key = "evaluator_template_conf"
	etf = make(map[string]map[string]*evaluatordto.EvaluatorContent)
	lo.Ternary(c.loader.UnmarshalKey(ctx, key, &etf) == nil, etf, DefaultEvaluatorTemplateConf())
	return etf
}

func DefaultEvaluatorTemplateConf() map[string]map[string]*evaluatordto.EvaluatorContent {
	return map[string]map[string]*evaluatordto.EvaluatorContent{}
}

func (c *configer) GetEvaluatorToolConf(ctx context.Context) (etf map[string]*evaluatordto.Tool) {
	const key = "evaluator_tool_conf"
	return lo.Ternary(c.loader.UnmarshalKey(ctx, key, &etf) == nil, etf, DefaultEvaluatorToolConf())
}

func DefaultEvaluatorToolConf() map[string]*evaluatordto.Tool {
	return make(map[string]*evaluatordto.Tool, 0)
}

func (c *configer) GetRateLimiterConf(ctx context.Context) (rlc []limiter.Rule) {
	const key = "rate_limiter_conf"
	return lo.Ternary(c.loader.UnmarshalKey(ctx, key, &rlc) == nil, rlc, DefaultRateLimiterConf())
}

func DefaultRateLimiterConf() []limiter.Rule {
	return make([]limiter.Rule, 0)
}

func (c *configer) GetEvaluatorToolMapping(ctx context.Context) (etf map[string]string) {
	const key = "evaluator_tool_mapping"
	return lo.Ternary(c.loader.UnmarshalKey(ctx, key, &etf) == nil, etf, DefaultEvaluatorToolMapping())
}

func DefaultEvaluatorToolMapping() map[string]string {
	return make(map[string]string)
}

func (c *configer) GetEvaluatorPromptSuffix(ctx context.Context) (suffix map[string]string) {
	const key = "evaluator_prompt_suffix"
	return lo.Ternary(c.loader.UnmarshalKey(ctx, key, &suffix) == nil, suffix, DefaultEvaluatorPromptSuffix())
}

func DefaultEvaluatorPromptSuffix() map[string]string {
	return make(map[string]string)
}

func (c *configer) GetEvaluatorPromptSuffixMapping(ctx context.Context) (suffix map[string]string) {
	const key = "evaluator_prompt_mapping"
	return lo.Ternary(c.loader.UnmarshalKey(ctx, key, &suffix) == nil, suffix, DefaultEvaluatorPromptMapping())
}

func DefaultEvaluatorPromptMapping() map[string]string {
	return make(map[string]string)
}
