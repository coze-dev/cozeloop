// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package consumer

import (
	"context"

	"github.com/coze-dev/cozeloop/backend/infra/middleware/session"
	"github.com/coze-dev/cozeloop/backend/infra/mq"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/service"
	"github.com/coze-dev/cozeloop/backend/pkg/json"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/conv"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

func newExptSchedulerConsumer(scheduler service.ExptSchedulerEvent) mq.IConsumerHandler {
	return &ExptSchedulerConsumer{
		scheduler: scheduler,
	}
}

type ExptSchedulerConsumer struct {
	scheduler service.ExptSchedulerEvent
}

func (e *ExptSchedulerConsumer) HandleMessage(ctx context.Context, msg *mq.MessageExt) error {
	event := &entity.ExptScheduleEvent{}
	body := msg.Body
	if err := json.Unmarshal(body, event); err != nil {
		logs.CtxError(ctx, "ExptExecEvent json unmarshal fail, raw: %v, err: %s", conv.UnsafeBytesToString(body), err)
		return nil
	}

	logs.CtxInfo(ctx, "ExptSchedulerConsumer consume message, event: %v, msg_id: %v", conv.UnsafeBytesToString(body), msg.MsgID)

	if event.Session != nil {
		ctx = session.WithCtxUser(ctx, &session.User{
			ID: event.Session.UserID,
		})
	}

	return e.scheduler.Schedule(ctx, event)
}
