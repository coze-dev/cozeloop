// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	druntime "github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/llm/domain/runtime"
	"github.com/coze-dev/cozeloop/backend/modules/llm/domain/entity"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/slices"
)

func ModelAndTools2OptionDOs(modelCfg *druntime.ModelConfig, tools []*druntime.Tool) []entity.Option {
	var opts []entity.Option
	if modelCfg != nil {
		if modelCfg.Temperature != nil {
			opts = append(opts, entity.WithTemperature(float32(*modelCfg.Temperature)))
		}
		if modelCfg.MaxTokens != nil {
			opts = append(opts, entity.WithMaxTokens(int(*modelCfg.MaxTokens)))
		}
		if modelCfg.TopP != nil {
			opts = append(opts, entity.WithTopP(float32(*modelCfg.TopP)))
		}
		if len(modelCfg.Stop) > 0 {
			opts = append(opts, entity.WithStop(modelCfg.Stop))
		}
		if modelCfg.ToolChoice != nil {
			opts = append(opts, entity.WithToolChoice(ToolChoiceDTO2DO(modelCfg.ToolChoice)))
		}
	}
	if len(tools) > 0 {
		toolsDTO := slices.Transform(tools, func(t *druntime.Tool, _ int) *entity.ToolInfo {
			return ToolDTO2DO(t)
		})
		opts = append(opts, entity.WithTools(toolsDTO))
	}
	return opts
}

func ToolsDTO2DO(ts []*druntime.Tool) []*entity.ToolInfo {
	return slices.Transform(ts, func(t *druntime.Tool, _ int) *entity.ToolInfo {
		return ToolDTO2DO(t)
	})
}

func ToolDTO2DO(t *druntime.Tool) *entity.ToolInfo {
	if t == nil {
		return nil
	}
	return &entity.ToolInfo{
		Name:        t.GetName(),
		Desc:        t.GetDesc(),
		ToolDefType: entity.ToolDefType(t.GetDefType()),
		Def:         t.GetDef(),
	}
}

func ToolChoiceDTO2DO(tc *druntime.ToolChoice) *entity.ToolChoice {
	if tc == nil {
		return nil
	}
	return ptr.Of(entity.ToolChoice(*tc))
}
