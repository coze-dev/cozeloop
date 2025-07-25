// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package unittest

import (
	"strings"
	"testing"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
)

// AssertErrorEqual is a helper function that compares errors properly based on their type
func AssertErrorEqual(t *testing.T, expected, actual error) {
	if expected == nil && actual == nil {
		return
	}

	if expected == nil || actual == nil {
		assert.Equal(t, expected, actual)
		return
	}

	// 处理errorx错误
	if expectedStatusErr, ok := kerrors.FromBizStatusError(expected); ok {
		actualStatusErr, ok := kerrors.FromBizStatusError(actual)
		assert.True(t, ok)
		// 比较错误码
		assert.Equal(t, expectedStatusErr.BizStatusCode(), actualStatusErr.BizStatusCode())
		return
	}

	// 处理标准错误，只比较错误消息，去掉堆栈
	assert.Equal(t, strings.TrimSpace(errorx.ErrorWithoutStack(expected)), strings.TrimSpace(errorx.ErrorWithoutStack(actual)))
}
