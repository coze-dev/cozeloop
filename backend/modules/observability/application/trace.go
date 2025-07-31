// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"strconv"

	"github.com/coze-dev/coze-loop/backend/infra/external/benefit"
	"github.com/coze-dev/coze-loop/backend/infra/middleware/session"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/observability/domain/common"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/observability/domain/filter"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/observability/domain/span"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/observability/domain/view"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/observability/trace"
	tconv "github.com/coze-dev/coze-loop/backend/modules/observability/application/convertor/trace"
	"github.com/coze-dev/coze-loop/backend/modules/observability/application/utils"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/component/config"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/component/metrics"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/component/rpc"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/loop_span"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/repo"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/service"
	obErrorx "github.com/coze-dev/coze-loop/backend/modules/observability/pkg/errno"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-loop/backend/pkg/logs"
)

const (
	MaxSpanLength     = 100
	MaxListSpansLimit = 1000
	QueryLimitDefault = 100
)

type ITraceApplication = trace.TraceService

func NewTraceApplication(
	traceService service.ITraceService,
	viewRepo repo.IViewRepo,
	authService rpc.IAuthProvider,
	benefitService benefit.IBenefitService,
	traceMetrics metrics.ITraceMetrics,
	traceConfig config.ITraceConfig,
) (ITraceApplication, error) {
	return &TraceApplication{
		traceService: traceService,
		viewRepo:     viewRepo,
		traceConfig:  traceConfig,
		auth:         authService,
		metrics:      traceMetrics,
		benefit:      benefitService,
	}, nil
}

type TraceApplication struct {
	traceService service.ITraceService
	viewRepo     repo.IViewRepo
	traceConfig  config.ITraceConfig
	auth         rpc.IAuthProvider
	metrics      metrics.ITraceMetrics
	benefit      benefit.IBenefitService
}

func (t *TraceApplication) ListSpans(ctx context.Context, req *trace.ListSpansRequest) (*trace.ListSpansResponse, error) {
	if err := t.validateListSpansReq(ctx, req); err != nil {
		return nil, err
	}
	if err := t.auth.CheckWorkspacePermission(ctx,
		rpc.AuthActionTraceRead,
		strconv.FormatInt(req.GetWorkspaceID(), 10)); err != nil {
		return nil, err
	}
	sReq, err := t.buildListSpansSvcReq(req)
	if err != nil {
		logs.CtxInfo(ctx, "invalid list spans request: %v", err)
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("list spans req is invalid"))
	}
	sResp, err := t.traceService.ListSpans(ctx, sReq)
	if err != nil {
		return nil, err
	}
	logs.CtxInfo(ctx, "List spans successfully, spans count: %d", len(sResp.Spans))
	return &trace.ListSpansResponse{
		Spans:         tconv.SpanListDO2DTO(sResp.Spans),
		NextPageToken: sResp.NextPageToken,
		HasMore:       sResp.HasMore,
	}, nil
}

func (t *TraceApplication) validateListSpansReq(ctx context.Context, req *trace.ListSpansRequest) error {
	if req == nil {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("no request provided"))
	} else if req.GetWorkspaceID() <= 0 {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid workspace_id"))
	} else if pageSize := req.GetPageSize(); pageSize < 0 || pageSize > MaxListSpansLimit {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid limit"))
	} else if len(req.GetOrderBys()) > 1 {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid order by %s"))
	}
	v := utils.DateValidator{
		Start:        req.GetStartTime(),
		End:          req.GetEndTime(),
		EarliestDays: t.traceConfig.GetTraceDataMaxDurationDay(ctx, req.PlatformType),
	}
	newStartTime, newEndTime, err := v.CorrectDate()
	if err != nil {
		return err
	}
	req.SetStartTime(newStartTime)
	req.SetEndTime(newEndTime)
	return nil
}

func (t *TraceApplication) buildListSpansSvcReq(req *trace.ListSpansRequest) (*service.ListSpansReq, error) {
	ret := &service.ListSpansReq{
		WorkspaceID:     req.GetWorkspaceID(),
		StartTime:       req.GetStartTime(),
		EndTime:         req.GetEndTime(),
		Limit:           QueryLimitDefault,
		DescByStartTime: len(req.GetOrderBys()) > 0,
		PageToken:       req.GetPageToken(),
	}
	if req.PageSize != nil {
		ret.Limit = *req.PageSize
	}
	platformType := loop_span.PlatformType(req.GetPlatformType())
	if req.PlatformType == nil {
		platformType = loop_span.PlatformCozeLoop
	}
	ret.PlatformType = platformType
	switch req.GetSpanListType() {
	case common.SpanListTypeRootSpan:
		ret.SpanListType = loop_span.SpanListTypeRootSpan
	case common.SpanListTypeAllSpan:
		ret.SpanListType = loop_span.SpanListTypeAllSpan
	case common.SpanListTypeLlmSpan:
		ret.SpanListType = loop_span.SpanListTypeLLMSpan
	default:
		ret.SpanListType = loop_span.SpanListTypeRootSpan
	}
	if req.Filters != nil {
		ret.Filters = tconv.FilterFieldsDTO2DO(req.Filters)
		if err := ret.Filters.Validate(); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (t *TraceApplication) GetTrace(ctx context.Context, req *trace.GetTraceRequest) (*trace.GetTraceResponse, error) {
	if err := t.validateGetTraceReq(ctx, req); err != nil {
		return nil, err
	}
	if err := t.auth.CheckWorkspacePermission(ctx,
		rpc.AuthActionTraceRead,
		strconv.FormatInt(req.GetWorkspaceID(), 10)); err != nil {
		return nil, err
	}
	logs.CtxInfo(ctx, "Get trace request: %+v", req)
	sReq := t.buildGetTraceSvcReq(req)
	sResp, err := t.traceService.GetTrace(ctx, sReq)
	if err != nil {
		return nil, err
	}
	inTokens, outTokens, err := sResp.Spans.Stat(ctx)
	if err != nil {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInternalErrorCodeCode)
	}
	logs.CtxInfo(ctx, "Get trace successfully, spans count %d", len(sResp.Spans))
	return &trace.GetTraceResponse{
		Spans: tconv.SpanListDO2DTO(sResp.Spans),
		TracesAdvanceInfo: &trace.TraceAdvanceInfo{
			TraceID: sResp.TraceId,
			Tokens: &trace.TokenCost{
				Input:  inTokens,
				Output: outTokens,
			},
		},
	}, nil
}

func (t *TraceApplication) validateGetTraceReq(ctx context.Context, req *trace.GetTraceRequest) error {
	if req == nil {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("no request provided"))
	} else if req.GetWorkspaceID() <= 0 {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid workspace_id"))
	} else if req.GetTraceID() == "" {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid trace_id"))
	}
	v := utils.DateValidator{
		Start:        req.GetStartTime(),
		End:          req.GetEndTime(),
		EarliestDays: t.traceConfig.GetTraceDataMaxDurationDay(ctx, req.PlatformType),
	}
	newStartTime, newEndTime, err := v.CorrectDate()
	if err != nil {
		return err
	}
	req.SetStartTime(newStartTime)
	req.SetEndTime(newEndTime)
	return nil
}

func (t *TraceApplication) buildGetTraceSvcReq(req *trace.GetTraceRequest) *service.GetTraceReq {
	ret := &service.GetTraceReq{
		WorkspaceID: req.GetWorkspaceID(),
		TraceID:     req.GetTraceID(),
		StartTime:   req.GetStartTime(),
		EndTime:     req.GetEndTime(),
	}
	platformType := loop_span.PlatformType(req.GetPlatformType())
	if req.PlatformType == nil {
		platformType = loop_span.PlatformCozeLoop
	}
	ret.PlatformType = platformType
	return ret
}

func (t *TraceApplication) BatchGetTracesAdvanceInfo(ctx context.Context, req *trace.BatchGetTracesAdvanceInfoRequest) (*trace.BatchGetTracesAdvanceInfoResponse, error) {
	if err := t.validateGetTracesAdvanceInfoReq(ctx, req); err != nil {
		return nil, err
	}
	if err := t.auth.CheckWorkspacePermission(ctx,
		rpc.AuthActionTraceRead,
		strconv.FormatInt(req.GetWorkspaceID(), 10)); err != nil {
		return nil, err
	}
	logs.CtxInfo(ctx, "Batch get traces advance info request: %+v", req)
	sReq := t.buildBatchGetTraceAdvanceInfoSvcReq(req)
	sResp, err := t.traceService.GetTracesAdvanceInfo(ctx, sReq)
	if err != nil {
		return nil, err
	}
	return &trace.BatchGetTracesAdvanceInfoResponse{
		TracesAdvanceInfo: tconv.BatchAdvanceInfoDO2DTO(sResp.Infos),
	}, nil
}

func (t *TraceApplication) validateGetTracesAdvanceInfoReq(ctx context.Context, req *trace.BatchGetTracesAdvanceInfoRequest) error {
	if req == nil {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("no request provided"))
	} else if req.GetWorkspaceID() <= 0 {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid workspace_id"))
	} else if len(req.GetTraces()) < 1 {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid traces"))
	}
	for _, tReq := range req.Traces {
		if tReq.GetTraceID() == "" {
			return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid trace_id"))
		}
		v := utils.DateValidator{
			Start:        tReq.GetStartTime(),
			End:          tReq.GetEndTime(),
			EarliestDays: t.traceConfig.GetTraceDataMaxDurationDay(ctx, req.PlatformType),
		}
		newStartTime, newEndTime, err := v.CorrectDate()
		if err != nil {
			return err
		}
		tReq.SetStartTime(newStartTime)
		tReq.SetEndTime(newEndTime)
	}
	return nil
}

func (t *TraceApplication) buildBatchGetTraceAdvanceInfoSvcReq(req *trace.BatchGetTracesAdvanceInfoRequest) *service.GetTracesAdvanceInfoReq {
	ret := &service.GetTracesAdvanceInfoReq{
		WorkspaceID: req.GetWorkspaceID(),
		Traces:      make([]*service.TraceQueryParam, len(req.GetTraces())),
	}
	for i, traceInfo := range req.GetTraces() {
		ret.Traces[i] = &service.TraceQueryParam{
			TraceID:   traceInfo.GetTraceID(),
			StartTime: traceInfo.GetStartTime(),
			EndTime:   traceInfo.GetEndTime(),
		}
	}
	platformType := loop_span.PlatformType(req.GetPlatformType())
	if req.PlatformType == nil {
		platformType = loop_span.PlatformCozeLoop
	}
	ret.PlatformType = platformType
	return ret
}

func (t *TraceApplication) IngestTraces(ctx context.Context, req *trace.IngestTracesRequest) (*trace.IngestTracesResponse, error) {
	if err := t.validateIngestTracesReq(ctx, req); err != nil {
		return nil, err
	}
	workspaceId := req.GetSpans()[0].WorkspaceID
	if err := t.auth.CheckWorkspacePermission(ctx,
		rpc.AuthActionTraceIngest,
		workspaceId); err != nil {
		return nil, err
	}
	workSpaceIdNum, err := strconv.ParseInt(workspaceId, 10, 64)
	if err != nil {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid workspace_id"))
	}
	connectorUid := session.UserIDInCtxOrEmpty(ctx)
	benefitRes, err := t.benefit.CheckTraceBenefit(ctx, &benefit.CheckTraceBenefitParams{
		ConnectorUID: connectorUid,
		SpaceID:      workSpaceIdNum,
	})
	if err != nil {
		logs.CtxError(ctx, "Fail to check benefit, %v", err)
	}
	if benefitRes == nil {
		benefitRes = &benefit.CheckTraceBenefitResult{
			AccountAvailable: true,
			IsEnough:         true,
			StorageDuration:  3,
			WhichIsEnough:    -1,
		}
	}
	if !benefitRes.IsEnough {
		return nil, errorx.NewByCode(obErrorx.TraceNoCapacityAvailableErrorCode)
	} else if !benefitRes.AccountAvailable {
		return nil, errorx.NewByCode(obErrorx.AccountNotAvailableErrorCode)
	}
	spans := tconv.SpanListDTO2DO(req.Spans)
	for _, s := range spans {
		s.CallType = "Custom"
	}
	if err := t.traceService.IngestTraces(ctx, &service.IngestTracesReq{
		TTL:              entity.TTLFromInteger(benefitRes.StorageDuration),
		WhichIsEnough:    benefitRes.WhichIsEnough,
		CozeAccountId:    connectorUid,
		VolcanoAccountID: benefitRes.VolcanoAccountID,
		Spans:            spans,
	}); err != nil {
		return nil, err
	}
	return trace.NewIngestTracesResponse(), nil
}

func (t *TraceApplication) IngestTracesInner(ctx context.Context, req *trace.IngestTracesRequest) (r *trace.IngestTracesResponse, err error) {
	if err := t.validateIngestTracesInnerReq(ctx, req); err != nil {
		return nil, err
	}
	// spaceId/UserId
	spansMap := make(map[string]map[string][]*span.InputSpan)
	for _, inputSpan := range req.Spans {
		if inputSpan == nil {
			continue
		}
		spaceId := inputSpan.WorkspaceID
		userId := inputSpan.TagsString[loop_span.SpanFieldUserID]
		if spansMap[spaceId] == nil {
			spansMap[spaceId] = make(map[string][]*span.InputSpan)
		}
		if spansMap[spaceId][userId] == nil {
			spansMap[spaceId][userId] = make([]*span.InputSpan, 0)
		}
		spansMap[spaceId][userId] = append(spansMap[spaceId][userId], inputSpan)
	}
	for spaceID, userIdSpansMap := range spansMap {
		for userId, spans := range userIdSpansMap {
			workspaceId := spaceID
			workSpaceIdNum, err := strconv.ParseInt(workspaceId, 10, 64)
			if err != nil {
				return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid workspace_id"))
			}
			benefitRes, err := t.benefit.CheckTraceBenefit(ctx, &benefit.CheckTraceBenefitParams{
				ConnectorUID: userId,
				SpaceID:      workSpaceIdNum,
			})
			if err != nil {
				logs.CtxError(ctx, "Fail to check benefit, %v", err)
			}
			if benefitRes == nil {
				benefitRes = &benefit.CheckTraceBenefitResult{
					AccountAvailable: true,
					IsEnough:         true,
					StorageDuration:  3,
					WhichIsEnough:    -1,
				}
			}
			if !benefitRes.IsEnough || !benefitRes.AccountAvailable {
				benefitRes.StorageDuration = 3
				logs.CtxWarn(ctx, "check benefit err: resource not enough")
			}
			spans := tconv.SpanListDTO2DO(spans)
			for _, s := range spans {
				callType, ok := s.TagsString[loop_span.SpanFieldCallType]
				if ok {
					s.CallType = callType
					delete(s.TagsString, loop_span.SpanFieldCallType)
				}
			}
			if err := t.traceService.IngestTraces(ctx, &service.IngestTracesReq{
				TTL:              entity.TTLFromInteger(benefitRes.StorageDuration),
				WhichIsEnough:    benefitRes.WhichIsEnough,
				CozeAccountId:    userId,
				VolcanoAccountID: benefitRes.VolcanoAccountID,
				Spans:            spans,
			}); err != nil {
				return nil, err
			}
		}
	}
	return trace.NewIngestTracesResponse(), nil
}

func (t *TraceApplication) validateIngestTracesInnerReq(ctx context.Context, req *trace.IngestTracesRequest) error {
	if req == nil {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("no request provided"))
	} else if len(req.Spans) > MaxSpanLength {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("max span length exceeded"))
	} else if len(req.Spans) < 1 {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("no spans provided"))
	}
	return nil
}

func (t *TraceApplication) validateIngestTracesReq(ctx context.Context, req *trace.IngestTracesRequest) error {
	if req == nil {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("no request provided"))
	} else if len(req.Spans) > MaxSpanLength {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("max span length exceeded"))
	} else if len(req.Spans) < 1 {
		return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("no spans provided"))
	}
	workspaceId := req.Spans[0].WorkspaceID
	for i := 1; i < len(req.Spans); i++ {
		if req.Spans[i].WorkspaceID != workspaceId {
			return errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("spans space id is not the same"))
		}
	}
	return nil
}

func (t *TraceApplication) GetTracesMetaInfo(ctx context.Context, req *trace.GetTracesMetaInfoRequest) (*trace.GetTracesMetaInfoResponse, error) {
	if err := t.auth.CheckWorkspacePermission(ctx,
		rpc.AuthActionTraceRead,
		strconv.FormatInt(req.GetWorkspaceID(), 10)); err != nil {
		return nil, err
	}
	logs.CtxInfo(ctx, "Get traces meta info request: %+v", req)
	sReq := t.buildGetTracesMetaInfoReq(req)
	sResp, err := t.traceService.GetTracesMetaInfo(ctx, sReq)
	if err != nil {
		return nil, err
	}
	fMeta := make(map[string]*trace.FieldMeta)
	for k, v := range sResp.FilesMetas {
		fMeta[k] = &trace.FieldMeta{
			ValueType:                 filter.FieldType(v.FieldType),
			SupportCustomizableOption: ptr.Of(v.SupportCustom),
		}
		if v.FieldOptions != nil {
			fMeta[k].FieldOptions = &filter.FieldOptions{
				I64List:    v.FieldOptions.I64List,
				F64List:    v.FieldOptions.F64List,
				StringList: v.FieldOptions.StringList,
			}
		}
		fTypes := make([]filter.FieldType, 0)
		for _, t := range v.FilterTypes {
			fTypes = append(fTypes, filter.FieldType(t))
		}
		fMeta[k].FilterTypes = fTypes
	}
	return &trace.GetTracesMetaInfoResponse{
		FieldMetas: fMeta,
	}, nil
}

func (t *TraceApplication) buildGetTracesMetaInfoReq(req *trace.GetTracesMetaInfoRequest) *service.GetTracesMetaInfoReq {
	ret := &service.GetTracesMetaInfoReq{
		WorkspaceID: req.GetWorkspaceID(),
	}
	platformType := loop_span.PlatformType(req.GetPlatformType())
	if req.PlatformType == nil {
		platformType = loop_span.PlatformCozeLoop
	}
	ret.PlatformType = platformType
	switch req.GetSpanListType() {
	case common.SpanListTypeRootSpan:
		ret.SpanListType = loop_span.SpanListTypeRootSpan
	case common.SpanListTypeAllSpan:
		ret.SpanListType = loop_span.SpanListTypeAllSpan
	case common.SpanListTypeLlmSpan:
		ret.SpanListType = loop_span.SpanListTypeLLMSpan
	default:
		ret.SpanListType = loop_span.SpanListTypeRootSpan
	}
	return ret
}

func (t *TraceApplication) CreateView(ctx context.Context, req *trace.CreateViewRequest) (*trace.CreateViewResponse, error) {
	if req == nil {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("no request provided"))
	} else if req.GetWorkspaceID() <= 0 {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid workspace_id"))
	} else if req.ViewName == "" {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid view_name"))
	}
	if err := t.auth.CheckWorkspacePermission(ctx,
		rpc.AuthActionTraceViewCreate,
		strconv.FormatInt(req.GetWorkspaceID(), 10)); err != nil {
		return nil, err
	}
	userID := session.UserIDInCtxOrEmpty(ctx)
	if userID == "" {
		return nil, errorx.NewByCode(obErrorx.UserParseFailedCode)
	}
	viewPO := tconv.CreateViewDTO2PO(req, userID)
	logs.CtxInfo(ctx, "Create view %v", *viewPO)
	id, err := t.viewRepo.CreateView(ctx, viewPO)
	if err != nil {
		return nil, err
	}
	logs.CtxInfo(ctx, "Create view successfully")
	return &trace.CreateViewResponse{
		ID: id,
	}, nil
}

func (t *TraceApplication) UpdateView(ctx context.Context, req *trace.UpdateViewRequest) (*trace.UpdateViewResponse, error) {
	if req == nil {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("no request provided"))
	} else if req.GetWorkspaceID() <= 0 {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid workspace_id"))
	}
	if err := t.auth.CheckViewPermission(ctx,
		rpc.AuthActionTraceViewEdit,
		strconv.FormatInt(req.GetWorkspaceID(), 10),
		strconv.FormatInt(req.GetID(), 10)); err != nil {
		return nil, err
	}
	userID := session.UserIDInCtxOrEmpty(ctx)
	if userID == "" {
		return nil, errorx.NewByCode(obErrorx.UserParseFailedCode)
	}
	viewDo, err := t.viewRepo.GetView(ctx, req.GetID(), ptr.Of(req.GetWorkspaceID()), ptr.Of(userID))
	if err != nil {
		return nil, err
	}
	logs.CtxInfo(ctx, "Get original view %v", *viewDo)
	if req.ViewName != nil {
		viewDo.ViewName = *req.ViewName
	}
	if req.Filters != nil {
		viewDo.Filters = *req.Filters
	}
	if req.PlatformType != nil {
		viewDo.PlatformType = *req.PlatformType
	}
	if req.SpanListType != nil {
		viewDo.SpanListType = *req.SpanListType
	}
	logs.CtxInfo(ctx, "Update view %d into %v", req.GetID(), *viewDo)
	if err := t.viewRepo.UpdateView(ctx, viewDo); err != nil {
		return nil, err
	}
	logs.CtxInfo(ctx, "Update view successfully")
	return trace.NewUpdateViewResponse(), nil
}

func (t *TraceApplication) DeleteView(ctx context.Context, req *trace.DeleteViewRequest) (*trace.DeleteViewResponse, error) {
	if req == nil {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("no request provided"))
	} else if req.GetID() <= 0 || req.GetWorkspaceID() <= 0 {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid workspace_id"))
	}
	if err := t.auth.CheckViewPermission(ctx,
		rpc.AuthActionTraceViewEdit,
		strconv.FormatInt(req.GetWorkspaceID(), 10),
		strconv.FormatInt(req.GetID(), 10)); err != nil {
		return nil, err
	}
	userID := session.UserIDInCtxOrEmpty(ctx)
	if userID == "" {
		return nil, errorx.NewByCode(obErrorx.UserParseFailedCode)
	}
	logs.CtxInfo(ctx, "Delete view %d at %d by %s", req.GetID(), req.GetWorkspaceID(), userID)
	if err := t.viewRepo.DeleteView(ctx, req.GetID(), req.GetWorkspaceID(), userID); err != nil {
		return nil, err
	}
	logs.CtxInfo(ctx, "Delete view successfully")
	return trace.NewDeleteViewResponse(), nil
}

func (t *TraceApplication) ListViews(ctx context.Context, req *trace.ListViewsRequest) (*trace.ListViewsResponse, error) {
	if req == nil {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("no request provided"))
	} else if req.GetWorkspaceID() <= 0 {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid workspace_id"))
	}
	if err := t.auth.CheckWorkspacePermission(ctx,
		rpc.AuthActionTraceViewList,
		strconv.FormatInt(req.GetWorkspaceID(), 10)); err != nil {
		return nil, err
	}
	systemViews, err := t.getSystemViews(ctx)
	if err != nil {
		return nil, err
	}
	userID := session.UserIDInCtxOrEmpty(ctx)
	if userID == "" {
		return nil, errorx.NewByCode(obErrorx.UserParseFailedCode)
	}
	logs.CtxInfo(ctx, "List views for %s at %d", userID, req.GetWorkspaceID())
	viewList, err := t.viewRepo.ListViews(ctx, req.WorkspaceID, userID)
	if err != nil {
		return nil, err
	}
	logs.CtxInfo(ctx, "List views successfully")
	return &trace.ListViewsResponse{
		Views:    append(systemViews, tconv.BatchViewPO2DTO(viewList)...),
		BaseResp: nil,
	}, nil
}

func (t *TraceApplication) getSystemViews(ctx context.Context) ([]*view.View, error) {
	systemViews, err := t.traceConfig.GetSystemViews(ctx)
	if err != nil {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInternalErrorCodeCode, errorx.WithExtraMsg("get system views failed"))
	}
	ret := make([]*view.View, 0)
	for _, v := range systemViews {
		ret = append(ret, &view.View{
			ID:           v.ID,
			ViewName:     v.ViewName,
			Filters:      v.Filters,
			PlatformType: ptr.Of(common.PlatformTypeCozeloop),
			SpanListType: ptr.Of(common.SpanListTypeRootSpan),
			IsSystem:     true,
		})
	}
	return ret, nil
}
