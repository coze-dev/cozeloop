// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package eino

import (
	"context"
	"time"

	ori_qianfan "github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/arkbot"
	"github.com/cloudwego/eino-ext/components/model/claude"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/model/qianfan"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	acl_openai "github.com/cloudwego/eino-ext/libs/acl/openai"
	einoModel "github.com/cloudwego/eino/components/model"
	"github.com/google/generative-ai-go/genai"
	"github.com/ollama/ollama/api"
	"github.com/pkg/errors"
	"google.golang.org/api/option"

	"github.com/coze-dev/coze-loop/backend/modules/llm/domain/entity"
	"github.com/coze-dev/coze-loop/backend/pkg/json"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
)

func NewLLM(ctx context.Context, model *entity.Model) (*LLM, error) {
	// 根据protocol导航到不同的builder
	var err error
	var chatModel einoModel.ToolCallingChatModel
	switch model.Protocol {
	case entity.ProtocolArk:
		chatModel, err = arkBuilder(ctx, model)
	case entity.ProtocolOpenAI:
		chatModel, err = openAIBuilder(ctx, model)
	case entity.ProtocolClaude:
		chatModel, err = claudeBuilder(ctx, model)
	case entity.ProtocolDeepseek:
		chatModel, err = deepSeekBuilder(ctx, model)
	case entity.ProtocolOllama:
		chatModel, err = ollamaBuilder(ctx, model)
	case entity.ProtocolGemini:
		chatModel, err = geminiBuilder(ctx, model)
	case entity.ProtocolQwen:
		chatModel, err = qwenBuilder(ctx, model)
	case entity.ProtocolQianfan:
		chatModel, err = qianfanBuilder(ctx, model)
	case entity.ProtocolArkBot:
		chatModel, err = arkBotBuilder(ctx, model)
	default:
		err = errors.Errorf("eino unsupport the protocol:%s", model.Protocol)
	}
	if err != nil {
		return nil, err
	}
	return &LLM{
		frame:     model.Frame,
		protocol:  model.Protocol,
		chatModel: chatModel,
	}, nil
}

func arkBuilder(ctx context.Context, model *entity.Model) (einoModel.ToolCallingChatModel, error) {
	if err := checkModelBeforeBuild(model); err != nil {
		return nil, err
	}
	p := model.ProtocolConfig
	cp := model.ParamConfig.GetCommonParamDefaultVal()
	cfg := &ark.ChatModelConfig{
		BaseURL:          p.BaseURL,
		APIKey:           p.APIKey,
		Model:            p.Model,
		MaxTokens:        cp.MaxTokens,
		Temperature:      cp.Temperature,
		TopP:             cp.TopP,
		Stop:             cp.Stop,
		FrequencyPenalty: cp.FrequencyPenalty,
		PresencePenalty:  cp.PresencePenalty,
	}
	if p.TimeoutMs != nil {
		cfg.Timeout = ptr.Of(time.Duration(*p.TimeoutMs) * time.Millisecond)
	}
	if arkCfg := p.ProtocolConfigArk; arkCfg != nil {
		cfg.Region = arkCfg.Region
		cfg.AccessKey = arkCfg.AccessKey
		cfg.SecretKey = arkCfg.SecretKey
		if arkCfg.RetryTimes != nil {
			cfg.RetryTimes = ptr.Of(int(*arkCfg.RetryTimes))
		}
		cfg.CustomHeader = arkCfg.CustomHeaders
	}
	return ark.NewChatModel(ctx, cfg)
}

func openAIBuilder(ctx context.Context, model *entity.Model) (einoModel.ToolCallingChatModel, error) {
	if err := checkModelBeforeBuild(model); err != nil {
		return nil, err
	}
	p := model.ProtocolConfig
	cp := model.ParamConfig.GetCommonParamDefaultVal()
	cfg := &openai.ChatModelConfig{
		APIKey:           p.APIKey,
		BaseURL:          p.BaseURL,
		Model:            p.Model,
		MaxTokens:        cp.MaxTokens,
		Temperature:      cp.Temperature,
		TopP:             cp.TopP,
		Stop:             cp.Stop,
		FrequencyPenalty: cp.FrequencyPenalty,
		PresencePenalty:  cp.PresencePenalty,
	}
	if p.TimeoutMs != nil {
		cfg.Timeout = time.Duration(*p.TimeoutMs) * time.Millisecond
	}
	if pc := p.ProtocolConfigOpenAI; pc != nil {
		cfg.ByAzure = pc.ByAzure
		cfg.APIVersion = pc.ApiVersion
		var js acl_openai.ChatCompletionResponseFormatJSONSchema
		if pc.ResponseFormatJsonSchema != "" {
			if err := sonic.UnmarshalString(pc.ResponseFormatJsonSchema, js); err != nil {
				return nil, err
			}
		}
		cfg.ResponseFormat = &acl_openai.ChatCompletionResponseFormat{
			Type:       acl_openai.ChatCompletionResponseFormatType(pc.ResponseFormatType),
			JSONSchema: &js,
		}
	}
	return openai.NewChatModel(ctx, cfg)
}

func claudeBuilder(ctx context.Context, model *entity.Model) (einoModel.ToolCallingChatModel, error) {
	if err := checkModelBeforeBuild(model); err != nil {
		return nil, err
	}
	p := model.ProtocolConfig
	cp := model.ParamConfig.GetCommonParamDefaultVal()
	cfg := &claude.Config{
		APIKey:        p.APIKey,
		Model:         p.Model,
		Temperature:   cp.Temperature,
		TopP:          cp.TopP,
		StopSequences: cp.Stop,
	}
	if p.BaseURL != "" {
		cfg.BaseURL = &p.BaseURL
	}
	if cp.MaxTokens != nil {
		cfg.MaxTokens = *cp.MaxTokens
	}
	if cp.TopK != nil {
		cfg.TopK = ptr.Of(int32(*cp.TopK))
	}
	if pc := p.ProtocolConfigClaude; pc != nil {
		cfg.ByBedrock = pc.ByBedrock
		cfg.AccessKey = pc.AccessKey
		cfg.SecretAccessKey = pc.SecretAccessKey
		cfg.SessionToken = pc.SessionToken
		cfg.Region = pc.Region
	}
	return claude.NewChatModel(ctx, cfg)
}

func deepSeekBuilder(ctx context.Context, model *entity.Model) (einoModel.ToolCallingChatModel, error) {
	if err := checkModelBeforeBuild(model); err != nil {
		return nil, err
	}
	p := model.ProtocolConfig
	cp := model.ParamConfig.GetCommonParamDefaultVal()
	cfg := &deepseek.ChatModelConfig{
		APIKey:  p.APIKey,
		BaseURL: p.BaseURL,
		Model:   p.Model,
		Stop:    cp.Stop,
	}
	if p.TimeoutMs != nil {
		cfg.Timeout = time.Duration(*p.TimeoutMs) * time.Millisecond
	}
	if cp.Temperature != nil {
		cfg.Temperature = *cp.Temperature
	}
	if cp.FrequencyPenalty != nil {
		cfg.FrequencyPenalty = *cp.FrequencyPenalty
	}
	if cp.PresencePenalty != nil {
		cfg.PresencePenalty = *cp.PresencePenalty
	}
	if cp.MaxTokens != nil {
		cfg.MaxTokens = *cp.MaxTokens
	}
	if cp.TopP != nil {
		cfg.TopP = *cp.TopP
	}
	if pc := p.ProtocolConfigDeepSeek; pc != nil {
		cfg.ResponseFormatType = deepseek.ResponseFormatType(pc.ResponseFormatType)
	}
	return deepseek.NewChatModel(ctx, cfg)
}

func checkModelBeforeBuild(model *entity.Model) error {
	if model == nil || model.ProtocolConfig == nil {
		return errors.Errorf("[checkModelBeforeBuild] failed as model:%s", json.MarshalStringIgnoreErr(model))
	}
	return nil
}

func geminiBuilder(ctx context.Context, model *entity.Model) (einoModel.ToolCallingChatModel, error) {
	if err := checkModelBeforeBuild(model); err != nil {
		return nil, err
	}
	p := model.ProtocolConfig
	cp := model.ParamConfig.GetCommonParamDefaultVal()
	cli, err := genai.NewClient(ctx, option.WithAPIKey(p.APIKey))
	if err != nil {
		return nil, err
	}
	cfg := &gemini.Config{
		Client:      cli,
		Model:       p.Model,
		MaxTokens:   cp.MaxTokens,
		Temperature: cp.Temperature,
		TopP:        cp.TopP,
	}
	if cp.TopK != nil {
		cfg.TopK = ptr.Of(int32(*cp.TopK))
	}
	if pc := p.ProtocolConfigGemini; pc != nil {
		if pc.ResponseSchema != nil && *pc.ResponseSchema != "" {
			if err := sonic.UnmarshalString(*pc.ResponseSchema, &cfg.ResponseSchema); err != nil {
				return nil, err
			}
		}
		cfg.EnableCodeExecution = pc.EnableCodeExecution
		for _, ss := range pc.SafetySettings {
			cfg.SafetySettings = append(cfg.SafetySettings, &genai.SafetySetting{
				Category:  genai.HarmCategory(ss.Category),
				Threshold: genai.HarmBlockThreshold(ss.Threshold),
			})
		}
	}
	return gemini.NewChatModel(ctx, cfg)
}

func ollamaBuilder(ctx context.Context, model *entity.Model) (einoModel.ToolCallingChatModel, error) {
	if err := checkModelBeforeBuild(model); err != nil {
		return nil, err
	}
	p := model.ProtocolConfig
	cp := model.ParamConfig.GetCommonParamDefaultVal()
	cfg := &ollama.ChatModelConfig{
		BaseURL:   p.BaseURL,
		Model:     p.Model,
		Format:    nil,
		KeepAlive: nil,
		Options: &api.Options{
			// NumKeep:          0,
			// Seed:             0,
			// NumPredict:       0,
			TopK: ptr.From(cp.TopK),
			TopP: ptr.From(cp.TopP),
			// MinP:             0,
			// TypicalP:         0,
			// RepeatLastN:      0,
			Temperature:      ptr.From(cp.Temperature),
			RepeatPenalty:    0,
			PresencePenalty:  ptr.From(cp.PresencePenalty),
			FrequencyPenalty: ptr.From(cp.FrequencyPenalty),
			Stop:             cp.Stop,
		},
	}
	if p.TimeoutMs != nil {
		cfg.Timeout = time.Duration(*p.TimeoutMs) * time.Millisecond
	}
	if pc := p.ProtocolConfigOllama; pc != nil {
		if pc.Format != nil && *pc.Format != "" {
			cfg.Format = []byte(*pc.Format)
		}
		if pc.KeepAliveMs != nil && *pc.KeepAliveMs > 0 {
			cfg.KeepAlive = ptr.Of(time.Duration(*pc.KeepAliveMs) * time.Millisecond)
		}
	}
	return ollama.NewChatModel(ctx, cfg)
}

func qwenBuilder(ctx context.Context, model *entity.Model) (einoModel.ToolCallingChatModel, error) {
	if err := checkModelBeforeBuild(model); err != nil {
		return nil, err
	}
	p := model.ProtocolConfig
	cp := model.ParamConfig.GetCommonParamDefaultVal()
	cfg := &qwen.ChatModelConfig{
		APIKey:           p.APIKey,
		BaseURL:          p.BaseURL,
		Model:            p.Model,
		MaxTokens:        cp.MaxTokens,
		Temperature:      cp.Temperature,
		TopP:             cp.TopP,
		Stop:             cp.Stop,
		PresencePenalty:  cp.PresencePenalty,
		FrequencyPenalty: cp.FrequencyPenalty,
	}
	if p.TimeoutMs != nil {
		cfg.Timeout = time.Duration(*p.TimeoutMs) * time.Millisecond
	}
	if pc := p.ProtocolConfigQwen; pc != nil {
		if pc.ResponseFormatType != nil && pc.ResponseFormatJsonSchema != nil {
			var js acl_openai.ChatCompletionResponseFormatJSONSchema
			if *pc.ResponseFormatJsonSchema != "" {
				if err := sonic.UnmarshalString(*pc.ResponseFormatJsonSchema, js); err != nil {
					return nil, err
				}
			}
			cfg.ResponseFormat = &acl_openai.ChatCompletionResponseFormat{
				Type:       acl_openai.ChatCompletionResponseFormatType(*pc.ResponseFormatType),
				JSONSchema: &js,
			}
		}
	}
	return qwen.NewChatModel(ctx, cfg)
}

func qianfanBuilder(ctx context.Context, model *entity.Model) (einoModel.ToolCallingChatModel, error) {
	if err := checkModelBeforeBuild(model); err != nil {
		return nil, err
	}
	p := model.ProtocolConfig
	cp := model.ParamConfig.GetCommonParamDefaultVal()
	cfg := &qianfan.ChatModelConfig{
		Model:               p.Model,
		Temperature:         cp.Temperature,
		TopP:                cp.TopP,
		MaxCompletionTokens: cp.MaxTokens,
		Stop:                cp.Stop,
	}
	if cp.FrequencyPenalty != nil {
		cfg.FrequencyPenalty = ptr.Of(float64(*cp.FrequencyPenalty))
	}
	if cp.PresencePenalty != nil {
		cfg.PresencePenalty = ptr.Of(float64(*cp.PresencePenalty))
	}
	if pc := p.ProtocolConfigQianfan; pc != nil {
		cfg.LLMRetryCount = pc.LLMRetryCount
		cfg.LLMRetryTimeout = pc.LLMRetryTimeout
		cfg.LLMRetryBackoffFactor = pc.LLMRetryBackoffFactor
		cfg.ParallelToolCalls = pc.ParallelToolCalls
		if pc.ResponseFormatType != nil && pc.ResponseFormatJsonSchema != nil {
			var js any
			if *pc.ResponseFormatJsonSchema != "" {
				if err := sonic.UnmarshalString(*pc.ResponseFormatJsonSchema, js); err != nil {
					return nil, err
				}
			}
			cfg.ResponseFormat = &ori_qianfan.ResponseFormat{
				FormatType: *pc.ResponseFormatType,
				JsonSchema: &js,
			}
		}

	}
	return qianfan.NewChatModel(ctx, cfg)
}

func arkBotBuilder(ctx context.Context, model *entity.Model) (einoModel.ToolCallingChatModel, error) {
	if err := checkModelBeforeBuild(model); err != nil {
		return nil, err
	}
	p := model.ProtocolConfig
	cp := model.ParamConfig.GetCommonParamDefaultVal()
	cfg := &arkbot.Config{
		BaseURL:          p.BaseURL,
		APIKey:           p.APIKey,
		Model:            p.Model,
		MaxTokens:        cp.MaxTokens,
		Temperature:      cp.Temperature,
		TopP:             cp.TopP,
		Stop:             cp.Stop,
		FrequencyPenalty: cp.FrequencyPenalty,
		PresencePenalty:  cp.PresencePenalty,
	}
	if p.TimeoutMs != nil {
		cfg.Timeout = ptr.Of(time.Duration(*p.TimeoutMs) * time.Millisecond)
	}
	if arkCfg := p.ProtocolConfigArkBot; arkCfg != nil {
		cfg.Region = arkCfg.Region
		cfg.AccessKey = arkCfg.AccessKey
		cfg.SecretKey = arkCfg.SecretKey
		if arkCfg.RetryTimes != nil {
			cfg.RetryTimes = ptr.Of(int(*arkCfg.RetryTimes))
		}
		cfg.CustomHeader = arkCfg.CustomHeaders
	}
	return arkbot.NewChatModel(ctx, cfg)
}
