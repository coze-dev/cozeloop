// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package backoff

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cenk/backoff"
	"github.com/stretchr/testify/assert"
)

func Test_Backoff(t *testing.T) {
	ctx := context.Background()
	fn := func() error { return fmt.Errorf("mock err") }
	p := backoff.NewExponentialBackOff()
	p.InitialInterval = defaultRetryInterval
	p.MaxElapsedTime = defaultRetryInterval

	// Since we can't mock without mockey, we'll test the actual functions
	// These should fail as expected since fn always returns an error
	assert.NotNil(t, RetryOneSecond(ctx, fn))
	assert.NotNil(t, RetryThreeSeconds(ctx, fn))
	assert.NotNil(t, RetryFiveSeconds(ctx, fn))
	assert.NotNil(t, RetryTenSeconds(ctx, fn))
}

func Test_backoff(t *testing.T) {
	ctx := context.Background()

	t.Run("test success", func(t *testing.T) {
		fn := func() error { return nil }
		assert.Nil(t, RetryWithElapsedTime(ctx, time.Second, fn))
	})

	t.Run("test ctx cancel", func(t *testing.T) {
		cc, cancelFn := context.WithCancel(ctx)
		cnt := 0
		fn := func() error {
			cancelFn()
			cnt++
			return fmt.Errorf("mock err")
		}

		start := time.Now()
		assert.NotNil(t, RetryWithElapsedTime(cc, time.Second, fn))
		assert.Equal(t, 1, cnt)
		assert.True(t, time.Since(start) < time.Second)
	})
}

func TestRetryWithMaxTimes(t *testing.T) {
	ctx := context.Background()

	t.Run("test success", func(t *testing.T) {
		var count int
		err := RetryWithMaxTimes(ctx, 3, func() error {
			count++
			return fmt.Errorf("error")
		})
		assert.NotNil(t, err)
		assert.Equal(t, 4, count)
	})
}
