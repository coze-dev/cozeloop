// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package backoff

import (
	"context"
	"time"

	"github.com/cenk/backoff"
)

const (
	defaultRetryInterval = time.Millisecond * 100

	one     = time.Second * 1
	three   = time.Second * 3
	five    = time.Second * 5
	ten     = time.Second * 10
	OneMin  = time.Minute
	FiveMin = time.Minute * 5
	OneHour = time.Hour
)

func RetryOneSecond(ctx context.Context, fn func() error) error {
	return RetryWithElapsedTime(ctx, one, fn)
}

func RetryThreeSeconds(ctx context.Context, fn func() error) error {
	return RetryWithElapsedTime(ctx, three, fn)
}

func RetryFiveSeconds(ctx context.Context, fn func() error) error {
	return RetryWithElapsedTime(ctx, five, fn)
}

func RetryOneMin(ctx context.Context, fn func() error) error {
	return RetryWithElapsedTime(ctx, OneMin, fn)
}

func RetryFiveMin(ctx context.Context, fn func() error) error {
	return RetryWithElapsedTime(ctx, FiveMin, fn)
}

func RetryOneHour(ctx context.Context, fn func() error) error {
	return RetryWithElapsedTime(ctx, OneHour, fn)
}

func RetryTenSeconds(ctx context.Context, fn func() error) error {
	return RetryWithElapsedTime(ctx, ten, fn)
}

func RetryWithElapsedTime(ctx context.Context, maxElapsedTime time.Duration, fn func() error) error {
	policy := backoff.NewExponentialBackOff()
	policy.InitialInterval = defaultRetryInterval
	policy.MaxElapsedTime = maxElapsedTime

	return backoffFn(ctx, fn, policy)
}

func backoffFn(ctx context.Context, fn func() error, policy *backoff.ExponentialBackOff) error {
	ctxWithCancel, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	backoffCtx := backoff.WithContext(policy, ctxWithCancel)

	return backoff.Retry(fn, backoffCtx)
}

func RetryWithMaxTimes(ctx context.Context, max int, fn func() error) error {
	return backoff.Retry(fn, backoff.WithMaxRetries(&backoff.ZeroBackOff{}, uint64(max)))
}
