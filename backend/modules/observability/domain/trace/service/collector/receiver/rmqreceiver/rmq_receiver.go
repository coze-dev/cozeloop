// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package rmqreceiver

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coze-dev/cozeloop/backend/infra/mq"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/component"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/consumer"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/loop_span"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

type rmqReceiver struct {
	componentID  component.ID
	nextConsumer consumer.Consumer
	config       *Config
	mqFactory    mq.IFactory
	consumer     mq.IConsumer
}

func (r *rmqReceiver) Start(ctx context.Context) error {
	logs.CtxInfo(ctx, "rmq receiver starting")
	if r.mqFactory == nil {
		return fmt.Errorf("rmq receiver factory not initialized")
	}
	mqConsumer, err := r.mqFactory.NewConsumer(mq.ConsumerConfig{
		Addr:           r.config.Addr,
		Topic:          r.config.Topic,
		ConsumerGroup:  r.config.ConsumerGroup,
		ConsumeTimeout: time.Duration(r.config.Timeout) * time.Second,
	})
	if err != nil {
		logs.CtxError(ctx, "rmq receiver failed to initialize consumer", err)
		return err
	}
	r.consumer = mqConsumer
	r.consumer.RegisterHandler(r)
	if err := r.consumer.Start(); err != nil {
		logs.CtxError(ctx, "rmq receiver consumer start err: %v", err)
		return err
	}
	return nil
}

func (r *rmqReceiver) Shutdown(ctx context.Context) error {
	if r.consumer == nil {
		return nil
	}
	logs.CtxInfo(ctx, "rmq receiver shutting down")
	return r.consumer.Close()
}

func (r *rmqReceiver) HandleMessage(ctx context.Context, msg *mq.MessageExt) error {
	traceData := new(entity.TraceData)
	if err := json.Unmarshal(msg.Body, traceData); err != nil {
		logs.CtxError(ctx, "fail to unmarshal message, %v", err)
		return err
	}
	spanList := make(loop_span.SpanList, 0)
	for _, span := range traceData.SpanList {
		if err := span.IsValidSpan(); err != nil {
			logs.CtxError(ctx, "rmqReceiver: invalid span found: %v", err)
			continue
		}
		spanList = append(spanList, span)
	}
	if len(spanList) == 0 {
		logs.CtxInfo(ctx, "rmqReceiver: no valid spans remains, just skip")
		return nil
	}
	td := consumer.Traces{
		Tenant:    traceData.Tenant,
		TraceData: []*entity.TraceData{traceData},
	}
	if err := r.nextConsumer.ConsumeTraces(ctx, td); err != nil {
		logs.CtxError(ctx, "rmqReceiver: next consumer consume traces failed: %v", err)
		return err
	}
	return nil
}
