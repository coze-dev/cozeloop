// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package goroutine

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	t.Run("recover from panic", func(t *testing.T) {
		ctx := context.Background()
		var err error

		func() {
			defer Recover(ctx, &err)
			panic("test panic")
		}()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "panic occurred")
	})

	t.Run("recover with existing error", func(t *testing.T) {
		ctx := context.Background()
		originalErr := errors.New("original error")
		err := originalErr

		func() {
			defer Recover(ctx, &err)
			panic("test panic")
		}()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "panic occured")
		assert.Contains(t, err.Error(), originalErr.Error())
	})
}

func TestRecovery(t *testing.T) {
	t.Run("recover from panic", func(t *testing.T) {
		ctx := context.Background()

		func() {
			defer Recovery(ctx)
			panic("test panic")
		}()
		// No assertion needed as Recovery only logs the error
	})
}
