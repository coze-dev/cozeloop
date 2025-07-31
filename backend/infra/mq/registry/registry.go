// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/infra/mq"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
	"github.com/coze-dev/coze-loop/backend/pkg/json"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/goroutine"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
)

type defaultConsumerRegistry struct {
	factory mq.IFactory
	workers []mq.IConsumerWorker
}

func NewConsumerRegistry(factory mq.IFactory) mq.ConsumerRegistry {
	return &defaultConsumerRegistry{factory: factory}
}

func (d *defaultConsumerRegistry) Register(worker []mq.IConsumerWorker) mq.ConsumerRegistry {
	d.workers = append(d.workers, worker...)
	return d
}

func (d *defaultConsumerRegistry) StartAll(ctx context.Context) error {
	for _, worker := range d.workers {
		cfg, err := worker.ConsumerCfg(ctx)
		if err != nil {
			return err
		}

		consumer, err := d.factory.NewConsumer(ptr.From(cfg))
		if err != nil {
			return errorx.Wrapf(err, "NewConsumer fail, cfg: %v", json.Jsonify(cfg))
		}

		consumer.RegisterHandler(newSafeConsumerWrapper(worker))
		if err := consumer.Start(); err != nil {
			return errorx.Wrapf(err, "StartConsumer fail, cfg: %v", json.Jsonify(cfg))
		}
	}
	return nil
}

type safeConsumerHandlerDecorator struct {
	handler mq.IConsumerHandler
}

func (s *safeConsumerHandlerDecorator) HandleMessage(ctx context.Context, msg *mq.MessageExt) error {
	defer goroutine.Recovery(ctx)
	return s.handler.HandleMessage(ctx, msg)
}

func newSafeConsumerWrapper(h mq.IConsumerHandler) mq.IConsumerHandler {
	return &safeConsumerHandlerDecorator{handler: h}
}
