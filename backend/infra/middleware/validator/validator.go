// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"context"
	"fmt"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"

	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
)

func KiteXValidatorMW(next endpoint.Endpoint) endpoint.Endpoint {
	validate := func(req any) error {
		withArg, ok := req.(interface{ GetFirstArgument() any })
		if !ok || withArg == nil {
			return nil
		}

		arg, ok := withArg.GetFirstArgument().(interface{ IsValid() error })
		if !ok || arg == nil {
			return nil
		}

		return arg.IsValid()
	}

	return func(ctx context.Context, req, resp any) (err error) {
		const unknown = "unknown"
		if err := validate(req); err != nil {
			method := unknown
			if info := rpcinfo.GetRPCInfo(ctx); info != nil && info.To() != nil {
				method = info.To().Method()
			}
			return errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg(fmt.Sprintf("validate request fail, method=%s, err=%s", method, err.Error())))
		}

		return next(ctx, req, resp)
	}
}
