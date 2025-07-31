// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package code

import (
	"github.com/coze-dev/coze-loop/backend/pkg/errorx/internal"
)

type RegisterOptionFn = internal.RegisterOption

func WithAffectStability(affectStability bool) RegisterOptionFn {
	return internal.WithAffectStability(affectStability)
}

func Register(code int32, msg string, opts ...RegisterOptionFn) {
	internal.Register(code, msg, opts...)
}

func SetDefaultErrorCode(code int32) {
	internal.SetDefaultErrorCode(code)
}
