// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/coze-loop/backend/pkg/consts"
	"github.com/coze-dev/coze-loop/backend/pkg/contexts"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
)

func TestLocaleMW(t *testing.T) {
	tests := []struct {
		name           string
		cookieValue    string
		expectedLocale string
	}{
		{
			name:           "Valid zh-CN locale from cookie",
			cookieValue:    "zh-CN",
			expectedLocale: "zh-CN",
		},
		{
			name:           "Valid en-US locale from cookie",
			cookieValue:    "en-US",
			expectedLocale: "en-US",
		},
		{
			name:           "Case insensitive zh-CN locale",
			cookieValue:    "ZH-CN",
			expectedLocale: "ZH-CN",
		},
		{
			name:           "Case insensitive en-US locale",
			cookieValue:    "EN-US",
			expectedLocale: "EN-US",
		},
		{
			name:           "Empty cookie value",
			cookieValue:    "",
			expectedLocale: consts.LocaleDefault,
		},
		{
			name:           "Invalid locale from cookie",
			cookieValue:    "fr-FR",
			expectedLocale: consts.LocaleDefault,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			c := &app.RequestContext{}
			c.Request = ptr.From(protocol.AcquireRequest())
			c.Response = ptr.From(protocol.AcquireResponse())

			if tt.cookieValue != "" {
				c.Request.Header.SetCookie(consts.CookieLanguageKey, tt.cookieValue)
			}

			middleware := LocaleMW()

			next := func() app.HandlerFunc {
				return func(ctx context.Context, c *app.RequestContext) {
					assert.Equal(t, tt.expectedLocale, contexts.CtxLocale(ctx))
				}
			}
			c.SetHandlers(append(c.Handlers(), next()))

			middleware(ctx, c)

			assert.NotPanics(t, func() {
				middleware(ctx, c)
			})
		})
	}
}

func TestParseLocale(t *testing.T) {
	tests := []struct {
		name           string
		cookieValue    string
		expectedLocale string
	}{
		{
			name:           "Valid zh-CN locale",
			cookieValue:    "zh-CN",
			expectedLocale: "zh-CN",
		},
		{
			name:           "Valid en-US locale",
			cookieValue:    "en-US",
			expectedLocale: "en-US",
		},
		{
			name:           "Case insensitive zh-CN locale",
			cookieValue:    "ZH-CN",
			expectedLocale: "ZH-CN",
		},
		{
			name:           "Case insensitive en-US locale",
			cookieValue:    "EN-US",
			expectedLocale: "EN-US",
		},
		{
			name:           "Empty cookie value",
			cookieValue:    "",
			expectedLocale: consts.LocaleDefault,
		},
		{
			name:           "Invalid locale",
			cookieValue:    "fr-FR",
			expectedLocale: consts.LocaleDefault,
		},
		{
			name:           "Partial match should fail",
			cookieValue:    "zh",
			expectedLocale: consts.LocaleDefault,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &app.RequestContext{}
			c.Request = ptr.From(protocol.AcquireRequest())
			c.Response = ptr.From(protocol.AcquireResponse())

			if tt.cookieValue != "" {
				c.Request.Header.SetCookie(consts.CookieLanguageKey, tt.cookieValue)
			}

			result := parseLocale(c)

			assert.Equal(t, tt.expectedLocale, result)
		})
	}
}

func TestIsValidLocale(t *testing.T) {
	tests := []struct {
		name     string
		locale   string
		expected bool
	}{
		{
			name:     "Valid zh-CN locale",
			locale:   "zh-CN",
			expected: true,
		},
		{
			name:     "Valid en-US locale",
			locale:   "en-US",
			expected: true,
		},
		{
			name:     "Case insensitive zh-CN locale",
			locale:   "ZH-CN",
			expected: true,
		},
		{
			name:     "Case insensitive en-US locale",
			locale:   "EN-US",
			expected: true,
		},
		{
			name:     "Empty string",
			locale:   "",
			expected: false,
		},
		{
			name:     "Invalid locale fr-FR",
			locale:   "fr-FR",
			expected: false,
		},
		{
			name:     "Partial match zh",
			locale:   "zh",
			expected: false,
		},
		{
			name:     "Mixed case zh-cn",
			locale:   "zh-cn",
			expected: true,
		},
		{
			name:     "Mixed case en-us",
			locale:   "en-us",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidLocale(tt.locale)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSupportedLocales(t *testing.T) {
	expectedLocales := []string{
		consts.LocaleZhCN,
		consts.LocalEnUS,
	}
	assert.Equal(t, expectedLocales, supportedLocales)
	assert.Equal(t, 2, len(supportedLocales))
}
