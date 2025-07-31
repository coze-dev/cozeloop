// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package goroutineutil

import (
	"context"
	"runtime"

	"github.com/coze-dev/coze-loop/backend/pkg/logs"
)

func GoWithDeferFunc(ctx context.Context, f func()) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				logs.CtxError(ctx, "goroutine panic: %s: %s", e, buf)
			}
		}()
		f()
	}()
}

func GoWithDefaultRecovery(ctx context.Context, f func()) {
	GoWithDeferFunc(ctx, f)
}
