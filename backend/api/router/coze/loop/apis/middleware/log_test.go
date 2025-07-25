// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
)

func TestAccessLogMW(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedLog    string
		expectedHeader string
	}{
		{
			name:           "success response",
			statusCode:     http.StatusOK,
			responseBody:   `{"data":"test"}`,
			expectedHeader: "X-Tt-Logid",
		},
		{
			name:           "bad request",
			statusCode:     http.StatusBadRequest,
			responseBody:   `{"code":400,"msg":"bad request"}`,
			expectedHeader: "X-Tt-Logid",
		},
		{
			name:           "internal server error",
			statusCode:     http.StatusInternalServerError,
			responseBody:   `{"code":500,"msg":"internal error"}`,
			expectedHeader: "X-Tt-Logid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &app.RequestContext{}
			ctx.Request = ptr.From(protocol.AcquireRequest())
			ctx.Response = ptr.From(protocol.AcquireResponse())

			ctx.Request.SetRequestURI("/test")
			ctx.Request.Header.SetMethod(consts.MethodGet)
			ctx.Request.SetBodyString("")

			ctx.Response.SetStatusCode(tt.statusCode)
			ctx.Response.SetBodyString(tt.responseBody)

			c, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			assert.NotPanics(t, func() {
				AccessLogMW()(c, ctx)
			})
		})
	}
}
