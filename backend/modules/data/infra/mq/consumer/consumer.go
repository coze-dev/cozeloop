// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package consumer

import (
	"github.com/coze-dev/cozeloop/backend/infra/mq"
	"github.com/coze-dev/cozeloop/backend/modules/data/application"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/consts"
	"github.com/coze-dev/cozeloop/backend/pkg/conf"
)

func NewConsumerWorkers(
	cfactory conf.IConfigLoaderFactory,
	handler application.IJobRunMsgHandler,
) ([]mq.IConsumerWorker, error) {
	loader, err := cfactory.NewConfigLoader(consts.DataConfigFileName)
	if err != nil {
		return nil, err
	}
	return []mq.IConsumerWorker{
		newDatasetJobConsumer(handler, loader),
	}, nil
}
