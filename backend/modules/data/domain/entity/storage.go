// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

type Provider string

const (
	ProviderUnknown Provider = ""
	ProviderTOS     Provider = "TOS"
	ProviderVETOS   Provider = "VETOS"
	ProviderHDFS    Provider = "HDFS"
	ProviderImageX  Provider = "ImageX"
	ProviderS3      Provider = "S3"
	ProviderLocalFS Provider = "LocalFS"
	ProviderAbase   Provider = "Abase"
	ProviderRDS     Provider = "RDS"
)

type UploadToken struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	ExpiredTime     string
	CurrentTime     string
}
