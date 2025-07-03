// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/bytedance/gg/gmap"
	"github.com/bytedance/gg/gptr"
	"github.com/coze-dev/cozeloop-go"
	"github.com/kaptinlin/jsonrepair"
	"github.com/valyala/fasttemplate"

	"github.com/coze-dev/cozeloop/backend/infra/looptracer"
	"github.com/coze-dev/cozeloop/backend/infra/middleware/session"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/application/convertor/evaluator"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/consts"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/component/metrics"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/component/rpc"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/infra/tracer"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/pkg/conf"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/pkg/json"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

const (
	TemplateStartTag = "{{"
	TemplateEndTag   = "}}"
)

var (
	evaluatorVersionServiceOnce      = sync.Once{}
	singletonEvaluatorVersionService EvaluatorSourceService
)

func NewEvaluatorSourcePromptServiceImpl(
	llmProvider rpc.ILLMProvider,
	metric metrics.EvaluatorExecMetrics,
	configer conf.IConfiger,
) EvaluatorSourceService {
	evaluatorVersionServiceOnce.Do(func() {
		singletonEvaluatorVersionService = &EvaluatorSourcePromptServiceImpl{
			llmProvider: llmProvider,
			metric:      metric,
			configer:    configer,
		}
	})
	return singletonEvaluatorVersionService
}

type EvaluatorSourcePromptServiceImpl struct {
	llmProvider rpc.ILLMProvider
	metric      metrics.EvaluatorExecMetrics
	configer    conf.IConfiger
}

func (p *EvaluatorSourcePromptServiceImpl) EvaluatorType() entity.EvaluatorType {
	return entity.EvaluatorTypePrompt
}

func (p *EvaluatorSourcePromptServiceImpl) Run(ctx context.Context, evaluator *entity.Evaluator, input *entity.EvaluatorInputData) (output *entity.EvaluatorOutputData, runStatus entity.EvaluatorRunStatus, traceID string) {
	var err error
	startTime := time.Now()
	rootSpan, ctx := newEvaluatorSpan(ctx, evaluator.Name, "LoopEvaluation", strconv.FormatInt(evaluator.SpaceID, 10), false)
	traceID = rootSpan.GetTraceID()
	defer func() {
		if output == nil {
			output = &entity.EvaluatorOutputData{
				EvaluatorRunError: &entity.EvaluatorRunError{},
			}
		}
		var errInfo error
		if err != nil {
			if output.EvaluatorRunError == nil {
				output.EvaluatorRunError = &entity.EvaluatorRunError{}
			}
			statusErr, ok := errorx.FromStatusError(err)
			if ok {
				output.EvaluatorRunError.Code = statusErr.Code()
				output.EvaluatorRunError.Message = statusErr.Error()
				errInfo = statusErr
			} else {
				output.EvaluatorRunError.Code = errno.RunEvaluatorFailCode
				output.EvaluatorRunError.Message = err.Error()
				errInfo = err
			}
		}
		rootSpan.reportRootSpan(ctx, &ReportRootSpanRequest{
			input:            input,
			output:           output,
			runStatus:        runStatus,
			evaluatorVersion: evaluator.PromptEvaluatorVersion,
			errInfo:          errInfo,
		})
	}()

	err = evaluator.GetEvaluatorVersion().ValidateBaseInfo()
	if err != nil {
		logs.CtxInfo(ctx, "[RunEvaluator] ValidateBaseInfo fail, err: %v", err)
		runStatus = entity.EvaluatorRunStatusFail
		return nil, runStatus, traceID
	}
	// 校验输入数据
	err = evaluator.GetEvaluatorVersion().ValidateInput(input)
	if err != nil {
		logs.CtxInfo(ctx, "[RunEvaluator] ValidateInput fail, err: %v", err)
		runStatus = entity.EvaluatorRunStatusFail
		return nil, runStatus, traceID
	}
	defer func() {
		var modelID string
		if evaluator.PromptEvaluatorVersion.ModelConfig.ModelID == 0 {
			modelID = ptr.From(evaluator.PromptEvaluatorVersion.ModelConfig.ProviderModelID)
		} else {
			modelID = strconv.FormatInt(evaluator.PromptEvaluatorVersion.ModelConfig.ModelID, 10)
		}

		p.metric.EmitRun(evaluator.SpaceID, err, startTime, modelID)
	}()
	// 渲染变量
	err = renderTemplate(ctx, evaluator.PromptEvaluatorVersion, input)
	if err != nil {
		logs.CtxError(ctx, "[RunEvaluator] renderTemplate fail, err: %v", err)
		runStatus = entity.EvaluatorRunStatusFail
		return nil, runStatus, traceID
	}
	// 执行评估逻辑
	userIDInContext := session.UserIDInCtxOrEmpty(ctx)
	llmResp, err := p.chat(ctx, evaluator.PromptEvaluatorVersion, userIDInContext)
	if err != nil {
		logs.CtxError(ctx, "[RunEvaluator] chat fail, err: %v", err)
		runStatus = entity.EvaluatorRunStatusFail
		return nil, runStatus, traceID
	}
	output, err = parseOutput(ctx, evaluator.PromptEvaluatorVersion, llmResp)
	if err != nil {
		logs.CtxWarn(ctx, "[RunEvaluator] parseOutput fail, err: %v", err)
		runStatus = entity.EvaluatorRunStatusFail
		return nil, runStatus, traceID
	}
	return output, entity.EvaluatorRunStatusSuccess, traceID
}

func (p *EvaluatorSourcePromptServiceImpl) chat(ctx context.Context, evaluatorVersion *entity.PromptEvaluatorVersion, userIDInContext string) (resp *entity.ReplyItem, err error) {
	modelSpan, modelCtx := newEvaluatorSpan(ctx, evaluatorVersion.ModelConfig.ModelName, "model", strconv.FormatInt(evaluatorVersion.SpaceID, 10), true)
	defer func() {
		modelSpan.reportModelSpan(modelCtx, evaluatorVersion, resp, err)
	}()
	modelTraceCtx := looptracer.GetTracer().Inject(modelCtx)
	if err != nil {
		logs.CtxWarn(ctx, "[RunEvaluator] Inject fail, err: %v", err)
	}

	llmCallParam := &entity.LLMCallParam{
		SpaceID:     evaluatorVersion.GetSpaceID(),
		EvaluatorID: strconv.FormatInt(evaluatorVersion.EvaluatorID, 10),
		UserID:      gptr.Of(userIDInContext),
		Scenario:    entity.ScenarioEvaluator,
		Messages:    evaluatorVersion.MessageList,
		ModelConfig: evaluatorVersion.ModelConfig,
	}
	if evaluatorVersion.ParseType == entity.ParseTypeFunctionCall {
		llmCallParam.Tools = evaluatorVersion.Tools
		llmCallParam.ToolCallConfig = &entity.ToolCallConfig{
			ToolChoice: entity.ToolChoiceTypeRequired,
		}
	}
	resp, err = p.llmProvider.Call(modelTraceCtx, llmCallParam)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type evaluatorSpan struct {
	looptracer.Span
}

func newEvaluatorSpan(ctx context.Context, spanName, spanType, spaceID string, asyncChild bool) (*evaluatorSpan, context.Context) {
	var evalSpan looptracer.Span
	var nctx context.Context
	if asyncChild {
		nctx, evalSpan = looptracer.GetTracer().StartSpan(ctx, spanName, spanType, cozeloop.WithSpanWorkspaceID(spaceID))
	} else {
		nctx, evalSpan = looptracer.GetTracer().StartSpan(ctx, spanName, spanType, cozeloop.WithStartNewTrace(), cozeloop.WithSpanWorkspaceID(spaceID))
	}

	return &evaluatorSpan{
		Span: evalSpan,
	}, nctx
}

type toolCallSpanContent struct {
	ToolCall *entity.ToolCall `json:"tool_call"`
}

type ReportRootSpanRequest struct {
	input            *entity.EvaluatorInputData
	output           *entity.EvaluatorOutputData
	runStatus        entity.EvaluatorRunStatus
	evaluatorVersion *entity.PromptEvaluatorVersion
	errInfo          error
}

func (e *evaluatorSpan) reportRootSpan(ctx context.Context, reportRootSpanRequest *ReportRootSpanRequest) {
	e.SetInput(ctx, tracer.Convert2TraceString(reportRootSpanRequest.input))
	if reportRootSpanRequest.output != nil {
		e.SetOutput(ctx, tracer.Convert2TraceString(reportRootSpanRequest.output.EvaluatorResult))
	}
	switch reportRootSpanRequest.runStatus {
	case entity.EvaluatorRunStatusSuccess:
		e.SetStatusCode(ctx, 0)
	case entity.EvaluatorRunStatusFail:
		e.SetStatusCode(ctx, int(entity.EvaluatorRunStatusFail))
		e.SetError(ctx, reportRootSpanRequest.errInfo)
	default:
		e.SetStatusCode(ctx, 0) // 默认为成功
	}
	tags := make(map[string]interface{}, 0)
	tags["evaluator_id"] = reportRootSpanRequest.evaluatorVersion.EvaluatorID
	tags["evaluator_version"] = reportRootSpanRequest.evaluatorVersion.Version
	e.SetCallType("Evaluator")
	userIDInContext := session.UserIDInCtxOrEmpty(ctx)
	if userIDInContext != "" {
		e.SetUserID(ctx, userIDInContext)
	}
	e.SetTags(ctx, tags)
	e.Finish(ctx)
}

func (e *evaluatorSpan) reportModelSpan(ctx context.Context, evaluatorVersion *entity.PromptEvaluatorVersion, replyItem *entity.ReplyItem, respErr error) {
	if respErr != nil {
		e.SetStatusCode(ctx, errno.InvalidOutputFromModelCode)
		e.SetError(ctx, respErr)
	}
	if evaluatorVersion.ParseType == entity.ParseTypeFunctionCall {
		if replyItem != nil && len(replyItem.ToolCalls) > 0 {
			e.SetOutput(ctx, tracer.Convert2TraceString(&toolCallSpanContent{
				ToolCall: replyItem.ToolCalls[0],
			}))
			if replyItem.TokenUsage != nil {
				e.SetInputTokens(ctx, int(replyItem.TokenUsage.InputTokens))
				e.SetOutputTokens(ctx, int(replyItem.TokenUsage.OutputTokens))
			}
		} else {
			e.SetStatusCode(ctx, errno.InvalidOutputFromModelCode)
			e.SetError(ctx, errorx.New("LLM response empty"))
		}
	} else {
		if replyItem != nil {
			e.SetOutput(ctx, replyItem.Content)
			if replyItem.TokenUsage != nil {
				e.SetInputTokens(ctx, int(replyItem.TokenUsage.InputTokens))
				e.SetOutputTokens(ctx, int(replyItem.TokenUsage.OutputTokens))
			}
		} else {
			e.SetStatusCode(ctx, errno.InvalidOutputFromModelCode)
			e.SetError(ctx, errorx.New("LLM response empty"))
		}
	}
	e.SetCallType("Evaluator")
	userIDInContext := session.UserIDInCtxOrEmpty(ctx)
	if userIDInContext != "" {
		e.SetUserID(ctx, userIDInContext)
	}
	tags := tracer.ConvertModel2Ob(evaluatorVersion.MessageList, evaluatorVersion.Tools)
	tags["model_config"] = tracer.Convert2TraceString(evaluatorVersion.ModelConfig)
	e.SetTags(ctx, tags)
	e.Finish(ctx)
}

func (e *evaluatorSpan) reportOutputParserSpan(ctx context.Context, replyItem *entity.ReplyItem, output *entity.EvaluatorOutputData, spaceID string, errInfo error) {
	if replyItem != nil && len(replyItem.ToolCalls) > 0 {
		e.SetInput(ctx, tracer.Convert2TraceString(&toolCallSpanContent{
			ToolCall: replyItem.ToolCalls[0],
		}))
	}
	if output != nil {
		e.SetOutput(ctx, tracer.Convert2TraceString(output.EvaluatorResult))
	}
	if errInfo != nil {
		e.SetStatusCode(ctx, int(entity.EvaluatorRunStatusFail))
		e.SetError(ctx, errInfo)
	} else {
		e.SetStatusCode(ctx, 0)
	}
	tags := make(map[string]interface{})
	e.SetCallType("Evaluator")
	userIDInContext := session.UserIDInCtxOrEmpty(ctx)
	if userIDInContext != "" {
		e.SetUserID(ctx, userIDInContext)
	}
	e.SetTags(ctx, tags)
	e.Finish(ctx)
}

func parseOutput(ctx context.Context, evaluatorVersion *entity.PromptEvaluatorVersion, replyItem *entity.ReplyItem) (output *entity.EvaluatorOutputData, err error) {
	// 输出数据全空直接返回
	outputParserSpan, ctx := newEvaluatorSpan(ctx, "ParseOutput", "LoopEvaluation", strconv.FormatInt(evaluatorVersion.SpaceID, 10), true)
	defer func() {
		outputParserSpan.reportOutputParserSpan(ctx, replyItem, output, strconv.FormatInt(evaluatorVersion.SpaceID, 10), err)
	}()
	output = &entity.EvaluatorOutputData{
		EvaluatorResult: &entity.EvaluatorResult{},
		EvaluatorUsage:  &entity.EvaluatorUsage{},
	}
	if replyItem == nil {
		logs.CtxWarn(ctx, "[RunEvaluator] parseOutput fail, err: resp is nil")
		return output, errorx.NewByCode(errno.LLMOutputEmptyCode, errorx.WithExtraMsg(" resp is nil"))
	}
	var repairArgs string
	if evaluatorVersion.ParseType == entity.ParseTypeContent {
		repairArgs, err = jsonrepair.JSONRepair(gptr.Indirect(replyItem.Content))
		if err != nil {
			logs.CtxWarn(ctx, "[RunEvaluator] parseOutput Content RepairJSON fail, origin content: %v, err: %v", gptr.Indirect(replyItem.Content), err)
			return output, errorx.NewByCode(errno.InvalidOutputFromModelCode)
		}
	} else {
		if len(replyItem.ToolCalls) == 0 {
			logs.CtxWarn(ctx, "[RunEvaluator] parseOutput fail, err: tool call empty")
			return output, errorx.NewByCode(errno.LLMToolCallFailCode)
		}
		repairArgs, err = jsonrepair.JSONRepair(gptr.Indirect(replyItem.ToolCalls[0].FunctionCall.Arguments))
		if err != nil {
			logs.CtxWarn(ctx, "[RunEvaluator] parseOutput ToolCalls RepairJSON fail, origin content: %v, err: %v", gptr.Indirect(replyItem.ToolCalls[0].FunctionCall.Arguments), err)
			return output, errorx.NewByCode(errno.InvalidOutputFromModelCode)
		}
	}
	// 解析输出数据
	params := evaluatorVersion.Tools[0].Function.Parameters
	var scoreFieldValue any
	scoreFieldValue, err = json.ExtractFieldValue(params, repairArgs, "score")
	if err != nil {
		logs.CtxWarn(ctx, "[RunEvaluator] parseOutput ExtractFieldValue score fail, repairArgs: %v, err: %v", repairArgs, err)
		err = errorx.NewByCode(errno.InvalidOutputFromModelCode)
	}
	if score, ok := scoreFieldValue.(float64); ok {
		output.EvaluatorResult.Score = &score
	} else {
		logs.CtxWarn(ctx, "[RunEvaluator] parseOutput fail, repairArgs: %v, err: score not float64", repairArgs)
		err = errorx.NewByCode(errno.InvalidOutputFromModelCode)
	}
	var reasonFieldValue any
	reasonFieldValue, err = json.ExtractFieldValue(params, repairArgs, "reason")
	if err != nil {
		logs.CtxWarn(ctx, "[RunEvaluator] parseOutput ReasonFieldValue reason fail, repairArgs: %v, err: %v", repairArgs, err)
		err = errorx.NewByCode(errno.InvalidOutputFromModelCode)
	}
	if reason, ok := reasonFieldValue.(string); ok {
		output.EvaluatorResult.Reasoning = reason
	} else {
		logs.CtxWarn(ctx, "[RunEvaluator] parseOutput fail, repairArgs: %v, err: reason not string", repairArgs)
		err = errorx.NewByCode(errno.InvalidOutputFromModelCode)
	}
	if replyItem.TokenUsage != nil {
		output.EvaluatorUsage.InputTokens = replyItem.TokenUsage.InputTokens
		output.EvaluatorUsage.OutputTokens = replyItem.TokenUsage.OutputTokens
	}

	return output, err
}

func renderTemplate(ctx context.Context, evaluatorVersion *entity.PromptEvaluatorVersion, input *entity.EvaluatorInputData) error {
	// 实现渲染模板的逻辑
	renderTemplateSpan, ctx := newEvaluatorSpan(ctx, "RenderTemplate", "prompt", strconv.FormatInt(evaluatorVersion.SpaceID, 10), true)
	renderTemplateSpan.SetInput(ctx, tracer.Convert2TraceString(tracer.ConvertPrompt2Ob(evaluatorVersion.MessageList,
		gmap.Map(input.InputFields, func(key string, value *entity.Content) (string, any) {
			if value == nil {
				return key, nil
			}
			return key, value.Text
		}))))
	for _, message := range evaluatorVersion.MessageList {
		// 现阶段只支持text类型模板渲染
		if gptr.Indirect(message.Content.ContentType) == entity.ContentTypeText {
			message.Content.Text = gptr.Of(fasttemplate.ExecuteFuncString(gptr.Indirect(message.Content.Text), TemplateStartTag, TemplateEndTag, func(w io.Writer, tag string) (int, error) {
				// 输入变量里没有就不做替换直接返回
				if v, ok := input.InputFields[tag]; !ok || v == nil {
					return w.Write([]byte(""))
				}
				// 目前仅适用text替换
				return w.Write([]byte(gptr.Indirect(input.InputFields[tag].Text)))
			}))
		}
	}
	if len(evaluatorVersion.MessageList) > 0 {
		evaluatorVersion.MessageList[0].Content.Text = gptr.Of(gptr.Indirect(evaluatorVersion.MessageList[0].Content.Text) + evaluatorVersion.PromptSuffix)
	}

	renderTemplateSpan.SetOutput(ctx, tracer.Convert2TraceString(tracer.ConvertPrompt2Ob(evaluatorVersion.MessageList, nil)))
	tags := make(map[string]interface{})
	renderTemplateSpan.SetTags(ctx, tags)
	renderTemplateSpan.SetCallType("Evaluator")
	userIDInContext := session.UserIDInCtxOrEmpty(ctx)
	if userIDInContext != "" {
		renderTemplateSpan.SetUserID(ctx, userIDInContext)
	}
	renderTemplateSpan.Finish(ctx)
	return nil
}

func (p *EvaluatorSourcePromptServiceImpl) Debug(ctx context.Context, evaluator *entity.Evaluator, input *entity.EvaluatorInputData) (output *entity.EvaluatorOutputData, err error) {
	// 实现调试评估的逻辑
	output, _, _ = p.Run(ctx, evaluator, input)
	if output != nil && output.EvaluatorRunError != nil {
		return nil, errorx.NewByCode(output.EvaluatorRunError.Code, errorx.WithExtraMsg(output.EvaluatorRunError.Message))
	}
	return output, nil
}

func (p *EvaluatorSourcePromptServiceImpl) PreHandle(ctx context.Context, evaluator *entity.Evaluator) error {
	p.injectPromptTools(ctx, evaluator)
	p.injectParseType(ctx, evaluator)
	return nil
}

func (p *EvaluatorSourcePromptServiceImpl) injectPromptTools(ctx context.Context, evaluatorDO *entity.Evaluator) {
	// 注入默认工具
	tools := make([]*entity.Tool, 0, len(p.configer.GetEvaluatorToolConf(ctx)))

	if toolKey, ok := p.configer.GetEvaluatorToolMapping(ctx)[evaluatorDO.GetEvaluatorVersion().GetPromptTemplateKey()]; ok {
		tools = append(tools, evaluator.ConvertToolDTO2DO(p.configer.GetEvaluatorToolConf(ctx)[toolKey]))
	} else {
		tools = append(tools, evaluator.ConvertToolDTO2DO(p.configer.GetEvaluatorToolConf(ctx)[consts.DefaultEvaluatorToolKey]))
	}
	evaluatorDO.GetEvaluatorVersion().SetTools(tools)
}

func (p *EvaluatorSourcePromptServiceImpl) injectParseType(ctx context.Context, evaluatorDO *entity.Evaluator) {
	// 注入后缀
	if suffixKey, ok := p.configer.GetEvaluatorPromptSuffixMapping(ctx)[strconv.FormatInt(evaluatorDO.GetEvaluatorVersion().GetModelConfig().ModelID, 10)]; ok {
		evaluatorDO.GetEvaluatorVersion().SetPromptSuffix(p.configer.GetEvaluatorPromptSuffix(ctx)[suffixKey])
		evaluatorDO.GetEvaluatorVersion().SetParseType(entity.ParseType(suffixKey))
	} else {
		evaluatorDO.GetEvaluatorVersion().SetPromptSuffix(p.configer.GetEvaluatorPromptSuffix(ctx)[consts.DefaultEvaluatorPromptSuffixKey])
		evaluatorDO.GetEvaluatorVersion().SetParseType(entity.ParseTypeFunctionCall)
	}
}
