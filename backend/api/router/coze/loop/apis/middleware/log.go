// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/cozeloop/backend/pkg/json"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/conv"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

func AccessLogMW() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		const headerKeyLogID = "X-Log-ID"

		start := time.Now()
		logID := logs.NewLogID()
		c = logs.SetLogID(c, logID)

		ctx.Next(c)

		ctx.Response.Header.Set(headerKeyLogID, logID)

		status := ctx.Response.StatusCode()
		path := conv.UnsafeBytesToString(ctx.Request.URI().PathOriginal())
		latency := time.Since(start)
		method := conv.UnsafeBytesToString(ctx.Request.Header.Method())
		clientIP := ctx.ClientIP()
		baseLog := fmt.Sprintf("| %d | %v | %s | %s | %v", status, latency, clientIP, method, path)

		ep := &errPacket{}
		_ = json.Unmarshal(ctx.GetResponse().Body(), ep)

		switch {
		case status >= http.StatusInternalServerError || ep.Code != 0:
			logs.CtxError(c, "%s | %s", baseLog, ctx.GetResponse().Body())
		case status >= http.StatusBadRequest:
			logs.CtxWarn(c, "%s | %s", baseLog, ctx.GetResponse().Body())
		default:
			urlQuery := ctx.Request.URI().QueryString()
			reqBody := conv.UnsafeBytesToString(ctx.Request.Body())
			respBody := conv.UnsafeBytesToString(ctx.Response.Body())

			logs.CtxDebug(c, "%s \nquery : %s \nreq : %s \nresp: %s",
				baseLog, urlQuery, reqBody, respBody)
		}
	}
}
