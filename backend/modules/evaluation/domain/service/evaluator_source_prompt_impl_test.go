// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/kaptinlin/jsonrepair"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/evaluation/domain/evaluator"
	metricsmocks "github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/component/metrics/mocks"
	rpcmocks "github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/component/rpc/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
	configmocks "github.com/coze-dev/cozeloop/backend/modules/evaluation/pkg/conf/mocks"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
)

// 修正后的 TestEvaluatorSourcePromptServiceImpl_Run 结构：
func TestEvaluatorSourcePromptServiceImpl_Run(t *testing.T) {
	ctrl := gomock.NewController(t) // Controller for the entire test function
	defer ctrl.Finish()

	// These mocks will be shared across all test cases due to the singleton nature of the service
	sharedMockLLMProvider := rpcmocks.NewMockILLMProvider(ctrl)
	sharedMockMetric := metricsmocks.NewMockEvaluatorExecMetrics(ctrl)
	sharedMockConfiger := configmocks.NewMockIConfiger(ctrl)

	// Instantiate the service once with the shared mocks
	// The singleton instance will use these mocks for all subsequent calls in this test function
	service := &EvaluatorSourcePromptServiceImpl{
		llmProvider: sharedMockLLMProvider,
		metric:      sharedMockMetric,
		configer:    sharedMockConfiger,
	}

	ctx := context.Background()
	baseMockEvaluator := &entity.Evaluator{
		ID:            100,
		SpaceID:       1,
		Name:          "Test Evaluator",
		EvaluatorType: entity.EvaluatorTypePrompt,
		PromptEvaluatorVersion: &entity.PromptEvaluatorVersion{
			ID:                100,
			EvaluatorID:       100,
			SpaceID:           1,
			PromptTemplateKey: "test-template-key",
			PromptSuffix:      "test-prompt-suffix",
			ModelConfig: &entity.ModelConfig{
				ModelID: 1,
			},
			ParseType: entity.ParseTypeFunctionCall,
			MessageList: []*entity.Message{
				{
					Role: entity.RoleSystem,
					Content: &entity.Content{
						ContentType: ptr.Of(entity.ContentTypeText),
						Text:        ptr.Of("{{test-content}}"),
					},
				},
			},
			Tools: []*entity.Tool{
				{
					Type: entity.ToolTypeFunction,
					Function: &entity.Function{
						Name:        "test_function",
						Description: "test description",
						Parameters:  "{\"type\": \"object\", \"properties\": {\"score\": {\"type\": \"number\"}, \"reasoning\": {\"type\": \"string\"}}}",
					},
				},
			},
		},
	}

	baseMockInput := &entity.EvaluatorInputData{
		InputFields: map[string]*entity.Content{
			"input": {
				ContentType: ptr.Of(entity.ContentTypeText),
				Text:        ptr.Of("test input"),
			},
		},
	}

	testCases := []struct {
		name            string
		evaluator       *entity.Evaluator
		input           *entity.EvaluatorInputData
		setupMocks      func() // setupMocks now configures the shared mocks
		expectedOutput  *entity.EvaluatorOutputData
		expectedStatus  entity.EvaluatorRunStatus
		checkOutputFunc func(t *testing.T, output *entity.EvaluatorOutputData, expected *entity.EvaluatorOutputData)
	}{
		{
			name:      "成功运行评估器",
			evaluator: baseMockEvaluator,
			input:     baseMockInput,
			setupMocks: func() {
				sharedMockLLMProvider.EXPECT().Call(gomock.Any(), gomock.Any()).Return(
					&entity.ReplyItem{
						ToolCalls: []*entity.ToolCall{
							{
								Type: entity.ToolTypeFunction,
								FunctionCall: &entity.FunctionCall{
									Name:      "test_function",
									Arguments: ptr.Of("{\"score\": 1.0, \"reason\": \"test response\"}"),
								},
							},
						},
						TokenUsage: &entity.TokenUsage{InputTokens: 10, OutputTokens: 10},
					}, nil)
				sharedMockMetric.EXPECT().EmitRun(int64(1), gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedOutput: &entity.EvaluatorOutputData{
				EvaluatorResult:   &entity.EvaluatorResult{Score: ptr.Of(1.0), Reasoning: "test response"},
				EvaluatorUsage:    &entity.EvaluatorUsage{InputTokens: 10, OutputTokens: 10}, // As per original test
				EvaluatorRunError: nil,
			},
			expectedStatus: entity.EvaluatorRunStatusSuccess,
			checkOutputFunc: func(t *testing.T, output *entity.EvaluatorOutputData, expected *entity.EvaluatorOutputData) {
				assert.NotNil(t, output.EvaluatorResult)
				assert.Equal(t, expected.EvaluatorResult.Score, output.EvaluatorResult.Score)
				assert.Equal(t, expected.EvaluatorResult.Reasoning, output.EvaluatorResult.Reasoning)
				assert.NotNil(t, output.EvaluatorUsage)
				assert.Equal(t, expected.EvaluatorUsage.InputTokens, output.EvaluatorUsage.InputTokens)
				assert.Equal(t, expected.EvaluatorUsage.OutputTokens, output.EvaluatorUsage.OutputTokens)
				assert.Nil(t, output.EvaluatorRunError)
				assert.GreaterOrEqual(t, output.TimeConsumingMS, int64(0))
			},
		},
		{
			name:      "LLM调用失败",
			evaluator: baseMockEvaluator,
			input:     baseMockInput,
			setupMocks: func() {
				expectedLlmError := errors.New("llm call failed")
				sharedMockLLMProvider.EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil, expectedLlmError)
				sharedMockMetric.EXPECT().EmitRun(int64(1), expectedLlmError, gomock.Any(), gomock.Any())
			},
			expectedOutput: &entity.EvaluatorOutputData{
				EvaluatorRunError: &entity.EvaluatorRunError{Message: "llm call failed"},
				EvaluatorResult:   nil,
				EvaluatorUsage:    &entity.EvaluatorUsage{},
			},
			expectedStatus: entity.EvaluatorRunStatusFail,
			checkOutputFunc: func(t *testing.T, output *entity.EvaluatorOutputData, expected *entity.EvaluatorOutputData) {
				assert.NotNil(t, output.EvaluatorRunError)
				assert.Contains(t, output.EvaluatorRunError.Message, expected.EvaluatorRunError.Message)
				assert.Nil(t, output.EvaluatorResult)
				assert.GreaterOrEqual(t, output.TimeConsumingMS, int64(0))
			},
		},
		{
			name:      "LLM返回ToolCalls为空",
			evaluator: baseMockEvaluator,
			input:     baseMockInput,
			setupMocks: func() {
				sharedMockLLMProvider.EXPECT().Call(gomock.Any(), gomock.Any()).Return(
					&entity.ReplyItem{
						ToolCalls: nil,
					}, nil)
				sharedMockMetric.EXPECT().EmitRun(int64(1), gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedOutput: &entity.EvaluatorOutputData{
				EvaluatorRunError: &entity.EvaluatorRunError{Message: "no tool calls returned from LLM"},
				EvaluatorResult:   nil,
				EvaluatorUsage:    &entity.EvaluatorUsage{InputTokens: 5, OutputTokens: 5},
			},
			expectedStatus: entity.EvaluatorRunStatusFail,
			checkOutputFunc: func(t *testing.T, output *entity.EvaluatorOutputData, expected *entity.EvaluatorOutputData) {
				assert.NotNil(t, output.EvaluatorRunError)
				assert.Nil(t, output.EvaluatorResult)
			},
		},
		{
			name:      "LLM返回FunctionCall Arguments 字段为空",
			evaluator: baseMockEvaluator,
			input:     baseMockInput,
			setupMocks: func() {
				sharedMockLLMProvider.EXPECT().Call(gomock.Any(), gomock.Any()).Return(
					&entity.ReplyItem{
						ToolCalls: []*entity.ToolCall{{Type: entity.ToolTypeFunction, FunctionCall: &entity.FunctionCall{
							Name:      "test_function",
							Arguments: ptr.Of(""),
						}}},
						TokenUsage: &entity.TokenUsage{InputTokens: 8, OutputTokens: 8},
					}, nil)
				sharedMockMetric.EXPECT().EmitRun(int64(1), gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedOutput: &entity.EvaluatorOutputData{
				EvaluatorRunError: &entity.EvaluatorRunError{Message: "function call arguments are nil"},
				EvaluatorResult:   nil,
				EvaluatorUsage:    &entity.EvaluatorUsage{InputTokens: 8, OutputTokens: 8},
			},
			expectedStatus: entity.EvaluatorRunStatusFail,
			checkOutputFunc: func(t *testing.T, output *entity.EvaluatorOutputData, expected *entity.EvaluatorOutputData) {
				assert.NotNil(t, output.EvaluatorRunError)
				assert.Nil(t, output.EvaluatorResult)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset expectations on shared mocks before each test case if necessary,
			// though gomock's Finish() at the end of TestEvaluatorSourcePromptServiceImpl_Run_Revised
			// should handle verification. New expectations are set in setupMocks.
			// If a mock method could be called by multiple test cases and needs specific
			// behavior per case, ensure .Times(1) or similar is used, or that expectations
			// are cleared/reset if the mocking framework supports it mid-test-function.
			// Gomock generally expects all defined EXPECT() calls to be met by the time ctrl.Finish() is called.
			// For table-driven tests sharing mocks, it's common to define all expectations for a given
			// mock call within the setupMocks of the specific test case that triggers it.

			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			output, status, _ := service.Run(ctx, tc.evaluator, tc.input)

			assert.Equal(t, tc.expectedStatus, status)
			if tc.checkOutputFunc != nil {
				tc.checkOutputFunc(t, output, tc.expectedOutput)
			}
			// Ensure all expectations set in setupMocks were met for this specific case.
			// This is tricky with shared mocks and a single ctrl.Finish().
			// A common pattern is one controller per sub-test (t.Run), but that
			// conflicts with a true singleton service that captures mocks at its creation.
			// The current setup relies on careful definition of EXPECT calls in each setupMocks
			// and that they don't unintentionally satisfy other test cases' calls.
		})
	}
}

// TestEvaluatorSourcePromptServiceImpl_PreHandle 测试 PreHandle 方法
func TestEvaluatorSourcePromptServiceImpl_PreHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLLMProvider := rpcmocks.NewMockILLMProvider(ctrl)
	mockMetric := metricsmocks.NewMockEvaluatorExecMetrics(ctrl)
	mockConfiger := configmocks.NewMockIConfiger(ctrl)

	service := &EvaluatorSourcePromptServiceImpl{
		llmProvider: mockLLMProvider,
		metric:      mockMetric,
		configer:    mockConfiger,
	}
	ctx := context.Background()

	testCases := []struct {
		name        string
		evaluator   *entity.Evaluator
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "成功预处理评估器",
			evaluator: &entity.Evaluator{
				ID:            100,
				SpaceID:       1,
				Name:          "Test Evaluator",
				EvaluatorType: entity.EvaluatorTypePrompt,
				PromptEvaluatorVersion: &entity.PromptEvaluatorVersion{
					ID:                100,
					EvaluatorID:       100,
					SpaceID:           1,
					PromptTemplateKey: "test-template-key",
					PromptSuffix:      "test-prompt-suffix",
					ModelConfig: &entity.ModelConfig{
						ModelID: 1,
					},
					ParseType: entity.ParseTypeFunctionCall,
				},
			},
			setupMocks: func() {
				// 如果需要设置 mock 期望
				mockConfiger.EXPECT().GetEvaluatorPromptSuffix(gomock.Any()).Return(map[string]string{
					"test-template-key": "test-prompt-suffix",
				}).Times(1)
				mockConfiger.EXPECT().GetEvaluatorToolConf(gomock.Any()).Return(map[string]*evaluator.Tool{
					"test_function": {
						Type: evaluator.ToolType(entity.ToolTypeFunction),
						Function: &evaluator.Function{
							Name:        "test_function",
							Description: ptr.Of("test description"),
							Parameters:  ptr.Of("{\"type\": \"object\", \"properties\": {\"score\": {\"type\": \"number\"}, \"reasoning\": {\"type\": \"string\"}}}"),
						},
					},
				}).Times(2)
				mockConfiger.EXPECT().GetEvaluatorToolMapping(gomock.Any()).Return(map[string]string{
					"test-template-key": "test-function",
				}).Times(1)
				mockConfiger.EXPECT().GetEvaluatorPromptSuffixMapping(gomock.Any()).Return(map[string]string{
					"1": "test-prompt-suffix",
				}).Times(1)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			err := service.PreHandle(ctx, tc.evaluator)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewEvaluatorSourcePromptServiceImpl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLLMProvider := rpcmocks.NewMockILLMProvider(ctrl)
	mockMetric := metricsmocks.NewMockEvaluatorExecMetrics(ctrl)
	mockConfiger := configmocks.NewMockIConfiger(ctrl)

	service := NewEvaluatorSourcePromptServiceImpl(
		mockLLMProvider,
		mockMetric,
		mockConfiger,
	)
	assert.NotNil(t, service)
	assert.Implements(t, (*EvaluatorSourceService)(nil), service)
}

func TestEvaluatorSourcePromptServiceImpl_Debug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLLMProvider := rpcmocks.NewMockILLMProvider(ctrl)
	mockMetric := metricsmocks.NewMockEvaluatorExecMetrics(ctrl)
	mockConfiger := configmocks.NewMockIConfiger(ctrl)

	service := &EvaluatorSourcePromptServiceImpl{
		llmProvider: mockLLMProvider,
		metric:      mockMetric,
		configer:    mockConfiger,
	}
	ctx := context.Background()

	baseMockEvaluator := &entity.Evaluator{
		ID:            100,
		SpaceID:       1,
		Name:          "Test Evaluator",
		EvaluatorType: entity.EvaluatorTypePrompt,
		PromptEvaluatorVersion: &entity.PromptEvaluatorVersion{
			ID:                100,
			EvaluatorID:       100,
			SpaceID:           1,
			PromptTemplateKey: "test-template-key",
			PromptSuffix:      "test-prompt-suffix",
			ModelConfig: &entity.ModelConfig{
				ModelID: 1,
			},
			ParseType: entity.ParseTypeFunctionCall,
			MessageList: []*entity.Message{
				{
					Role: entity.RoleSystem,
					Content: &entity.Content{
						ContentType: ptr.Of(entity.ContentTypeText),
						Text:        ptr.Of("{{test-content}}"),
					},
				},
			},
			Tools: []*entity.Tool{
				{
					Type: entity.ToolTypeFunction,
					Function: &entity.Function{
						Name:        "test_function",
						Description: "test description",
						Parameters:  "{\"type\": \"object\", \"properties\": {\"score\": {\"type\": \"number\"}, \"reasoning\": {\"type\": \"string\"}}}",
					},
				},
			},
		},
	}

	baseMockInput := &entity.EvaluatorInputData{
		InputFields: map[string]*entity.Content{
			"input": {
				ContentType: ptr.Of(entity.ContentTypeText),
				Text:        ptr.Of("test input"),
			},
		},
	}

	t.Run("成功调试评估器", func(t *testing.T) {
		mockLLMProvider.EXPECT().Call(gomock.Any(), gomock.Any()).Return(
			&entity.ReplyItem{
				ToolCalls: []*entity.ToolCall{
					{
						Type: entity.ToolTypeFunction,
						FunctionCall: &entity.FunctionCall{
							Name:      "test_function",
							Arguments: ptr.Of("{\"score\": 1.0, \"reason\": \"test response\"}"),
						},
					},
				},
				TokenUsage: &entity.TokenUsage{InputTokens: 10, OutputTokens: 10},
			}, nil)
		mockMetric.EXPECT().EmitRun(int64(1), gomock.Any(), gomock.Any(), gomock.Any())
		output, err := service.Debug(ctx, baseMockEvaluator, baseMockInput)
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotNil(t, output.EvaluatorResult)
		assert.Equal(t, 1.0, *output.EvaluatorResult.Score)
		assert.Equal(t, "test response", output.EvaluatorResult.Reasoning)
	})

	t.Run("调试评估器失败", func(t *testing.T) {
		mockLLMProvider.EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil, errors.New("llm call failed"))
		mockMetric.EXPECT().EmitRun(int64(1), gomock.Any(), gomock.Any(), gomock.Any())
		output, err := service.Debug(ctx, baseMockEvaluator, baseMockInput)
		assert.Error(t, err)
		assert.Nil(t, output)
	})
}

func TestEvaluatorSourcePromptServiceImpl_EvaluatorType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLLMProvider := rpcmocks.NewMockILLMProvider(ctrl)
	mockMetric := metricsmocks.NewMockEvaluatorExecMetrics(ctrl)
	mockConfiger := configmocks.NewMockIConfiger(ctrl)
	service := &EvaluatorSourcePromptServiceImpl{
		llmProvider: mockLLMProvider,
		metric:      mockMetric,
		configer:    mockConfiger,
	}
	assert.Equal(t, entity.EvaluatorTypePrompt, service.EvaluatorType())
}

func TestJSONRepair(t *testing.T) {
	t.Run("非法JSON应能修复", func(t *testing.T) {
		json := "{name: 'John'}"
		repaired, err := jsonrepair.JSONRepair(json)
		assert.NoError(t, err)
		assert.Equal(t, "{\"name\": \"John\"}", repaired)
	})

	t.Run("合法JSON应原样返回", func(t *testing.T) {
		json := "{\"name\":\"John\"}"
		repaired, err := jsonrepair.JSONRepair(json)
		assert.NoError(t, err)
		assert.Equal(t, json, repaired)
	})

	t.Run("完全不合法", func(t *testing.T) {
		json := "{name: John"
		referenceJson := "{\"name\": \"John\"}"

		repaired, err := jsonrepair.JSONRepair(json)
		assert.NoError(t, err)
		assert.Equal(t, referenceJson, repaired)
	})

	t.Run("空字符串应报错", func(t *testing.T) {
		json := ""
		repaired, err := jsonrepair.JSONRepair(json)
		assert.Error(t, err)
		assert.Empty(t, repaired)
	})

	t.Run("部分修复但仍不合法应报错", func(t *testing.T) {
		json := "{name: 'John', age: }"
		referenceJson := "{\"name\": \"John\", \"age\": null}"

		repaired, err := jsonrepair.JSONRepair(json)
		assert.NoError(t, err)
		assert.Equal(t, repaired, referenceJson)
	})

	t.Run("嵌套对象修复", func(t *testing.T) {
		json := "{user: {name: 'John', age: 30}}"
		repaired, err := jsonrepair.JSONRepair(json)
		assert.NoError(t, err)
		assert.Equal(t, "{\"user\": {\"name\": \"John\", \"age\": 30}}", repaired)
	})

	t.Run("数组修复", func(t *testing.T) {
		json := "[{name: 'John'}, {name: 'Jane'}]"
		repaired, err := jsonrepair.JSONRepair(json)
		assert.NoError(t, err)
		assert.Equal(t, "[{\"name\": \"John\"}, {\"name\": \"Jane\"}]", repaired)
	})
}

func TestParseOutput_ParseTypeContent(t *testing.T) {
	t.Run("ParseTypeContent-正常修复", func(t *testing.T) {
		evaluatorVersion := &entity.PromptEvaluatorVersion{
			ParseType: entity.ParseTypeContent,
			SpaceID:   1,
			Tools: []*entity.Tool{
				{
					Function: &entity.Function{
						Parameters: "{\"type\": \"object\", \"properties\": {\"score\": {\"type\": \"number\"}, \"reason\": {\"type\": \"string\"}}}",
					},
				},
			},
		}
		replyItem := &entity.ReplyItem{
			Content:    ptr.Of("{score: 1.5, reason: 'good'}"),
			TokenUsage: &entity.TokenUsage{InputTokens: 5, OutputTokens: 6},
		}
		output, err := parseOutput(context.Background(), evaluatorVersion, replyItem)
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotNil(t, output.EvaluatorResult)
		assert.Equal(t, 1.5, *output.EvaluatorResult.Score)
		assert.Equal(t, "good", output.EvaluatorResult.Reasoning)
		assert.Equal(t, int64(5), output.EvaluatorUsage.InputTokens)
		assert.Equal(t, int64(6), output.EvaluatorUsage.OutputTokens)
	})

	t.Run("ParseTypeContent-修复失败", func(t *testing.T) {
		evaluatorVersion := &entity.PromptEvaluatorVersion{
			ParseType: entity.ParseTypeContent,
			SpaceID:   1,
			Tools: []*entity.Tool{
				{
					Function: &entity.Function{
						Parameters: "{\"type\": \"object\", \"properties\": {\"score\": {\"type\": \"number\"}, \"reason\": {\"type\": \"string\"}}}",
					},
				},
			},
		}
		replyItem := &entity.ReplyItem{
			Content: ptr.Of("{score: 1.5, reason: }"), // reason缺失值
		}
		output, err := parseOutput(context.Background(), evaluatorVersion, replyItem)
		assert.Error(t, err)
		assert.NotNil(t, output)
	})

	t.Run("ParseTypeContent-字段类型错误", func(t *testing.T) {
		evaluatorVersion := &entity.PromptEvaluatorVersion{
			ParseType: entity.ParseTypeContent,
			SpaceID:   1,
			Tools: []*entity.Tool{
				{
					Function: &entity.Function{
						Parameters: "{\"type\": \"object\", \"properties\": {\"score\": {\"type\": \"number\"}, \"reason\": {\"type\": \"string\"}}}",
					},
				},
			},
		}
		replyItem := &entity.ReplyItem{
			Content: ptr.Of("{score: 'not-a-number', reason: 123}"),
		}
		output, err := parseOutput(context.Background(), evaluatorVersion, replyItem)
		assert.Error(t, err)
		assert.NotNil(t, output)
	})
}
