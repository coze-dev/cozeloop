// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	prompterr "github.com/coze-dev/coze-loop/backend/modules/prompt/pkg/errno"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-loop/backend/pkg/unittest"
)

func TestPromptTemplate_formatMessages(t *testing.T) {
	tests := []struct {
		name          string
		template      *PromptTemplate
		messages      []*Message
		variableVals  []*VariableVal
		expectedMsgs  []*Message
		expectedError error
	}{
		{
			name:          "nil template",
			template:      nil,
			messages:      []*Message{},
			variableVals:  []*VariableVal{},
			expectedMsgs:  nil,
			expectedError: nil,
		},
		{
			name: "empty messages",
			template: &PromptTemplate{
				TemplateType: TemplateTypeNormal,
				Messages: []*Message{
					{
						Role:    RoleSystem,
						Content: ptr.Of("You are a helpful assistant."),
					},
				},
			},
			messages:     []*Message{},
			variableVals: []*VariableVal{},
			expectedMsgs: []*Message{
				{
					Role:    RoleSystem,
					Content: ptr.Of("You are a helpful assistant."),
				},
			},
			expectedError: nil,
		},
		{
			name: "nil variable values",
			template: &PromptTemplate{
				TemplateType: TemplateTypeNormal,
				Messages: []*Message{
					{
						Role:    RoleSystem,
						Content: ptr.Of("You are a {{role}}."),
					},
				},
				VariableDefs: []*VariableDef{
					{
						Key:  "role",
						Desc: "role",
						Type: VariableTypeString,
					},
				},
			},
			messages: []*Message{
				{
					Role:    RoleUser,
					Content: ptr.Of("Hello"),
				},
			},
			variableVals: nil,
			expectedMsgs: []*Message{
				{
					Role:    RoleSystem,
					Content: ptr.Of("You are a ."),
				},
				{
					Role:    RoleUser,
					Content: ptr.Of("Hello"),
				},
			},
			expectedError: nil,
		},
		{
			name: "nil variable defs",
			template: &PromptTemplate{
				TemplateType: TemplateTypeNormal,
				Messages: []*Message{
					{
						Role:    RoleSystem,
						Content: ptr.Of("You are a {{role}}."),
					},
				},
				VariableDefs: nil,
			},
			messages: []*Message{
				{
					Role:    RoleUser,
					Content: ptr.Of("Hello"),
				},
			},
			variableVals: []*VariableVal{
				{
					Key:   "role",
					Value: ptr.Of("helpful assistant"),
				},
			},
			expectedMsgs: []*Message{
				{
					Role:    RoleSystem,
					Content: ptr.Of("You are a {{role}}."),
				},
				{
					Role:    RoleUser,
					Content: ptr.Of("Hello"),
				},
			},
			expectedError: nil,
		},
		{
			name: "placeholder role with valid messages",
			template: &PromptTemplate{
				TemplateType: TemplateTypeNormal,
				Messages: []*Message{
					{
						Role:    RolePlaceholder,
						Content: ptr.Of("greeting"),
					},
				},
			},
			messages: []*Message{},
			variableVals: []*VariableVal{
				{
					Key: "greeting",
					PlaceholderMessages: []*Message{
						{
							Role:    RoleSystem,
							Content: ptr.Of("Hello!"),
						},
					},
				},
			},
			expectedMsgs: []*Message{
				{
					Role:    RoleSystem,
					Content: ptr.Of("Hello!"),
				},
			},
			expectedError: nil,
		},
		{
			name: "placeholder role with invalid message role",
			template: &PromptTemplate{
				TemplateType: TemplateTypeNormal,
				Messages: []*Message{
					{
						Role:    RolePlaceholder,
						Content: ptr.Of("greeting"),
					},
				},
			},
			messages: []*Message{},
			variableVals: []*VariableVal{
				{
					Key: "greeting",
					PlaceholderMessages: []*Message{
						{
							Role:    RolePlaceholder,
							Content: ptr.Of("Hello!"),
						},
					},
				},
			},
			expectedMsgs:  nil,
			expectedError: errorx.NewByCode(prompterr.CommonInvalidParamCode),
		},
		{
			name: "normal message with variable replacement",
			template: &PromptTemplate{
				TemplateType: TemplateTypeNormal,
				Messages: []*Message{
					{
						Role:    RoleSystem,
						Content: ptr.Of("You are a {{role}}."),
					},
				},
				VariableDefs: []*VariableDef{
					{
						Key:  "role",
						Desc: "role",
						Type: VariableTypeString,
					},
				},
			},
			messages: []*Message{
				{
					Role:    RoleUser,
					Content: ptr.Of("Hello"),
				},
			},
			variableVals: []*VariableVal{
				{
					Key:   "role",
					Value: ptr.Of("helpful assistant"),
				},
			},
			expectedMsgs: []*Message{
				{
					Role:    RoleSystem,
					Content: ptr.Of("You are a helpful assistant."),
				},
				{
					Role:    RoleUser,
					Content: ptr.Of("Hello"),
				},
			},
			expectedError: nil,
		},
		{
			name: "message with parts",
			template: &PromptTemplate{
				TemplateType: TemplateTypeNormal,
				Messages: []*Message{
					{
						Role:    RoleSystem,
						Content: ptr.Of("You are a {{role}}."),
						Parts: []*ContentPart{
							{
								Type: ContentTypeText,
								Text: ptr.Of("Additional info: {{info}}"),
							},
						},
					},
				},
				VariableDefs: []*VariableDef{
					{
						Key:  "role",
						Desc: "role",
						Type: VariableTypeString,
					},
					{
						Key:  "info",
						Desc: "info",
						Type: VariableTypeString,
					},
				},
			},
			messages: []*Message{
				{
					Role:    RoleUser,
					Content: ptr.Of("Hello"),
				},
			},
			variableVals: []*VariableVal{
				{
					Key:   "role",
					Value: ptr.Of("helpful assistant"),
				},
				{
					Key:   "info",
					Value: ptr.Of("some info"),
				},
			},
			expectedMsgs: []*Message{
				{
					Role:    RoleSystem,
					Content: ptr.Of("You are a helpful assistant."),
					Parts: []*ContentPart{
						{
							Type: ContentTypeText,
							Text: ptr.Of("Additional info: some info"),
						},
					},
				},
				{
					Role:    RoleUser,
					Content: ptr.Of("Hello"),
				},
			},
			expectedError: nil,
		},
		{
			name: "message with empty content",
			template: &PromptTemplate{
				TemplateType: TemplateTypeNormal,
				Messages: []*Message{
					{
						Role:    RoleSystem,
						Content: ptr.Of(""),
					},
				},
			},
			messages: []*Message{
				{
					Role:    RoleUser,
					Content: ptr.Of("Hello"),
				},
			},
			variableVals: []*VariableVal{},
			expectedMsgs: []*Message{
				{
					Role:    RoleSystem,
					Content: ptr.Of(""),
				},
				{
					Role:    RoleUser,
					Content: ptr.Of("Hello"),
				},
			},
			expectedError: nil,
		},
		{
			name: "message with nil content",
			template: &PromptTemplate{
				TemplateType: TemplateTypeNormal,
				Messages: []*Message{
					{
						Role: RoleSystem,
					},
				},
			},
			messages: []*Message{
				{
					Role:    RoleUser,
					Content: ptr.Of("Hello"),
				},
			},
			variableVals: []*VariableVal{},
			expectedMsgs: []*Message{
				{
					Role: RoleSystem,
				},
				{
					Role:    RoleUser,
					Content: ptr.Of("Hello"),
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formattedMsgs, err := tt.template.formatMessages(tt.messages, tt.variableVals)
			unittest.AssertErrorEqual(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedMsgs, formattedMsgs)
		})
	}
}

func TestCmpEqual(t *testing.T) {
	var pd1 *PromptDetail
	var pd2 *PromptDetail
	fmt.Printf("nil cmp nil = %t\n", cmp.Equal(pd1, pd2))              // true
	fmt.Printf("nil cmp !nil = %t\n", cmp.Equal(pd1, &PromptDetail{})) // false
	fmt.Printf("!nil cmp nil = %t\n", cmp.Equal(&PromptDetail{}, pd2)) // false
}
