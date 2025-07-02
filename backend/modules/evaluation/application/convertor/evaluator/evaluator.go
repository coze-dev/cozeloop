// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package evaluator

import (
	"github.com/bytedance/gg/gptr"

	commondto "github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/evaluation/domain/common"
	evaluatordto "github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/evaluation/domain/evaluator"
	commonconvertor "github.com/coze-dev/cozeloop/backend/modules/evaluation/application/convertor/common"
	evaluatordo "github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
)

func ConvertEvaluatorDTO2DO(evaluatorDTO *evaluatordto.Evaluator) *evaluatordo.Evaluator {
	// 从DTO转换为DO
	evaluatorDO := &evaluatordo.Evaluator{
		ID:                     evaluatorDTO.GetEvaluatorID(),
		SpaceID:                evaluatorDTO.GetWorkspaceID(),
		Name:                   evaluatorDTO.GetName(),
		Description:            evaluatorDTO.GetDescription(),
		DraftSubmitted:         evaluatorDTO.GetDraftSubmitted(),
		EvaluatorType:          evaluatordo.EvaluatorType(evaluatorDTO.GetEvaluatorType()),
		LatestVersion:          evaluatorDTO.GetLatestVersion(),
		PromptEvaluatorVersion: nil,
		BaseInfo:               commonconvertor.ConvertBaseInfoDTO2DO(evaluatorDTO.GetBaseInfo()),
	}
	if evaluatorDTO.CurrentVersion != nil {
		switch evaluatorDTO.GetEvaluatorType() {
		case evaluatordto.EvaluatorType_Prompt:
			evaluatorDO.PromptEvaluatorVersion = ConvertPromptEvaluatorVersionDTO2DO(evaluatorDO.ID, evaluatorDO.SpaceID, evaluatorDTO.GetCurrentVersion())
		}
	}
	return evaluatorDO
}

func ConvertEvaluatorDOList2DTO(doList []*evaluatordo.Evaluator) []*evaluatordto.Evaluator {
	dtoList := make([]*evaluatordto.Evaluator, 0, len(doList))
	for _, evaluatorDO := range doList {
		dtoList = append(dtoList, ConvertEvaluatorDO2DTO(evaluatorDO))
	}
	return dtoList
}

// ConvertEvaluatorDO2DTO 将 evaluatordo.Evaluator 转换为 evaluatordto.Evaluator
func ConvertEvaluatorDO2DTO(do *evaluatordo.Evaluator) *evaluatordto.Evaluator {
	if do == nil {
		return nil
	}
	dto := &evaluatordto.Evaluator{
		EvaluatorID:    gptr.Of(do.ID),
		WorkspaceID:    gptr.Of(do.SpaceID),
		Name:           gptr.Of(do.Name),
		Description:    gptr.Of(do.Description),
		DraftSubmitted: gptr.Of(do.DraftSubmitted),
		EvaluatorType:  evaluatordto.EvaluatorTypePtr(evaluatordto.EvaluatorType(do.EvaluatorType)),
		LatestVersion:  gptr.Of(do.LatestVersion),
		BaseInfo:       commonconvertor.ConvertBaseInfoDO2DTO(do.BaseInfo),
	}

	switch do.EvaluatorType {
	case evaluatordo.EvaluatorTypePrompt:
		if do.PromptEvaluatorVersion != nil {
			versionDTO := ConvertPromptEvaluatorVersionDO2DTO(do.PromptEvaluatorVersion)
			dto.CurrentVersion = versionDTO
		}
	}
	return dto
}

func ConvertPromptEvaluatorVersionDTO2DO(evaluatorID, spaceID int64, dto *evaluatordto.EvaluatorVersion) *evaluatordo.PromptEvaluatorVersion {
	promptEvaluatorVersion := &evaluatordo.PromptEvaluatorVersion{
		ID:                dto.GetID(),
		SpaceID:           spaceID,
		EvaluatorType:     evaluatordo.EvaluatorTypePrompt,
		EvaluatorID:       evaluatorID,
		Description:       dto.GetDescription(),
		Version:           dto.GetVersion(),
		PromptSourceType:  evaluatordo.PromptSourceType(dto.EvaluatorContent.PromptEvaluator.GetPromptSourceType()),
		PromptTemplateKey: dto.EvaluatorContent.PromptEvaluator.GetPromptTemplateKey(),
		BaseInfo:          commonconvertor.ConvertBaseInfoDTO2DO(dto.GetBaseInfo()),
	}
	if dto.EvaluatorContent != nil {
		promptEvaluatorVersion.ReceiveChatHistory = dto.EvaluatorContent.ReceiveChatHistory
		if len(dto.EvaluatorContent.InputSchemas) > 0 {
			promptEvaluatorVersion.InputSchemas = make([]*evaluatordo.ArgsSchema, 0)
			for _, v := range dto.EvaluatorContent.InputSchemas {
				args := commonconvertor.ConvertArgsSchemaDTO2DO(v)
				promptEvaluatorVersion.InputSchemas = append(promptEvaluatorVersion.InputSchemas, args)
			}
		}
		if dto.EvaluatorContent.PromptEvaluator != nil {
			promptEvaluatorVersion.PromptSourceType = evaluatordo.PromptSourceType(dto.EvaluatorContent.PromptEvaluator.GetPromptSourceType())
			promptEvaluatorVersion.PromptTemplateKey = dto.EvaluatorContent.PromptEvaluator.GetPromptTemplateKey()
			promptEvaluatorVersion.MessageList = make([]*evaluatordo.Message, 0)
			for _, originMessage := range dto.EvaluatorContent.PromptEvaluator.GetMessageList() {
				message := commonconvertor.ConvertMessageDTO2DO(originMessage)
				promptEvaluatorVersion.MessageList = append(promptEvaluatorVersion.MessageList, message)
			}
			promptEvaluatorVersion.ModelConfig = commonconvertor.ConvertModelConfigDTO2DO(dto.EvaluatorContent.PromptEvaluator.ModelConfig)
			promptEvaluatorVersion.Tools = make([]*evaluatordo.Tool, 0)
			for _, doTool := range dto.EvaluatorContent.PromptEvaluator.Tools {
				promptEvaluatorVersion.Tools = append(promptEvaluatorVersion.Tools, ConvertToolDTO2DO(doTool))
			}
		}
	}
	return promptEvaluatorVersion
}

// ConvertPromptEvaluatorVersionDO2DTO 将 prompt.PromptEvaluatorVersion 转换为 evaluatordto.EvaluatorVersion
func ConvertPromptEvaluatorVersionDO2DTO(do *evaluatordo.PromptEvaluatorVersion) *evaluatordto.EvaluatorVersion {
	if do == nil {
		return nil
	}
	dto := &evaluatordto.EvaluatorVersion{
		ID:          gptr.Of(do.ID),
		Version:     gptr.Of(do.Version),
		Description: gptr.Of(do.Description),
		BaseInfo:    commonconvertor.ConvertBaseInfoDO2DTO(do.BaseInfo),
		EvaluatorContent: &evaluatordto.EvaluatorContent{
			ReceiveChatHistory: do.ReceiveChatHistory,
			PromptEvaluator: &evaluatordto.PromptEvaluator{
				ModelConfig:       commonconvertor.ConvertModelConfigDO2DTO(do.ModelConfig),
				PromptSourceType:  evaluatordto.PromptSourceTypePtr(evaluatordto.PromptSourceType(do.PromptSourceType)),
				PromptTemplateKey: gptr.Of(do.PromptTemplateKey),
			},
		},
	}
	if len(do.InputSchemas) > 0 {
		dto.EvaluatorContent.InputSchemas = make([]*commondto.ArgsSchema, 0, len(do.InputSchemas))
		for _, v := range do.InputSchemas {
			dto.EvaluatorContent.InputSchemas = append(dto.EvaluatorContent.InputSchemas, commonconvertor.ConvertArgsSchemaDO2DTO(v))
		}
	}
	if len(do.MessageList) > 0 {
		dto.EvaluatorContent.PromptEvaluator.MessageList = make([]*commondto.Message, 0, len(do.MessageList))
		for _, v := range do.MessageList {
			dto.EvaluatorContent.PromptEvaluator.MessageList = append(dto.EvaluatorContent.PromptEvaluator.MessageList, commonconvertor.ConvertMessageDO2DTO(v))
		}
	}
	if len(do.Tools) > 0 {
		dto.EvaluatorContent.PromptEvaluator.Tools = make([]*evaluatordto.Tool, 0, len(do.Tools))
		for _, v := range do.Tools {
			dto.EvaluatorContent.PromptEvaluator.Tools = append(dto.EvaluatorContent.PromptEvaluator.Tools, ConvertToolDO2DTO(v))
		}
	}

	return dto
}
