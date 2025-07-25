// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	errors2 "github.com/pkg/errors"

	"github.com/coze-dev/cozeloop/backend/infra/fileserver"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/file"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/localos"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

//go:generate mockgen -destination=mocks/file_service.go -package=mocks . FileService
type FileService interface {
	UploadLoopFile(ctx context.Context, fileHeader *multipart.FileHeader, spaceID string) (key string, err error)
	SignUploadFile(ctx context.Context, req *file.SignUploadFileRequest) (uris []string, heads []*file.SignHead, err error)
	SignDownLoadFile(ctx context.Context, req *file.SignDownloadFileRequest) (uris []string, err error)
}

type fileService struct {
	client fileserver.BatchObjectStorage
}

func NewFileService(objectStorage fileserver.BatchObjectStorage) FileService {
	return &fileService{
		client: objectStorage,
	}
}

func (fs fileService) UploadLoopFile(ctx context.Context, fileHeader *multipart.FileHeader, spaceID string) (key string, err error) {
	fileName := fileHeader.Filename
	if fileName == "" {
		return "", errorx.NewByCode(errno.CommonInvalidParamCode)
	}
	fileName = filepath.Join(spaceID, "/", fileName)

	f, err := fileHeader.Open()
	if err != nil {
		logs.CtxError(ctx, "open file failed, err: %v", err)
		return "", errorx.NewByCode(errno.CommonInternalErrorCode)
	}
	defer func(f multipart.File) {
		if f == nil {
			return
		}
		err := f.Close()
		if err != nil {
			logs.CtxError(ctx, "close file failed, err: %v", err)
		}
	}(f)

	// read file header to get file content_type
	part := make([]byte, 512)
	n, err := io.ReadFull(f, part)
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		err = nil
	}
	if err != nil {
		return "", errors2.Wrapf(err, "read reader to upload '%s'", key)
	}
	fileContentType := http.DetectContentType(part[:n])
	if _, err = f.Seek(0, io.SeekStart); err != nil { // seek to origin
		logs.CtxError(ctx, "upload file seek fail, err: %v", err)
	}

	logs.CtxDebug(ctx, "start upload, fileName: %s", fileName)
	if err = fs.client.Upload(ctx, fileName, f, fileserver.UploadWithContentType(fileContentType)); err != nil {
		logs.CtxError(ctx, "upload file failed, err: %v", err)
		return "", err
	}

	return fileName, nil
}

func (fs fileService) SignUploadFile(ctx context.Context, req *file.SignUploadFileRequest) (uris []string, heads []*file.SignHead, err error) {
	signOpt := make([]fileserver.SignOpt, 0)
	if req.Option != nil {
		if req.Option.TTL != nil && *req.Option.TTL > 0 {
			signOpt = append(signOpt, fileserver.SignWithTTL(time.Duration(*req.Option.TTL)*time.Second))
		}
	}

	signUrls, headers, err := fs.client.BatchSignUploadReq(ctx, req.Keys, signOpt...)
	if err != nil {
		return nil, nil, err
	}
	for _, signUrl := range signUrls {
		parsedURL, err := url.Parse(signUrl)
		if err != nil {
			return nil, nil, err
		}
		if parsedURL.Host == localos.GetLocalOSHost() {
			signUrl = fmt.Sprintf("%s?%s", parsedURL.Path, parsedURL.RawQuery)
		}
		uris = append(uris, signUrl)
	}

	heads = make([]*file.SignHead, 0)
	for _, header := range headers {
		h := &file.SignHead{}
		for key, value := range header {
			if len(value) == 0 {
				continue
			}
			if key == HeaderAccessKeyId {
				h.AccessKeyID = &value[0]
			}
			if key == HeaderSecretAccessKey {
				h.SecretAccessKey = &value[0]
			}
			if key == HeaderSessionToken {
				h.SessionToken = &value[0]
			}
			if key == HeaderExpiredTime {
				h.ExpiredTime = &value[0]
			}
			if key == HeaderCurrentTime {
				h.CurrentTime = &value[0]
			}
		}
		heads = append(heads, h)
	}

	return uris, heads, nil
}

func (fs fileService) SignDownLoadFile(ctx context.Context, req *file.SignDownloadFileRequest) (uris []string, err error) {
	signOpt := make([]fileserver.SignOpt, 0)
	if req.Option != nil {
		if req.Option.TTL != nil && *req.Option.TTL > 0 {
			signOpt = append(signOpt, fileserver.SignWithTTL(time.Duration(*req.Option.TTL)*time.Second))
		}
	}

	signUrls, _, err := fs.client.BatchSignDownloadReq(ctx, req.Keys, signOpt...)
	if err != nil {
		return nil, err
	}
	for _, signUrl := range signUrls {
		parsedURL, err := url.Parse(signUrl)
		if err != nil {
			return nil, err
		}
		if parsedURL.Host == localos.GetLocalOSHost() {
			signUrl = fmt.Sprintf("%s?%s", parsedURL.Path, parsedURL.RawQuery)
		}
		uris = append(uris, signUrl)
	}

	return uris, nil
}
