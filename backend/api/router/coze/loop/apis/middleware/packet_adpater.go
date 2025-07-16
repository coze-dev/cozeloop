// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/mitchellh/mapstructure"

	"github.com/coze-dev/cozeloop/backend/infra/i18n"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/consts"
	"github.com/coze-dev/cozeloop/backend/pkg/json"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

const (
	headerContentTypeKey  = "content-type"
	headerContentTypeJSON = "application/json"
	headerContentTypeSSE  = "text/event-stream"

	jsonBodyCodeKey = "code"
	jsonBodyMsgKey  = "msg"

	camelCaseBodyMapKeyBaseResp  = "BaseResp"
	baseRespExtraAffectStableKey = "biz_err_affect_stability"
	affectStableValue            = "1"
)

func PacketAdapterMW(translater i18n.ITranslater) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)

		if ep := parseErrPacket(ctx, c); ep != nil {
			c.JSON(http.StatusOK, ep.localizeMessage(ctx, string(c.Cookie(consts.CookieLanguageKey)), translater))
			return
		}

		adapter := &respPacketAdapter{ctx: c}
		packet, err := adapter.getRespPacket(ctx)
		if err != nil {
			logs.CtxWarn(ctx, "parse hertz resp packet fail, err: %v", err)
			return
		}

		if err := adapter.wrapPacket(ctx, packet); err != nil {
			logs.CtxWarn(ctx, "wrap packet fail, err: %v", err)
			return
		}
	}
}

func parseErrPacket(ctx context.Context, c *app.RequestContext) *errPacket {
	if len(c.Errors) == 0 {
		return nil
	}

	logs.CtxError(ctx, "found hertz resp err: %v", c.Errors.String())

	for i := len(c.Errors) - 1; i >= 0; i-- {
		berr, ok := kerrors.FromBizStatusError(c.Errors[i].Err)
		if ok {
			return &errPacket{
				Code:    berr.BizStatusCode(),
				Message: berr.BizMessage(),
			}
		}
	}

	return &errPacket{
		Code:    errno.CommonInternalErrorCode,
		Message: "Service Internal Error",
	}
}

type respPacket struct {
	bodyMap         map[string]any
	baseResp        *baseResp
	baseRespBodyKey string
}

func (r *respPacket) parseBaseResp(ctx context.Context) *respPacket {
	if r.bodyMap == nil {
		return r
	}

	for _, key := range []string{camelCaseBodyMapKeyBaseResp} {
		if val, exists := r.bodyMap[key]; exists {
			br := &baseResp{}
			if err := mapstructure.Decode(val, br); err != nil {
				logs.CtxError(ctx, "decode map to base_resp fail, raw: %v, err: %v", val, err)
			} else {
				r.baseResp = br
				r.baseRespBodyKey = key
			}
		}
	}

	return r
}

type baseResp struct {
	StatusMessage string
	StatusCode    int32
	Extra         map[string]string
}

func (b baseResp) IsSuccessStatus() bool {
	return b.StatusCode == 0
}

func (b baseResp) AffectStability() bool {
	val, ok := b.Extra[baseRespExtraAffectStableKey]
	if !ok {
		return true
	}
	return val == affectStableValue
}

type errPacket struct {
	Code    int32  `json:"code"`
	Message string `json:"msg"`
}

func (e *errPacket) localizeMessage(ctx context.Context, locale string, translater i18n.ITranslater) *errPacket {
	if translater == nil || len(locale) == 0 {
		return e
	}
	msg := translater.MustTranslate(ctx, strconv.Itoa(int(e.Code)), locale)
	if len(msg) > 0 {
		e.Message = msg
	}
	return e
}

type respPacketAdapter struct {
	ctx        *app.RequestContext
	translater i18n.ITranslater
}

func (r *respPacketAdapter) getRespPacket(ctx context.Context) (*respPacket, error) {
	if r.ctx == nil || r.ctx.GetResponse() == nil {
		return nil, fmt.Errorf("hertz req with invalid response")
	}

	ct := strings.ToLower(r.ctx.GetResponse().Header.Get(headerContentTypeKey))
	if !strings.Contains(ct, headerContentTypeJSON) && !strings.Contains(ct, headerContentTypeSSE) {
		return nil, fmt.Errorf("response content type is not json or event-stream")
	}

	bodyMap := make(map[string]any)
	if bodyBytes := r.ctx.GetResponse().Body(); len(bodyBytes) > 0 {
		if err := sonic.Unmarshal(bodyBytes, &bodyMap); err != nil {
			return nil, fmt.Errorf("parse hertz resp map body fail: %w", err)
		}
	}

	rp := &respPacket{bodyMap: bodyMap}
	return rp.parseBaseResp(ctx), nil
}

func (r *respPacketAdapter) wrapPacket(ctx context.Context, packet *respPacket) error {
	bodyMap := packet.bodyMap
	if _, ok := bodyMap[jsonBodyCodeKey]; ok { // 若返回 body 中包含 code 则不处理
		return nil
	}

	wrapFn := func(bodyMap map[string]any, code int32, msg string) map[string]any {
		if _, ok1 := bodyMap[jsonBodyCodeKey]; !ok1 {
			bodyMap[jsonBodyCodeKey] = code
			if _, ok2 := bodyMap[jsonBodyMsgKey]; !ok2 {
				bodyMap[jsonBodyMsgKey] = msg
			}
		}
		return bodyMap
	}
	setPacketFn := func(body any) error {
		bytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("json marshal error fail: %w", err)
		}
		r.ctx.Response.SetBody(bytes)
		logs.CtxDebug(ctx, "repack body with err packet: "+string(bytes))
		return nil
	}

	br := packet.baseResp
	if br == nil {
		return setPacketFn(wrapFn(bodyMap, 0, ""))
	}

	if br.IsSuccessStatus() {
		delete(bodyMap, packet.baseRespBodyKey)
		return setPacketFn(wrapFn(bodyMap, 0, ""))
	}

	ep := &errPacket{
		Code:    br.StatusCode,
		Message: br.StatusMessage,
	}
	return setPacketFn(ep.localizeMessage(ctx, string(r.ctx.Cookie(consts.CookieLanguageKey)), r.translater))
}
