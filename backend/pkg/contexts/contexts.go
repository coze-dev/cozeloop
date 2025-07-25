// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package contexts

import (
	"context"

	"github.com/coze-dev/cozeloop/backend/pkg/consts"
)

type ctxLocaleKeyType struct{}

var ctxLocaleKey = ctxLocaleKeyType{}

func WithLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, ctxLocaleKey, locale)
}

func CtxLocale(ctx context.Context) string {
	locale, ok := ctx.Value(ctxLocaleKey).(string)
	if !ok || len(locale) == 0 {
		return consts.LocaleDefault
	}
	return locale
}
