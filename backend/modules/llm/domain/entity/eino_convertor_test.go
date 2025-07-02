// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"reflect"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
	"github.com/coze-dev/cozeloop/backend/pkg/unittest"
)

func TestFromDOResponseMeta(t *testing.T) {
	type args struct {
		rm *ResponseMeta
	}
	tests := []struct {
		name string
		args args
		want *schema.ResponseMeta
	}{
		{
			name: "TestFromDOResponseMeta",
			args: args{
				rm: &ResponseMeta{
					FinishReason: "stop",
					Usage: &TokenUsage{
						PromptTokens:     100,
						CompletionTokens: 10,
						TotalTokens:      110,
					},
				},
			},
			want: &schema.ResponseMeta{
				FinishReason: "stop",
				Usage: &schema.TokenUsage{
					PromptTokens:     100,
					CompletionTokens: 10,
					TotalTokens:      110,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromDOResponseMeta(tt.args.rm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromDOResponseMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromDOToolCalls(t *testing.T) {
	type args struct {
		ts []*ToolCall
	}
	tests := []struct {
		name string
		args args
		want []schema.ToolCall
	}{
		{
			name: "test from do tool calls",
			args: args{
				ts: []*ToolCall{
					{
						Index: ptr.Of(int64(1)),
						ID:    "id",
						Type:  "function",
						Function: &FunctionCall{
							Name:      "name",
							Arguments: "args",
						},
						Extra: nil,
					},
				},
			},
			want: []schema.ToolCall{
				{
					Index: ptr.Of(1),
					ID:    "id",
					Type:  "function",
					Function: schema.FunctionCall{
						Name:      "name",
						Arguments: "args",
					},
					Extra: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, len(tt.want), len(FromDOToolCalls(tt.args.ts)))
			assert.Equal(t, tt.want[0].Function.Arguments, FromDOToolCalls(tt.args.ts)[0].Function.Arguments)
		})
	}
}

func TestFromDOToolChoice(t *testing.T) {
	type args struct {
		do ToolChoice
	}
	tests := []struct {
		name               string
		args               args
		wantEinoToolChoice schema.ToolChoice
	}{
		{
			name: "TestFromDOToolChoice_none",
			args: args{
				do: ToolChoiceNone,
			},
			wantEinoToolChoice: schema.ToolChoiceForbidden,
		},
		{
			name: "TestFromDOToolChoice_auto",
			args: args{
				do: ToolChoiceAuto,
			},
			wantEinoToolChoice: schema.ToolChoiceAllowed,
		},
		{
			name: "TestFromDOToolChoice_required",
			args: args{
				do: ToolChoiceRequired,
			},
			wantEinoToolChoice: schema.ToolChoiceForced,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEinoToolChoice := FromDOToolChoice(tt.args.do); gotEinoToolChoice != tt.wantEinoToolChoice {
				t.Errorf("FromDOToolChoice() = %v, want %v", gotEinoToolChoice, tt.wantEinoToolChoice)
			}
		})
	}
}

func TestToDOMessages(t *testing.T) {
	type args struct {
		msgs []*schema.Message
	}
	tests := []struct {
		name    string
		args    args
		want    []*Message
		wantErr error
	}{
		{
			name: "success",
			args: args{
				msgs: []*schema.Message{
					{
						Role:    schema.Assistant,
						Content: "there is content",
					},
				},
			},
			want: []*Message{
				{
					Role:    RoleAssistant,
					Content: "there is content",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToDOMessages(tt.args.msgs)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want[0].Content, got[0].Content)
		})
	}
}

func TestToDOToolCalls(t *testing.T) {
	type args struct {
		tcs []schema.ToolCall
	}
	tests := []struct {
		name string
		args args
		want []*ToolCall
	}{
		{
			name: "success",
			args: args{
				tcs: []schema.ToolCall{
					{
						Index: ptr.Of(1),
						ID:    "id",
						Type:  "function",
						Function: schema.FunctionCall{
							Name:      "name",
							Arguments: "args",
						},
						Extra: nil,
					},
				},
			},
			want: []*ToolCall{
				{
					Index: ptr.Of(int64(1)),
					ID:    "id",
					Type:  "function",
					Function: &FunctionCall{
						Name:      "name",
						Arguments: "args",
					},
					Extra: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToDOToolCalls(tt.args.tcs)
			assert.Equal(t, len(tt.want), len(got))
			assert.Equal(t, tt.want[0].Function.Arguments, got[0].Function.Arguments)
		})
	}
}
