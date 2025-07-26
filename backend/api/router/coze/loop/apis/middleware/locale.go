// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-loop/backend/pkg/consts"
	"github.com/coze-dev/coze-loop/backend/pkg/contexts"
)

var supportedLocales = []string{
	consts.LocaleZhCN,
	consts.LocalEnUS,
}

func LocaleMW() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(contexts.WithLocale(ctx, parseLocale(c)))
	}
}

func parseLocale(c *app.RequestContext) string {
	if locale := string(c.Cookie(consts.CookieLanguageKey)); isValidLocale(locale) {
		return locale
	}
	return consts.LocaleDefault
}

func isValidLocale(locale string) bool {
	if len(locale) == 0 {
		return false
	}
	for _, supported := range supportedLocales {
		if strings.EqualFold(locale, supported) {
			return true
		}
	}
	return false
}
