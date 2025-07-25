// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package fileserver

import (
	"context"
	"io"
	"net/http"
	"time"
)

// ObjectStorage defines the interface for object storage backend(e.g. S3, OSS, TOS etc.)
//
//go:generate mockgen -destination=mocks/object_storage.go -package=mocks . ObjectStorage
type ObjectStorage interface {
	Stat(ctx context.Context, key string, opts ...StatOpt) (*ObjectInfo, error)
	Upload(ctx context.Context, key string, r io.Reader, opts ...UploadOpt) error
	Download(ctx context.Context, key string, w io.WriterAt, opts ...DownloadOpt) error
	Read(ctx context.Context, key string, opts ...DownloadOpt) (Reader, error)
	Remove(ctx context.Context, key string, opts ...RemoveOpt) error

	SignUploadReq(ctx context.Context, key string, opts ...SignOpt) (url string, header http.Header, err error)
	SignDownloadReq(ctx context.Context, key string, opts ...SignOpt) (url string, header http.Header, err error)
}

//go:generate mockgen -destination=mocks/batch_object_storage.go -package=mocks . BatchObjectStorage
type BatchObjectStorage interface {
	ObjectStorage
	BatchUpload(ctx context.Context, keys []string, readers []io.Reader, opts ...UploadOpt) error
	BatchDownload(ctx context.Context, keys []string, writers []io.WriterAt, opts ...DownloadOpt) error
	BatchRead(ctx context.Context, keys []string, opts ...DownloadOpt) ([]Reader, error)

	BatchSignDownloadReq(ctx context.Context, keys []string, opts ...SignOpt) (urls []string, headers []http.Header, err error)
	BatchSignUploadReq(ctx context.Context, keys []string, opts ...SignOpt) (urls []string, headers []http.Header, err error)
}

type Reader interface {
	io.ReadCloser
	io.ReaderAt
}

type Opt func(*Option)

type Option struct {
	Bucket      string
	Concurrency int
	Timeout     time.Duration
}

func (o *Option) ContextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if t := o.Timeout; t > 0 {
		return context.WithTimeout(ctx, t)
	}
	return ctx, func() {}
}

type StatOpt = Opt

type StatOption = Option

func NewStatOption(opts ...StatOpt) *StatOption {
	o := &StatOption{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func StatWithBucket(bucket string) StatOpt {
	return func(o *StatOption) { o.Bucket = bucket }
}

type DownloadOpt = Opt

type DownloadOption = Option

func NewDownloadOption(opts ...DownloadOpt) *DownloadOption {
	o := &DownloadOption{Concurrency: DefaultConcurrency}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func DownloadWithBucket(bucket string) DownloadOpt {
	return func(o *DownloadOption) { o.Bucket = bucket }
}

func DownloadWithConcurrency(concurrency int) DownloadOpt {
	return func(o *DownloadOption) { o.Concurrency = concurrency }
}

func DownloadWithTimeout(timeout time.Duration) DownloadOpt {
	return func(o *DownloadOption) { o.Timeout = timeout }
}

type RemoveOpt = Opt

type RemoveOption = Option

func NewRemoveOption(opts ...RemoveOpt) *RemoveOption {
	o := &RemoveOption{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func RemoveWithBucket(bucket string) RemoveOpt {
	return func(o *RemoveOption) { o.Bucket = bucket }
}

type SignOpt func(*SignOption)

type SignOption struct {
	Bucket        string
	ImageTemplate string
	Format        string
	TTL           time.Duration
}

func NewSignOption(opts ...SignOpt) *SignOption {
	o := &SignOption{TTL: DefaultSignTTL}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// SignWithImageTemplate sets the image template for the signed URL.
// Notice: not all backends support this option.
func SignWithImageTemplate(template string) SignOpt {
	return func(o *SignOption) { o.ImageTemplate = template }
}

// SignWithFormat sets the format for the signed URL.
// Notice: not all backends support this option.
func SignWithFormat(format string) SignOpt {
	return func(o *SignOption) { o.Format = format }
}

// SignWithTTL specifies the TTL for the signed URL.
// Default TTL is 24 hours.
func SignWithTTL(ttl time.Duration) SignOpt {
	return func(o *SignOption) { o.TTL = ttl }
}

// SignWithBucket specifies the bucket for the signed URL.
// Default bucket is the bucket of the object storage.
func SignWithBucket(bucket string) SignOpt {
	return func(o *SignOption) { o.Bucket = bucket }
}

type UploadOpt func(*UploadOption)

type UploadOption struct {
	Bucket       string
	Metadata     map[string]string
	Concurrency  int
	ContentTypes []string
}

func NewUploadOption(opts ...UploadOpt) *UploadOption {
	o := &UploadOption{Concurrency: DefaultConcurrency}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// UploadWithMetadata sets the metadata for the uploaded object.
func UploadWithMetadata(metadata map[string]string) UploadOpt {
	return func(o *UploadOption) {
		// make a copy
		m := make(map[string]string)
		for k, v := range metadata {
			m[k] = v
		}
		o.Metadata = m
	}
}

// UploadWithBucket specifies the bucket for the upload.
// Default bucket is the bucket of the object storage.
func UploadWithBucket(bucket string) UploadOpt {
	return func(o *UploadOption) { o.Bucket = bucket }
}

// UploadWithConcurrency sets the concurrency for the upload.
// Default concurrency is 3, used in batch upload.
func UploadWithConcurrency(concurrency int) UploadOpt {
	return func(o *UploadOption) { o.Concurrency = concurrency }
}

// UploadWithContentType sets the content_type for the upload.
// Default content_type is None.
func UploadWithContentType(contentTypes ...string) UploadOpt {
	return func(o *UploadOption) {
		o.ContentTypes = make([]string, len(contentTypes))
		copy(o.ContentTypes, contentTypes)
	}
}
