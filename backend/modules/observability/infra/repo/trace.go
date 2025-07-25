// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/component/config"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/loop_span"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/repo"
	"github.com/coze-dev/cozeloop/backend/modules/observability/infra/repo/ck"
	"github.com/coze-dev/cozeloop/backend/modules/observability/infra/repo/ck/convertor"
	obErrorx "github.com/coze-dev/cozeloop/backend/modules/observability/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/json"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
	time_util "github.com/coze-dev/cozeloop/backend/pkg/time"
)

func NewTraceCKRepoImpl(spanDao ck.ISpansDao, traceConfig config.ITraceConfig) (repo.ITraceRepo, error) {
	return &TraceCkRepoImpl{
		spansDao:    spanDao,
		traceConfig: traceConfig,
	}, nil
}

type TraceCkRepoImpl struct {
	spansDao    ck.ISpansDao
	traceConfig config.ITraceConfig
}

type PageToken struct {
	StartTime int64  `json:"StartTime"`
	SpanID    string `json:"SpanID"`
}

func (t *TraceCkRepoImpl) InsertSpans(ctx context.Context, param *repo.InsertTraceParam) error {
	table, err := t.getInsertTenantTable(ctx, param.Tenant, param.TTL)
	if err != nil {
		return err
	}
	if err := t.spansDao.Insert(ctx, &ck.InsertParam{
		Table: table,
		Spans: convertor.SpanListDO2PO(param.Spans, param.TTL),
	}); err != nil {
		logs.CtxError(ctx, "fail to insert spans, %v", err)
		return err
	}
	logs.CtxInfo(ctx, "insert spans into table %s successfully, count %d", table, len(param.Spans))
	return nil
}

func (t *TraceCkRepoImpl) ListSpans(ctx context.Context, req *repo.ListSpansParam) (*repo.ListSpansResult, error) {
	pageToken, err := parsePageToken(req.PageToken)
	if err != nil {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("invalid list spans request"))
	}
	if pageToken != nil {
		req.Filters = t.addPageTokenFilter(pageToken, req.Filters)
	}
	tables, err := t.getQueryTenantTables(ctx, req.Tenants)
	if err != nil {
		return nil, err
	}
	spans, err := t.spansDao.Get(ctx, &ck.QueryParam{
		QueryType:        ck.QueryTypeListSpans,
		Tables:           tables,
		StartTime:        time_util.MillSec2MicroSec(req.StartAt),
		EndTime:          time_util.MillSec2MicroSec(req.EndAt),
		Filters:          req.Filters,
		Limit:            req.Limit + 1,
		OrderByStartTime: req.DescByStartTime,
	})
	if err != nil {
		return nil, err
	}
	result := &repo.ListSpansResult{
		Spans:   convertor.SpanListPO2DO(spans),
		HasMore: len(spans) > int(req.Limit),
	}
	if result.HasMore {
		result.Spans = result.Spans[:len(result.Spans)-1]
	}
	if req.DescByStartTime && len(result.Spans) > 0 {
		lastSpan := result.Spans[len(result.Spans)-1]
		pageToken := &PageToken{
			StartTime: lastSpan.StartTime,
			SpanID:    lastSpan.SpanID,
		}
		pt, _ := json.Marshal(pageToken)
		result.PageToken = base64.StdEncoding.EncodeToString(pt)
	}
	return result, nil
}

func (t *TraceCkRepoImpl) GetTrace(ctx context.Context, req *repo.GetTraceParam) (loop_span.SpanList, error) {
	tables, err := t.getQueryTenantTables(ctx, req.Tenants)
	if err != nil {
		return nil, err
	}
	filter := &loop_span.FilterFields{
		QueryAndOr: ptr.Of(loop_span.QueryAndOrEnumAnd),
		FilterFields: []*loop_span.FilterField{
			{
				FieldName: loop_span.SpanFieldTraceId,
				FieldType: loop_span.FieldTypeString,
				Values:    []string{req.TraceID},
				QueryType: ptr.Of(loop_span.QueryTypeEnumEq),
			},
		},
	}
	st := time.Now()
	spans, err := t.spansDao.Get(ctx, &ck.QueryParam{
		QueryType: ck.QueryTypeGetTrace,
		Tables:    tables,
		StartTime: time_util.MillSec2MicroSec(req.StartAt),
		EndTime:   time_util.MillSec2MicroSec(req.EndAt),
		Filters:   filter,
		Limit:     req.Limit,
	})
	if err != nil {
		return nil, err
	}
	logs.CtxInfo(ctx, "get trace %s successfully, spans count %d, cost %v",
		req.TraceID, len(spans), time.Since(st))
	return convertor.SpanListPO2DO(spans), nil
}

func (t *TraceCkRepoImpl) getQueryTenantTables(ctx context.Context, tenants []string) ([]string, error) {
	tenantTableCfg, err := t.traceConfig.GetTenantConfig(ctx)
	if err != nil {
		logs.CtxError(ctx, "fail to get tenant table config, %v", err)
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInternalErrorCodeCode)
	}
	if len(tenants) == 0 {
		return nil, errorx.NewByCode(obErrorx.CommercialCommonInvalidParamCodeCode, errorx.WithExtraMsg("no tenants configured"))
	}
	tablesMap := make(map[string]bool)
	ret := make([]string, 0)
	for _, tenant := range tenants {
		tables := tenantTableCfg.QueryTables[tenant]
		if tables == nil {
			return nil, fmt.Errorf("no table config found for tenant %s", tenant)
		}
		for _, table := range tables {
			if tablesMap[table] {
				continue
			}
			ret = append(ret, table)
			tablesMap[table] = true
		}
	}
	return ret, nil
}

func (t *TraceCkRepoImpl) getInsertTenantTable(ctx context.Context, tenant string, ttl entity.TTL) (string, error) {
	tenantTableCfg, err := t.traceConfig.GetTenantConfig(ctx)
	if err != nil {
		logs.CtxError(ctx, "fail to get tenant config, %v", err)
		return "", err
	}
	table, ok := tenantTableCfg.InsertTable[tenant][ttl]
	if !ok {
		return "", fmt.Errorf("no table config found for tenant %s with ttl %s", tenant, ttl)
	}
	return table, nil
}

func (t *TraceCkRepoImpl) addPageTokenFilter(pageToken *PageToken, filter *loop_span.FilterFields) *loop_span.FilterFields {
	timeStr := strconv.FormatInt(pageToken.StartTime, 10)
	filterFields := &loop_span.FilterFields{
		QueryAndOr: ptr.Of(loop_span.QueryAndOrEnumOr),
		FilterFields: []*loop_span.FilterField{
			{
				FieldName: loop_span.SpanFieldStartTime,
				FieldType: loop_span.FieldTypeLong,
				Values:    []string{timeStr},
				QueryType: ptr.Of(loop_span.QueryTypeEnumLt),
			},
			{
				FieldName:  loop_span.SpanFieldStartTime,
				FieldType:  loop_span.FieldTypeLong,
				Values:     []string{timeStr},
				QueryType:  ptr.Of(loop_span.QueryTypeEnumEq),
				QueryAndOr: ptr.Of(loop_span.QueryAndOrEnumAnd),
				SubFilter: &loop_span.FilterFields{
					QueryAndOr: ptr.Of(loop_span.QueryAndOrEnumAnd),
					FilterFields: []*loop_span.FilterField{
						{
							FieldName: loop_span.SpanFieldSpanId,
							FieldType: loop_span.FieldTypeString,
							Values:    []string{pageToken.SpanID},
							QueryType: ptr.Of(loop_span.QueryTypeEnumLt),
						},
					},
				},
			},
		},
	}
	if filter == nil {
		return filterFields
	} else {
		return &loop_span.FilterFields{
			QueryAndOr: ptr.Of(loop_span.QueryAndOrEnumAnd),
			FilterFields: []*loop_span.FilterField{
				{
					SubFilter: filterFields,
				},
				{
					SubFilter: filter,
				},
			},
		}
	}
}

func parsePageToken(pageToken string) (*PageToken, error) {
	if pageToken == "" {
		return nil, nil
	}
	ptStr, err := base64.StdEncoding.DecodeString(pageToken)
	if err != nil {
		return nil, fmt.Errorf("fail to decode pageToken %s, %v", pageToken, err)
	}
	pt := new(PageToken)
	if err := json.Unmarshal(ptStr, pt); err != nil {
		return nil, fmt.Errorf("fail to unmarshal pageToken %s, %v", string(ptStr), err)
	}
	return pt, nil
}
