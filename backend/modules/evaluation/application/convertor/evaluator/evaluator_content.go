// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package evaluator

import (
	"github.com/bytedance/gg/gptr"

	evaluatordto "github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/evaluation/domain/evaluator"
	commonentity "github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
)

// ConvertToolDTO2DO 将 DTO 转换为 Tool 结构体
func ConvertToolDTO2DO(dto *evaluatordto.Tool) *commonentity.Tool {
	if dto == nil {
		return nil
	}
	return &commonentity.Tool{
		Type:     commonentity.ToolType(dto.Type),
		Function: ConvertFunctionDTO2DO(dto.Function),
	}
}

// ConvertToolDO2DTO 将 Tool 结构体转换为 DTO
func ConvertToolDO2DTO(do *commonentity.Tool) *evaluatordto.Tool {
	if do == nil {
		return nil
	}
	return &evaluatordto.Tool{
		Type:     evaluatordto.ToolType(do.Type),
		Function: ConvertFunctionDO2DTO(do.Function),
	}
}

// ConvertFunctionDTO2DO 将 DTO 转换为 Function 结构体
func ConvertFunctionDTO2DO(dto *evaluatordto.Function) *commonentity.Function {
	if dto == nil {
		return nil
	}
	description := ""
	if dto.Description != nil {
		description = *dto.Description
	}
	parameters := ""
	if dto.Parameters != nil {
		parameters = *dto.Parameters
	}
	return &commonentity.Function{
		Name:        dto.Name,
		Description: description,
		Parameters:  parameters,
	}
}

// ConvertFunctionDO2DTO 将 Function 结构体转换为 DTO
func ConvertFunctionDO2DTO(do *commonentity.Function) *evaluatordto.Function {
	if do == nil {
		return nil
	}
	return &evaluatordto.Function{
		Name:        do.Name,
		Description: gptr.Of(do.Description),
		Parameters:  gptr.Of(do.Parameters),
	}
}
