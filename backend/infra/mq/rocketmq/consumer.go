// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package rocketmq

import (
	"context"
	"errors"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/coze-dev/cozeloop/backend/infra/mq"
)

type Consumer struct {
	consumer rocketmq.PushConsumer
	handler  mq.IConsumerHandler
	topic    string
	selector *consumer.MessageSelector
}

func (c *Consumer) Start() error {
	if c.handler == nil {
		return errors.New("handler not set")
	}
	selector := consumer.MessageSelector{}
	if c.selector != nil {
		selector = *c.selector
	}

	err := c.consumer.Subscribe(c.topic, selector, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			// 转换消息格式
			ext := &mq.MessageExt{
				Message: mq.Message{
					Topic:        msg.Topic,
					Body:         msg.Body,
					Tag:          msg.GetTags(),
					PartitionKey: msg.GetShardingKey(),
					Properties:   msg.GetProperties(),
				},
				MsgID: msg.MsgId,
			}

			// 处理业务逻辑
			if err := c.handler.HandleMessage(ctx, ext); err != nil {
				return consumer.ConsumeRetryLater, err
			}
		}
		return consumer.ConsumeSuccess, nil
	})

	if err == nil {
		err = c.consumer.Start()
	}
	return err
}

func (c *Consumer) Close() error {
	return c.consumer.Shutdown()
}

func (c *Consumer) RegisterHandler(h mq.IConsumerHandler) {
	c.handler = h
}
