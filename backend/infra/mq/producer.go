// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mq

import (
	"context"
)

//go:generate mockgen -destination=mocks/producer.go -package=mocks . IProducer
type IProducer interface {
	Start() error
	Close() error
	Send(ctx context.Context, message *Message) (SendResponse, error)
	SendBatch(ctx context.Context, messages []*Message) (SendResponse, error)
	SendAsync(ctx context.Context, callback AsyncSendCallback, message *Message) error
}

type AsyncSendCallback func(ctx context.Context, sendResponse SendResponse, err error)

type SendResponse struct {
	MessageID string
	Offset    int64
}
