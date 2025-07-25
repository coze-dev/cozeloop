// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

type LoopPrompt struct {
	PromptID     int64
	Version      string
	PromptKey    string       `json:"-"`
	Name         string       `json:"-"`
	SubmitStatus SubmitStatus `json:"-"`
	Description  string       `json:"-"`
}

type SubmitStatus int64

const (
	SubmitStatus_Undefined SubmitStatus = 0
	// 未提交
	SubmitStatus_UnSubmit SubmitStatus = 1
	// 已提交
	SubmitStatus_Submitted SubmitStatus = 2
)
