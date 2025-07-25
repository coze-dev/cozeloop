// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

// IEvaluatorVersion 定义 Evaluator 的 DO 接口
//
//go:generate mockgen -destination mocks/evaluator_version_mock.go -package mocks . IEvaluatorVersion
type IEvaluatorVersion interface {
	SetID(id int64)
	GetID() int64
	SetEvaluatorID(evaluatorID int64)
	GetEvaluatorID() int64
	SetVersion(version string)
	GetVersion() string
	SetSpaceID(spaceID int64)
	GetSpaceID() int64
	SetDescription(description string)
	GetDescription() string
	SetBaseInfo(baseInfo *BaseInfo)
	GetBaseInfo() *BaseInfo
	SetTools(tools []*Tool)
	GetPromptTemplateKey() string
	SetPromptSuffix(promptSuffix string)
	GetModelConfig() *ModelConfig
	SetParseType(parseType ParseType)

	ValidateInput(input *EvaluatorInputData) error
	ValidateBaseInfo() error
}
