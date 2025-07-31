// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/loop_span"
)

type GetTraceParam struct {
	Tenants []string
	TraceID string
	StartAt int64 // ms
	EndAt   int64 // ms
	Limit   int32
}

type ListSpansParam struct {
	Tenants         []string
	Filters         *loop_span.FilterFields
	StartAt         int64 // ms
	EndAt           int64 // ms
	Limit           int32
	DescByStartTime bool
	PageToken       string
}

type ListSpansResult struct {
	Spans     loop_span.SpanList
	PageToken string
	HasMore   bool
}

type InsertTraceParam struct {
	Spans  loop_span.SpanList
	Tenant string
	TTL    entity.TTL
}

//go:generate mockgen -destination=mocks/trace.go -package=mocks . ITraceRepo
type ITraceRepo interface {
	InsertSpans(context.Context, *InsertTraceParam) error
	ListSpans(context.Context, *ListSpansParam) (*ListSpansResult, error)
	GetTrace(context.Context, *GetTraceParam) (loop_span.SpanList, error)
}
