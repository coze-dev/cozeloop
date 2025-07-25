// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package consumer

import (
	"context"

	"github.com/bytedance/sonic"

	"github.com/coze-dev/cozeloop/backend/infra/middleware/session"
	"github.com/coze-dev/cozeloop/backend/infra/mq"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/service"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/conv"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

func NewExptRecordEvalConsumer(recordEval service.ExptItemEvalEvent) mq.IConsumerHandler {
	return &ExptRecordEvalConsumer{
		ExptItemEvalEvent: recordEval,
	}
}

type ExptRecordEvalConsumer struct {
	service.ExptItemEvalEvent
}

func (e *ExptRecordEvalConsumer) HandleMessage(ctx context.Context, ext *mq.MessageExt) error {
	body := ext.Body

	event := &entity.ExptItemEvalEvent{}
	if err := sonic.Unmarshal(body, event); err != nil {
		logs.CtxError(ctx, "ExptItemEvalEvent json unmarshal fail, raw: %v, err: %s", conv.UnsafeBytesToString(body), err)
		return nil
	}

	logs.CtxInfo(ctx, "ExptRecordEvalConsumer consume message, event: %v, msg_id: %v", conv.UnsafeBytesToString(body), ext.MsgID)

	if event.Session != nil {
		ctx = session.WithCtxUser(ctx, &session.User{
			ID: event.Session.UserID,
		})
	}

	return e.Eval(ctx, event)
}
