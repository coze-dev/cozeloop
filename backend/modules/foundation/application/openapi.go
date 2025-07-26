// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"bytes"
	"context"
	"errors"
	"mime"
	"mime/multipart"

	"github.com/samber/lo"

	"github.com/coze-dev/coze-loop/backend/infra/fileserver"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/base"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/foundation/openapi"
	"github.com/coze-dev/coze-loop/backend/modules/foundation/domain/component/rpc"
	"github.com/coze-dev/coze-loop/backend/modules/foundation/domain/file/service"
	"github.com/coze-dev/coze-loop/backend/modules/foundation/pkg/errno"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
	"github.com/coze-dev/coze-loop/backend/pkg/logs"
)

type FoundationOpenAPIApplicationImpl struct {
	auth        rpc.IAuthProvider
	fileService service.FileService
}

func NewFoundationOpenAPIApplication(objectStorage fileserver.BatchObjectStorage, authService rpc.IAuthProvider) openapi.FoundationOpenAPIService {
	return &FoundationOpenAPIApplicationImpl{
		auth:        authService,
		fileService: service.NewFileService(objectStorage),
	}
}

func (f FoundationOpenAPIApplicationImpl) UploadLoopFile(ctx context.Context, req *openapi.UploadLoopFileRequest) (r *openapi.UploadLoopFileResponse, err error) {
	if req == nil || req.ContentType == "" || len(req.Body) == 0 {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode)
	}

	form, err := parseMultipartFormData(ctx, req.ContentType, req.Body)
	if err != nil {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode)
	}
	if form == nil || len(form.File["file"]) == 0 {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode)
	}
	fileHeader := form.File["file"][0]
	spaceID := ""
	if form.Value != nil && len(form.Value["workspace_id"]) > 0 {
		spaceID = form.Value["workspace_id"][0]
	}
	if spaceID == "" {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode)
	}

	err = f.auth.CheckWorkspacePermission(ctx, rpc.AuthActionOpenAPIFileUpload, spaceID)
	if err != nil {
		return nil, err
	}

	key, err := f.fileService.UploadLoopFile(ctx, fileHeader, spaceID)
	if err != nil {
		return nil, err
	}

	return &openapi.UploadLoopFileResponse{
		Data: &openapi.FileData{
			Bytes:    lo.ToPtr(int64(len(req.Body))),
			FileName: lo.ToPtr(key),
		},
		BaseResp: base.NewBaseResp(),
	}, nil
}

func parseMultipartFormData(ctx context.Context, httpContentType string, data []byte) (*multipart.Form, error) {
	_, params, err := mime.ParseMediaType(httpContentType)
	if err != nil {
		logs.CtxError(ctx, "parse ContentType failed, err: %v, contentType: %v", err, httpContentType)
		return nil, errorx.NewByCode(errno.CommonInternalErrorCode)
	}
	br := bytes.NewReader(data)
	// boundary is in contentType
	mr := multipart.NewReader(br, params["boundary"])
	form, err := mr.ReadForm(int64(1 * 1024 * 1024 * 1024))
	if errors.Is(err, multipart.ErrMessageTooLarge) {
		return nil, errorx.NewByCode(errno.FileSizeExceedLimitCode)
	} else if err != nil {
		logs.CtxError(ctx, "read form failed, err: %v", err)
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode)
	}
	return form, nil
}
