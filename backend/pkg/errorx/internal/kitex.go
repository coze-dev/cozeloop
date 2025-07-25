// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/kerrors"

	"github.com/coze-dev/cozeloop/backend/pkg/lang/conv"
)

var _ KiteXBizStatusError = &statusError{}

const (
	BizExtraKeyAffectStability = "biz_err_affect_stability"
	BizExtraKeyCustomExtra     = "biz_err_custom_extra"
)

// KiteXBizStatusError satisfy the interface requirement of Kitex biz exception handling
// Kitex-biz exception usage doc: https://www.cloudwego.io/zh/docs/kitex/tutorials/basic-feature/bizstatuserr/
type KiteXBizStatusError = kerrors.BizStatusErrorIface

// BizStatusCode implements kerrors.BizStatusErrorIface, support Kitex biz exception handling.
func (w *statusError) BizStatusCode() int32 {
	return w.statusCode
}

// BizMessage implements kerrors.BizStatusErrorIface, support Kitex biz exception handling.
func (w *statusError) BizMessage() string {
	return w.message
}

// BizExtra implements kerrors.BizStatusErrorIface, support Kitex biz exception handling.
func (w *statusError) BizExtra() map[string]string {
	bizExtra := map[string]string{}
	if len(w.ext.Extra) > 0 {
		extraStr, _ := sonic.MarshalString(w.ext.Extra)
		bizExtra[BizExtraKeyCustomExtra] = extraStr
	}

	affectStability := "0"
	if w.ext.IsAffectStability {
		affectStability = "1"
	}
	bizExtra[BizExtraKeyAffectStability] = affectStability

	return bizExtra
}

// handlerErr withCode 抛出的错误类型：
// Error() 打印堆栈；kerrors.FromBizStatusError 后转为 *kitexBizErrWrapper 类型；errors.Is 根据 Code 比较
type handlerErr struct {
	kitexBizErrWrapper *kitexBizErrWrapper
}

func (h *handlerErr) Error() string {
	return h.kitexBizErrWrapper.withStatus.Error()
}

func (h *handlerErr) Format(f fmt.State, c rune) {
	_, _ = f.Write(conv.UnsafeStringToBytes(h.kitexBizErrWrapper.withStatus.Error()))
}

// Unwrap supports go errors.Unwrap().
func (h *handlerErr) Unwrap() error {
	return h.kitexBizErrWrapper.Unwrap()
}

// Is supports go errors.Is().
func (h *handlerErr) Is(target error) bool {
	return h.kitexBizErrWrapper.Is(target)
}

// As supports go errors.As().
func (h *handlerErr) As(target interface{}) bool {
	return errors.As(h.kitexBizErrWrapper, target)
}

func (h *handlerErr) StackTrace() string {
	return h.kitexBizErrWrapper.stack
}

// kitexBizErrWrapper 为了实现kitex中间件可打印堆栈，但堆栈不流转出服务
type kitexBizErrWrapper struct {
	*withStatus
}

func (h *kitexBizErrWrapper) BizStatusCode() int32 {
	return h.status.BizStatusCode()
}

func (h *kitexBizErrWrapper) BizMessage() string {
	return h.status.BizMessage()
}

func (h *kitexBizErrWrapper) BizExtra() map[string]string {
	return h.status.BizExtra()
}

func (h *kitexBizErrWrapper) Error() string {
	return h.status.Error()
}

func (h *kitexBizErrWrapper) Format(f fmt.State, c rune) {
	if c == 'v' && f.Flag('+') {
		_, _ = f.Write(conv.UnsafeStringToBytes(h.withStatus.Error()))
	} else {
		_, _ = f.Write(conv.UnsafeStringToBytes(h.status.Error()))
	}
}
