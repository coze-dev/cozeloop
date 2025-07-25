// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
)

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name       string
		preVersion string
		newVersion string
		wantErr    bool
		errType    error
	}{
		{
			name:       "有效的新版本号",
			preVersion: "",
			newVersion: "1.0.0",
			wantErr:    false,
		},
		{
			name:       "有效的版本号更新",
			preVersion: "1.0.0",
			newVersion: "1.0.1",
			wantErr:    false,
		},
		{
			name:       "新版本号无效",
			preVersion: "",
			newVersion: "invalid",
			wantErr:    true,
			errType:    errno.InvalidParamErr(nil, "version 'invalid' not a valid semantic version", "invalid"),
		},
		{
			name:       "历史版本号无效",
			preVersion: "invalid",
			newVersion: "1.0.0",
			wantErr:    true,
			errType:    errno.InternalErr(nil, "previous version 'invalid' not a valid semantic version", "invalid"),
		},
		{
			name:       "新版本号小于历史版本号",
			preVersion: "2.0.0",
			newVersion: "1.0.0",
			wantErr:    true,
			errType:    errno.InvalidParamErrorf("new version '1.0.0' should be greater than '2.0.0'"),
		},
		{
			name:       "新版本号等于历史版本号",
			preVersion: "1.0.0",
			newVersion: "1.0.0",
			wantErr:    true,
			errType:    errno.InvalidParamErrorf("new version '1.0.0' should be greater than '1.0.0'"),
		},
		{
			name:       "复杂版本号比较",
			preVersion: "1.2.3",
			newVersion: "1.2.4-alpha.1",
			wantErr:    false,
		},
		{
			name:       "主版本号更新",
			preVersion: "1.9.9",
			newVersion: "2.0.0",
			wantErr:    false,
		},
		{
			name:       "次版本号更新",
			preVersion: "1.1.9",
			newVersion: "1.2.0",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVersion(tt.preVersion, tt.newVersion)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
