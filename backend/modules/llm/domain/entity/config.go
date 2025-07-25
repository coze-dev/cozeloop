// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

type RuntimeConfig struct {
	NeedCvtURLToBase64 bool   `json:"need_cvt_url_to_base_64" yaml:"need_cvt_url_to_base_64" mapstructure:"need_cvt_url_to_base_64"`
	QianfanAk          string `json:"qianfan_ak" yaml:"qianfan_ak" mapstructure:"qianfan_ak"`
	QianfanSk          string `json:"qianfan_sk" yaml:"qianfan_sk" mapstructure:"qianfan_sk"`
}
