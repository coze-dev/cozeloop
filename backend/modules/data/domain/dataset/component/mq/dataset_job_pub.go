// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mq

import (
	"context"
	"time"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
)

//go:generate mockgen -destination=mocks/dataset_job_pub.go -package=mocks . IDatasetJobPublisher
type IDatasetJobPublisher interface {
	Send(ctx context.Context, msg *entity.JobRunMessage, opts ...MessageOpt) error
}

type MessageOption struct {
	Key           string
	RetryInterval time.Duration
}

type MessageOpt func(*MessageOption)

func WithDelayTimeLevel(retryInterval time.Duration) MessageOpt {
	return func(m *MessageOption) {
		m.RetryInterval = retryInterval
	}
}

func WithKey(key string) MessageOpt {
	return func(m *MessageOption) {
		m.Key = key
	}
}
