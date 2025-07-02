// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package producer

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"

	"github.com/coze-dev/cozeloop/backend/infra/mq"
	config "github.com/coze-dev/cozeloop/backend/modules/data/domain/component/conf"
	imq "github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/component/mq"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

type DatasetJobPublisher struct {
	Topic    string
	Tag      string
	producer mq.IProducer
}

func NewDatasetJobPublisher(iconfiger config.IConfig, mqFactory mq.IFactory) (imq.IDatasetJobPublisher, error) {
	conf := iconfiger.GetProducerConfig()
	producer, err := mqFactory.NewProducer(mq.ProducerConfig{
		Addr:           conf.Addr,
		ProduceTimeout: conf.ProduceTimeout,
		ProducerGroup:  &conf.ProducerGroup,
	})
	if err != nil {
		logs.Error("new mq producer fail, cfg: %v", conf)
		return nil, err
	}
	publisher := &DatasetJobPublisher{
		Topic:    conf.Topic,
		Tag:      conf.Tag,
		producer: producer,
	}
	if err := publisher.producer.Start(); err != nil {
		logs.Error("start mq producer fail, cfg: %v", conf)
		return nil, err
	}
	return publisher, nil
}

func (p *DatasetJobPublisher) Send(ctx context.Context, msg *entity.JobRunMessage, opts ...imq.MessageOpt) error {
	body, err := sonic.Marshal(msg)
	if err != nil {
		return errors.WithMessage(err, "marshal message")
	}
	opt := &imq.MessageOption{}
	for _, o := range opts {
		o(opt)
	}
	mqMsg := &mq.Message{
		Topic:         p.Topic,
		Body:          body,
		Tag:           p.Tag,
		PartitionKey:  opt.Key,
		DeferDuration: opt.RetryInterval,
	}
	result, err := p.producer.Send(ctx, mqMsg)
	if err != nil {
		return errno.InternalErr(err, "send msg failed, body=%v, pk=%v", body)
	}
	logs.CtxInfo(
		ctx,
		"rmq message sent, msg_id=%s, offset=%s, topic=%s, tag='%s'",
		result.MessageID,
		result.Offset,
		p.Topic,
		p.Tag,
	)
	return nil
}
