// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package goroutine

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/samber/lo"

	"github.com/coze-dev/coze-loop/backend/pkg/logs"
)

func Recover(ctx context.Context, errPtr *error) {
	e := recover()
	if e == nil {
		return
	}

	var tmpErr error
	if errPtr != nil && *errPtr != nil {
		tmpErr = fmt.Errorf("panic occured, originErr=%v, reason=%v", *errPtr, e)
	} else {
		tmpErr = fmt.Errorf("panic occurred, reason=%v", e)
	}

	if errPtr != nil {
		*errPtr = tmpErr
	}

	ctx = lo.Ternary(ctx == nil, context.Background(), ctx)

	logs.CtxError(ctx, "[catch panic] err = %v \n stacktrace:\n%s", fmt.Errorf("%v", e), debug.Stack())
}

func Recovery(ctx context.Context) {
	e := recover()
	if e == nil {
		return
	}

	ctx = lo.Ternary(ctx == nil, context.Background(), ctx)

	logs.CtxError(ctx, "[catch panic] err = %v \n stacktrace:\n%s", fmt.Errorf("%v", e), debug.Stack())
}
