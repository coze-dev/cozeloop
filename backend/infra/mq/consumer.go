// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mq

import (
	"context"
)

//go:generate mockgen -destination=mocks/consumer.go -package=mocks . IConsumer
type IConsumer interface {
	Start() error
	Close() error
	RegisterHandler(IConsumerHandler)
}

type IConsumerHandler interface {
	HandleMessage(context.Context, *MessageExt) error
}
