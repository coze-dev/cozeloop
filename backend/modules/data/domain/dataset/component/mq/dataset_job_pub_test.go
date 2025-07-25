// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mq

import (
	"testing"
	"time"
)

func TestWithDelayTimeLevel(t *testing.T) {
	interval := 5 * time.Second
	opt := WithDelayTimeLevel(interval)
	option := &MessageOption{}
	opt(option)

	if option.RetryInterval != interval {
		t.Errorf("Expected RetryInterval to be %v, but got %v", interval, option.RetryInterval)
	}
}

func TestWithKey(t *testing.T) {
	key := "test-key"
	opt := WithKey(key)
	option := &MessageOption{}
	opt(option)

	if option.Key != key {
		t.Errorf("Expected Key to be %s, but got %s", key, option.Key)
	}
}
