// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"strconv"
	"time"

	"github.com/coze-dev/cozeloop/backend/infra/metrics"
	metrics2 "github.com/coze-dev/cozeloop/backend/modules/observability/domain/component/metrics"
)

const (
	traceSpansMetricsName = "trace_spans"

	getTraceSuffix  = "get_trace"
	listSpansSuffix = "list_spans"

	throughputSuffix = ".throughput"
	latencySuffix    = ".latency"
)

const (
	tagSpaceID  = "workspace_id"
	tagSpanType = "span_type"
	tagIsErr    = "is_err"
)

func traceQueryTagNames() []string {
	return []string{
		tagSpaceID,
		tagSpanType,
		tagIsErr,
	}
}

func NewTraceMetricsImpl(meter metrics.Meter) metrics2.ITraceMetrics {
	ret := &TraceMetricsImpl{}
	if meter == nil {
		return ret
	}
	spansMetrics, err := meter.NewMetric(traceSpansMetricsName, []metrics.MetricType{metrics.MetricTypeCounter, metrics.MetricTypeTimer}, traceQueryTagNames())
	if err != nil {
		return ret
	}
	return &TraceMetricsImpl{
		spansMetrics: spansMetrics,
	}
}

type TraceMetricsImpl struct {
	spansMetrics metrics.Metric
}

func (t *TraceMetricsImpl) EmitListSpans(workspaceId int64, spanType string, start time.Time, isError bool) {
	if t.spansMetrics == nil {
		return
	}
	t.spansMetrics.Emit(
		[]metrics.T{
			{Name: tagSpaceID, Value: strconv.FormatInt(workspaceId, 10)},
			{Name: tagIsErr, Value: strconv.FormatBool(isError)},
			{Name: tagSpanType, Value: spanType},
		},
		metrics.Counter(1, metrics.WithSuffix(listSpansSuffix+throughputSuffix)),
		metrics.Timer(time.Since(start).Microseconds(), metrics.WithSuffix(listSpansSuffix+latencySuffix)))
}

func (t *TraceMetricsImpl) EmitGetTrace(workspaceId int64, start time.Time, isError bool) {
	if t.spansMetrics == nil {
		return
	}
	t.spansMetrics.Emit(
		[]metrics.T{
			{Name: tagSpaceID, Value: strconv.FormatInt(workspaceId, 10)},
			{Name: tagIsErr, Value: strconv.FormatBool(isError)},
		},
		metrics.Counter(1, metrics.WithSuffix(getTraceSuffix+throughputSuffix)),
		metrics.Timer(time.Since(start).Microseconds(), metrics.WithSuffix(getTraceSuffix+latencySuffix)))
}
