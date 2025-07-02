// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package looptracer

import (
	"context"

	cozeloop "github.com/coze-dev/cozeloop-go"
)

var _ Tracer = (*TracerImpl)(nil)

type TracerImpl struct {
	cozeloop.Client
}

func NewTracer(client cozeloop.Client) Tracer {
	return &TracerImpl{Client: client}
}

func (t *TracerImpl) StartSpan(ctx context.Context, name, spanType string, opts ...cozeloop.StartSpanOption) (context.Context, Span) {
	ctx, span := t.Client.StartSpan(ctx, name, spanType, opts...)
	return ctx, SpanImpl{
		LoopSpan: span,
	}
}

func (t *TracerImpl) GetSpanFromContext(ctx context.Context) Span {
	span := t.Client.GetSpanFromContext(ctx)
	return SpanImpl{
		LoopSpan: span,
	}
}

func (t *TracerImpl) Inject(ctx context.Context) context.Context {
	return ctx
}

type noopTracer struct {
	c cozeloop.Client
}

func (d *noopTracer) StartSpan(ctx context.Context, name, spanType string, opts ...cozeloop.StartSpanOption) (context.Context, Span) {
	return ctx, &noopSpan{}
}

func (d *noopTracer) GetSpanFromContext(ctx context.Context) Span {
	return &noopSpan{}
}

func (d *noopTracer) Flush(ctx context.Context) {}

func (d *noopTracer) Inject(ctx context.Context) context.Context {
	return ctx
}

func (d *noopTracer) SetCallType(callType string) {
}
