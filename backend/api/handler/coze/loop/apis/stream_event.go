// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package apis

import (
	"context"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/hertz-contrib/sse"

	middleware "github.com/coze-dev/coze-loop/backend/infra/middleware/errors"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/js_conv"
	"github.com/coze-dev/coze-loop/backend/pkg/logs"
)

func publishDataEvent(ctx context.Context, s *sse.Stream, data any) error {
	if data == nil {
		return nil
	}
	bytes, err := js_conv.GetMarshaler()(data)
	if err != nil {
		logs.CtxError(ctx, "marshal data packet error: %s", err.Error())
		return err
	}
	event := &sse.Event{
		Event: "data",
		Data:  bytes,
	}
	return s.Publish(event)
}

func publishErrEvent(ctx context.Context, s *sse.Stream, err error) {
	if err == nil {
		return
	}
	var bErr *bizErr
	if statusErr, ok := kerrors.FromBizStatusError(err); ok {
		bErr = &bizErr{
			Code:     statusErr.BizStatusCode(),
			Msg:      statusErr.BizMessage(),
			BizExtra: statusErr.BizExtra(),
		}
	} else {
		bErr = &bizErr{
			Code: middleware.ServiceInternalErrorCode,
			Msg:  middleware.DefaultErrorMsg,
		}
	}
	var errData []byte
	errData, err = js_conv.GetMarshaler()(bErr)
	if err != nil {
		logs.CtxError(ctx, "marshal err packet error: %s", err.Error())
	}
	publishErr := s.Publish(&sse.Event{
		Event: "error",
		Data:  errData,
	})
	if publishErr != nil {
		logs.CtxError(ctx, "publish event error: %s", publishErr.Error())
		return
	}
}

type bizErr struct {
	Code     int32             `json:"code"`
	Msg      string            `json:"msg"`
	BizExtra map[string]string `json:"biz_extra"`
}
