// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package errors

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

const (
	ServiceInternalErrorCode int32 = 1
	DefaultErrorMsg                = "Service Internal Error"
)

func KiteXSvrCompatMW() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			err = next(ctx, req, resp)
			if err == nil {
				return nil
			}

			if bizErr, ok := kerrors.FromBizStatusError(err); ok {
				ri := rpcinfo.GetRPCInfo(ctx)
				if setter, ok := ri.Invocation().(rpcinfo.InvocationSetter); ok {
					setter.SetBizStatusErr(bizErr)
					return nil
				}
			}

			return err
		}
	}
}

func KiteXSvrErrorWrapMW() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			err = next(ctx, req, resp)
			if err == nil {
				return nil
			}

			if kerrors.IsKitexError(err) && !errors.Is(err, kerrors.ErrBiz) {
				return err
			}

			be := getBizStatusError(err)

			ri := rpcinfo.GetRPCInfo(ctx)
			if setter, ok := ri.Invocation().(rpcinfo.InvocationSetter); ok {
				setter.SetBizStatusErr(be)
				return nil
			}

			return be
		}
	}
}

func getBizStatusError(err error) kerrors.BizStatusErrorIface {
	var detailedErr *kerrors.DetailedError
	if errors.As(err, &detailedErr) {
		unknownErr := detailedErr.Unwrap()
		return kerrors.NewBizStatusError(ServiceInternalErrorCode, fmt.Sprintf("%s:%v", DefaultErrorMsg, unknownErr))
	}
	return kerrors.NewBizStatusError(ServiceInternalErrorCode, fmt.Sprintf("%s:%v", DefaultErrorMsg, err))
}
