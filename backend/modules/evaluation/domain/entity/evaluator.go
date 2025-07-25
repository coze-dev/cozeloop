// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

type Evaluator struct {
	ID             int64
	SpaceID        int64
	Name           string
	Description    string
	DraftSubmitted bool
	EvaluatorType  EvaluatorType
	LatestVersion  string
	BaseInfo       *BaseInfo

	PromptEvaluatorVersion *PromptEvaluatorVersion
}

type EvaluatorType int64

const (
	EvaluatorTypePrompt EvaluatorType = 1
	EvaluatorTypeCode   EvaluatorType = 2
)

var EvaluatorTypeSet = map[EvaluatorType]struct{}{
	EvaluatorTypePrompt: {},
	EvaluatorTypeCode:   {},
}

func (e *Evaluator) GetEvaluatorVersion() IEvaluatorVersion {
	switch e.EvaluatorType {
	case EvaluatorTypePrompt:
		return e.PromptEvaluatorVersion
	default:
		return nil
	}
}

func (e *Evaluator) SetEvaluatorVersion(version *Evaluator) {
	switch e.EvaluatorType {
	case EvaluatorTypePrompt:
		e.PromptEvaluatorVersion = version.PromptEvaluatorVersion
	default:
		return
	}
}
