// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/component/config"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/component/metrics"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/component/mq"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/loop_span"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/repo"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/service/trace/span_filter"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/service/trace/span_processor"
	obErrorx "github.com/coze-dev/cozeloop/backend/modules/observability/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/goroutine"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
	time_util "github.com/coze-dev/cozeloop/backend/pkg/time"
)

type ListSpansReq struct {
	WorkspaceID     int64
	StartTime       int64 // ms
	EndTime         int64 // ms
	Filters         *loop_span.FilterFields
	Limit           int32
	DescByStartTime bool
	PageToken       string
	PlatformType    loop_span.PlatformType
	SpanListType    loop_span.SpanListType
}

type ListSpansResp struct {
	Spans         loop_span.SpanList
	NextPageToken string
	HasMore       bool
}

type GetTraceReq struct {
	WorkspaceID  int64
	TraceID      string
	StartTime    int64 // ms
	EndTime      int64 // ms
	PlatformType loop_span.PlatformType
}

type GetTraceResp struct {
	TraceId string
	Spans   loop_span.SpanList
}

type TraceQueryParam struct {
	TraceID   string
	StartTime int64 // ms
	EndTime   int64 // ms
}

type GetTracesAdvanceInfoReq struct {
	WorkspaceID  int64
	Traces       []*TraceQueryParam
	PlatformType loop_span.PlatformType
}

type GetTracesAdvanceInfoResp struct {
	Infos []*loop_span.TraceAdvanceInfo
}

type IngestTracesReq struct {
	TTL              entity.TTL
	WhichIsEnough    int
	CozeAccountId    string
	VolcanoAccountID int64
	Spans            loop_span.SpanList
}

type SendTraceResp struct{}

type GetTracesMetaInfoReq struct {
	WorkspaceID  int64
	PlatformType loop_span.PlatformType
	SpanListType loop_span.SpanListType
}

type GetTracesMetaInfoResp struct {
	FilesMetas map[string]*config.FieldMeta
}

//go:generate mockgen -destination=mocks/trace_service.go -package=mocks . ITraceService
type ITraceService interface {
	ListSpans(ctx context.Context, req *ListSpansReq) (*ListSpansResp, error)
	GetTrace(ctx context.Context, req *GetTraceReq) (*GetTraceResp, error)
	GetTracesAdvanceInfo(ctx context.Context, req *GetTracesAdvanceInfoReq) (*GetTracesAdvanceInfoResp, error)
	IngestTraces(ctx context.Context, req *IngestTracesReq) error
	GetTracesMetaInfo(ctx context.Context, req *GetTracesMetaInfoReq) (*GetTracesMetaInfoResp, error)
}

func NewTraceServiceImpl(
	tRepo repo.ITraceRepo,
	traceConfig config.ITraceConfig,
	traceProducer mq.ITraceProducer,
	metrics metrics.ITraceMetrics,
	buildHelper TraceFilterProcessorBuilder,
) (ITraceService, error) {
	return &TraceServiceImpl{
		traceRepo:     tRepo,
		traceConfig:   traceConfig,
		traceProducer: traceProducer,
		buildHelper:   buildHelper,
		metrics:       metrics,
	}, nil
}

type TraceServiceImpl struct {
	traceRepo     repo.ITraceRepo
	traceConfig   config.ITraceConfig
	traceProducer mq.ITraceProducer
	metrics       metrics.ITraceMetrics
	buildHelper   TraceFilterProcessorBuilder
}

func (r *TraceServiceImpl) GetTrace(ctx context.Context, req *GetTraceReq) (*GetTraceResp, error) {
	tenants, err := r.getTenants(ctx, req.PlatformType)
	if err != nil {
		return nil, err
	}
	st := time.Now()
	spans, err := r.traceRepo.GetTrace(ctx, &repo.GetTraceParam{
		Tenants: tenants,
		TraceID: req.TraceID,
		StartAt: req.StartTime,
		EndAt:   req.EndTime,
		Limit:   1000,
	})
	r.metrics.EmitGetTrace(req.WorkspaceID, st, err != nil)
	if err != nil {
		return nil, err
	}
	processors, err := r.buildHelper.BuildGetTraceProcessors(ctx, span_processor.Settings{
		WorkspaceId:    req.WorkspaceID,
		PlatformType:   req.PlatformType,
		QueryStartTime: req.StartTime,
		QueryEndTime:   req.EndTime,
	})
	if err != nil {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInternalErrorCodeCode)
	}
	for _, p := range processors {
		spans, err = p.Transform(ctx, spans)
		if err != nil {
			return nil, err
		}
	}
	spans.SortByStartTime(false)
	return &GetTraceResp{
		TraceId: req.TraceID,
		Spans:   spans,
	}, nil
}

func (r *TraceServiceImpl) ListSpans(ctx context.Context, req *ListSpansReq) (*ListSpansResp, error) {
	if err := req.Filters.Traverse(processSpecificFilter); err != nil {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid filter"))
	}
	platformFilter, err := r.buildHelper.BuildPlatformRelatedFilter(ctx, req.PlatformType)
	if err != nil {
		return nil, err
	}
	builtinFilter, err := r.buildBuiltinFilters(ctx, platformFilter, req)
	if err != nil {
		return nil, err
	} else if builtinFilter == nil {
		return &ListSpansResp{Spans: loop_span.SpanList{}}, nil
	}
	filters := r.combineFilters(builtinFilter, req.Filters)
	tenants, err := r.getTenants(ctx, req.PlatformType)
	if err != nil {
		return nil, err
	}
	st := time.Now()
	tRes, err := r.traceRepo.ListSpans(ctx, &repo.ListSpansParam{
		Tenants:         tenants,
		Filters:         filters,
		StartAt:         req.StartTime,
		EndAt:           req.EndTime,
		Limit:           req.Limit,
		DescByStartTime: req.DescByStartTime,
		PageToken:       req.PageToken,
	})
	r.metrics.EmitListSpans(req.WorkspaceID, string(req.SpanListType), st, err != nil)
	if err != nil {
		return nil, err
	}
	processors, err := r.buildHelper.BuildListSpansProcessors(ctx, span_processor.Settings{
		WorkspaceId:    req.WorkspaceID,
		PlatformType:   req.PlatformType,
		QueryStartTime: req.StartTime,
		QueryEndTime:   req.EndTime,
	})
	if err != nil {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInternalErrorCodeCode)
	}
	spans := tRes.Spans
	for _, p := range processors {
		spans, err = p.Transform(ctx, spans)
		if err != nil {
			return nil, err
		}
	}
	return &ListSpansResp{
		Spans:         spans,
		NextPageToken: tRes.PageToken,
		HasMore:       tRes.HasMore,
	}, nil
}

func (r *TraceServiceImpl) IngestTraces(ctx context.Context, req *IngestTracesReq) error {
	traceData := &entity.TraceData{
		Tenant: r.traceConfig.GetDefaultTraceTenant(ctx),
		TenantInfo: entity.TenantInfo{
			TTL:              req.TTL,
			WorkspaceId:      req.Spans[0].WorkspaceID,
			CozeAccountID:    req.CozeAccountId,
			WhichIsEnough:    req.WhichIsEnough,
			VolcanoAccountID: req.VolcanoAccountID,
		},
		SpanList: req.Spans,
	}
	if err := r.traceProducer.IngestSpans(ctx, traceData); err != nil {
		return err
	}
	logs.CtxInfo(ctx, "Send msg successfully, spans count %d", len(req.Spans))
	return nil
}

func (r *TraceServiceImpl) GetTracesAdvanceInfo(ctx context.Context, req *GetTracesAdvanceInfoReq) (*GetTracesAdvanceInfoResp, error) {
	var (
		g                errgroup.Group
		lock             sync.Mutex
		defaultTimeRange = int64(15 * 60 * 1000) // ms
	)
	tenants, err := r.getTenants(ctx, req.PlatformType)
	if err != nil {
		return nil, err
	}
	resp := &GetTracesAdvanceInfoResp{
		Infos: []*loop_span.TraceAdvanceInfo{},
	}
	// use one processor...
	processors, err := r.buildHelper.BuildAdvanceInfoProcessors(ctx, span_processor.Settings{
		WorkspaceId:  req.WorkspaceID,
		PlatformType: req.PlatformType,
	})
	if err != nil {
		logs.CtxError(ctx, "Fail to build advance info processor, %v", err)
		return nil, err
	}
	for _, v := range req.Traces {
		g.Go(func() error {
			defer goroutine.Recovery(ctx)
			qReq := &repo.GetTraceParam{
				Tenants: tenants,
				TraceID: v.TraceID,
				StartAt: v.StartTime,
				EndAt:   v.StartTime + defaultTimeRange,
				Limit:   1000,
			}
			st := time.Now()
			spans, err := r.traceRepo.GetTrace(ctx, qReq)
			r.metrics.EmitGetTrace(req.WorkspaceID, st, err != nil)
			if err != nil {
				logs.CtxError(ctx, "Fail to get trace %v, %v", *qReq, err)
				return err
			}
			for _, p := range processors {
				spans, err = p.Transform(ctx, spans)
				if err != nil {
					logs.CtxWarn(ctx, "Fail to transform span, %v", err)
					return nil
				}
			}
			inputTokens, outputTokens, err := spans.Stat(ctx)
			if err != nil {
				logs.CtxWarn(ctx, "Fail to get spans stat, %v", err)
				return nil
			}
			lock.Lock()
			defer lock.Unlock()
			resp.Infos = append(resp.Infos, &loop_span.TraceAdvanceInfo{
				TraceId:    qReq.TraceID,
				InputCost:  inputTokens,
				OutputCost: outputTokens,
			})
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		logs.CtxError(ctx, "fail to get all trace advance info, %v", err)
		return nil, err
	}
	return resp, nil
}

func (r *TraceServiceImpl) GetTracesMetaInfo(ctx context.Context, req *GetTracesMetaInfoReq) (*GetTracesMetaInfoResp, error) {
	cfg, err := r.traceConfig.GetTraceFieldMetaInfo(ctx)
	if err != nil {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInternalErrorCodeCode)
	}
	fields, ok := cfg.FieldMetas[req.PlatformType][req.SpanListType]
	if !ok {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("meta info not found"))
	}
	fieldMetas := make(map[string]*config.FieldMeta)
	for _, field := range fields {
		fieldMta, ok := cfg.AvailableFields[field]
		if !ok || fieldMta == nil {
			return nil, errorx.NewByCode(obErrorx.CommercialCommonInternalErrorCodeCode)
		}
		fieldMetas[field] = fieldMta
	}
	return &GetTracesMetaInfoResp{
		FilesMetas: fieldMetas,
	}, nil
}

func (r *TraceServiceImpl) buildBuiltinFilters(ctx context.Context, f span_filter.Filter, req *ListSpansReq) (*loop_span.FilterFields, error) {
	filters := make([]*loop_span.FilterField, 0)
	env := &span_filter.SpanEnv{
		WorkspaceId: req.WorkspaceID,
	}
	basicFilter, err := f.BuildBasicSpanFilter(ctx, env)
	if err != nil {
		return nil, err
	} else if len(basicFilter) == 0 { // if it's null, no need to query from ck
		return nil, nil
	}
	filters = append(filters, basicFilter...)
	switch req.SpanListType {
	case loop_span.SpanListTypeRootSpan:
		subFilter, err := f.BuildRootSpanFilter(ctx, env)
		if err != nil {
			return nil, err
		}
		filters = append(filters, subFilter...)
	case loop_span.SpanListTypeLLMSpan:
		subFilter, err := f.BuildLLMSpanFilter(ctx, env)
		if err != nil {
			return nil, err
		}
		filters = append(filters, subFilter...)
	case loop_span.SpanListTypeAllSpan:
		subFilter, err := f.BuildALLSpanFilter(ctx, env)
		if err != nil {
			return nil, err
		}
		filters = append(filters, subFilter...)
	default:
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid span list type: %s"))
	}
	// not supposed to be here
	if len(filters) == 0 {
		return nil, nil
	}
	filterAggr := &loop_span.FilterFields{
		QueryAndOr:   ptr.Of(loop_span.QueryAndOrEnumAnd),
		FilterFields: filters,
	}
	return filterAggr, nil
}

func (r *TraceServiceImpl) combineFilters(filters ...*loop_span.FilterFields) *loop_span.FilterFields {
	filterAggr := &loop_span.FilterFields{
		QueryAndOr: ptr.Of(loop_span.QueryAndOrEnumAnd),
	}
	for _, f := range filters {
		if f == nil {
			continue
		}
		filterAggr.FilterFields = append(filterAggr.FilterFields, &loop_span.FilterField{
			QueryAndOr: ptr.Of(loop_span.QueryAndOrEnumAnd),
			SubFilter:  f,
		})
	}
	return filterAggr
}

func (r *TraceServiceImpl) getTenants(ctx context.Context, platform loop_span.PlatformType) ([]string, error) {
	cfg, err := r.traceConfig.GetPlatformTenants(ctx)
	if err != nil {
		logs.CtxError(ctx, "fail to get platform tenants, %v", err)
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInternalErrorCodeCode)
	}
	if tenants, ok := cfg.Config[string(platform)]; ok {
		return tenants, nil
	} else {
		logs.CtxError(ctx, "tenant not found for platform %s", platform)
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("tenant not found for the platform"))
	}
}

func processSpecificFilter(f *loop_span.FilterField) error {
	switch f.FieldName {
	case loop_span.SpanFieldStatus:
		if err := processStatusFilter(f); err != nil {
			return err
		}
	case loop_span.SpanFieldDuration,
		loop_span.SpanFieldLatencyFirstResp,
		loop_span.SpanFieldStartTimeFirstResp,
		loop_span.SpanFieldStartTimeFirstTokenResp,
		loop_span.SpanFieldLatencyFirstTokenResp,
		loop_span.SpanFieldReasoningDuration:
		if err := processLatencyFilter(f); err != nil {
			return err
		}
	}
	return nil
}

func processStatusFilter(f *loop_span.FilterField) error {
	if f.QueryType == nil || *f.QueryType != loop_span.QueryTypeEnumIn {
		return fmt.Errorf("status filter should use in operator")
	}
	f.FieldName = loop_span.SpanFieldStatusCode
	f.FieldType = loop_span.FieldTypeLong
	checkSuccess, checkError := false, false
	for _, val := range f.Values {
		switch val {
		case loop_span.SpanStatusSuccess:
			checkSuccess = true
		case loop_span.SpanStatusError:
			checkError = true
		default:
			return fmt.Errorf("invalid status code field value")
		}
	}
	if checkSuccess && checkError {
		f.QueryType = ptr.Of(loop_span.QueryTypeEnumAlwaysTrue)
		f.Values = nil
	} else if checkSuccess {
		f.Values = []string{"0"}
	} else if checkError {
		f.QueryType = ptr.Of(loop_span.QueryTypeEnumNotIn)
		f.Values = []string{"0"}
	} else {
		return fmt.Errorf("invalid status code query")
	}
	return nil
}

// ms -> us
func processLatencyFilter(f *loop_span.FilterField) error {
	if f.FieldType != loop_span.FieldTypeLong {
		return fmt.Errorf("latency field type should be long ")
	}
	micros := make([]string, 0)
	for _, val := range f.Values {
		integer, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fmt.Errorf("fail to parse long value %s, %v", val, err)
		}
		integer = time_util.MillSec2MicroSec(integer)
		micros = append(micros, strconv.FormatInt(integer, 10))
	}
	f.Values = micros
	return nil
}

//go:generate mockgen -destination=mocks/span_processor.go -package=mocks . TraceFilterProcessorBuilder
type TraceFilterProcessorBuilder interface {
	BuildPlatformRelatedFilter(context.Context, loop_span.PlatformType) (span_filter.Filter, error)
	BuildGetTraceProcessors(context.Context, span_processor.Settings) ([]span_processor.Processor, error)
	BuildListSpansProcessors(context.Context, span_processor.Settings) ([]span_processor.Processor, error)
	BuildAdvanceInfoProcessors(context.Context, span_processor.Settings) ([]span_processor.Processor, error)
}

type TraceFilterProcessorBuilderImpl struct {
	platformFilterFactory         span_filter.PlatformFilterFactory
	getTraceProcessorFactories    []span_processor.Factory
	listSpansProcessorFactories   []span_processor.Factory
	advanceInfoProcessorFactories []span_processor.Factory
}

func (t *TraceFilterProcessorBuilderImpl) BuildPlatformRelatedFilter(
	ctx context.Context,
	platformType loop_span.PlatformType,
) (span_filter.Filter, error) {
	return t.platformFilterFactory.GetFilter(ctx, platformType)
}

func (t *TraceFilterProcessorBuilderImpl) BuildGetTraceProcessors(
	ctx context.Context,
	set span_processor.Settings,
) ([]span_processor.Processor, error) {
	ret := make([]span_processor.Processor, 0)
	for _, factory := range t.getTraceProcessorFactories {
		p, err := factory.CreateProcessor(ctx, set)
		if err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	return ret, nil
}

func (t *TraceFilterProcessorBuilderImpl) BuildListSpansProcessors(
	ctx context.Context,
	set span_processor.Settings,
) ([]span_processor.Processor, error) {
	ret := make([]span_processor.Processor, 0)
	for _, factory := range t.listSpansProcessorFactories {
		p, err := factory.CreateProcessor(ctx, set)
		if err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	return ret, nil
}

func (t *TraceFilterProcessorBuilderImpl) BuildAdvanceInfoProcessors(
	ctx context.Context,
	set span_processor.Settings,
) ([]span_processor.Processor, error) {
	ret := make([]span_processor.Processor, 0)
	for _, factory := range t.advanceInfoProcessorFactories {
		p, err := factory.CreateProcessor(ctx, set)
		if err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	return ret, nil
}

func NewTraceFilterProcessorBuilder(
	platformFilterFactory span_filter.PlatformFilterFactory,
	getTraceProcessorFactories []span_processor.Factory,
	listSpansProcessorFactories []span_processor.Factory,
	advanceInfoProcessorFactories []span_processor.Factory,
) TraceFilterProcessorBuilder {
	return &TraceFilterProcessorBuilderImpl{
		platformFilterFactory:         platformFilterFactory,
		getTraceProcessorFactories:    getTraceProcessorFactories,
		listSpansProcessorFactories:   listSpansProcessorFactories,
		advanceInfoProcessorFactories: advanceInfoProcessorFactories,
	}
}
