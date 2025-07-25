// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/stretchr/testify/assert"
)

func TestKiteXBizStatusError(t *testing.T) {
	t.Run("statusError", func(t *testing.T) {
		oriErr := &statusError{
			statusCode: 1,
			message:    "err msg",
			ext: Extension{
				IsAffectStability: true,
			},
		}
		var bizErr KiteXBizStatusError
		ok := errors.As(oriErr, &bizErr)

		assert.True(t, ok)
		assert.Equal(t, int32(1), bizErr.BizStatusCode())
		assert.Equal(t, "err msg", bizErr.BizMessage())
		assert.Equal(t, "1", bizErr.BizExtra()[BizExtraKeyAffectStability])
	})

	t.Run("NewByCode", func(t *testing.T) {
		err := NewByCode(1)
		var bizErr KiteXBizStatusError
		ok := errors.As(err, &bizErr)

		assert.True(t, ok)
		assert.Equal(t, int32(1), bizErr.BizStatusCode())
		assert.Equal(t, DefaultErrorMsg, bizErr.BizMessage())
		assert.Equal(t, "1", bizErr.BizExtra()[BizExtraKeyAffectStability])
	})

	t.Run("WrapByCode", func(t *testing.T) {
		oriErr := errors.New("original error")
		err := WrapByCode(oriErr, 1)
		var bizErr KiteXBizStatusError
		ok := errors.As(err, &bizErr)

		assert.True(t, ok)
		assert.Equal(t, int32(1), bizErr.BizStatusCode())
		assert.Equal(t, DefaultErrorMsg, bizErr.BizMessage())
		assert.Equal(t, "1", bizErr.BizExtra()[BizExtraKeyAffectStability])
	})

	t.Run("wrapf", func(t *testing.T) {
		oriErr := &statusError{
			statusCode: 1,
			message:    "err msg",
			ext: Extension{
				IsAffectStability: true,
			},
		}
		err := wrapf(oriErr, "wrap err")
		var bizErr KiteXBizStatusError
		ok := errors.As(err, &bizErr)

		assert.True(t, ok)
		assert.Equal(t, int32(1), bizErr.BizStatusCode())
		assert.Equal(t, "err msg", bizErr.BizMessage())
		assert.Equal(t, "1", bizErr.BizExtra()[BizExtraKeyAffectStability])
	})

	t.Run("IsAffectStability = false", func(t *testing.T) {
		err := &statusError{
			statusCode: 1,
			message:    "msg",
			ext: Extension{
				IsAffectStability: false,
			},
		}
		var bizErr KiteXBizStatusError
		ok := errors.As(err, &bizErr)

		assert.True(t, ok)
		assert.Equal(t, int32(1), bizErr.BizStatusCode())
		assert.Equal(t, "msg", bizErr.BizMessage())
		assert.Equal(t, "0", bizErr.BizExtra()[BizExtraKeyAffectStability])
	})
}

func TestKitexBizErrWrapper(t *testing.T) {
	err := NewByCode(1)
	bizErr, ok := kerrors.FromBizStatusError(err)
	assert.True(t, ok)
	errStr := bizErr.Error()
	assert.False(t, strings.Contains(errStr, "kitex_test.go"))
	verbErrStr := fmt.Sprintf("%+v", bizErr)
	assert.True(t, strings.Contains(verbErrStr, "kitex_test.go"))
	fmtErrStr := fmt.Sprintf("%v", bizErr)
	assert.False(t, strings.Contains(fmtErrStr, "kitex_test.go"))

	compared := NewByCode(1)
	assert.True(t, errors.Is(bizErr, compared))

	as, ok := bizErr.BizExtra()[BizExtraKeyAffectStability]
	assert.True(t, ok)
	assert.Equal(t, as, "1")

	se, ok := FromStatusError(err)
	assert.True(t, ok)
	assert.Equal(t, se.Code(), int32(1))

	mockUpstreamErr := &testUpstreamBizErr{bizErr}
	se, ok = FromStatusError(mockUpstreamErr)
	assert.True(t, ok)
	assert.Equal(t, se.Code(), int32(1))
}

type testUpstreamBizErr struct {
	raw kerrors.BizStatusErrorIface
}

func (t *testUpstreamBizErr) BizStatusCode() int32 {
	return t.raw.BizStatusCode()
}

func (t *testUpstreamBizErr) BizMessage() string {
	return t.raw.BizMessage()
}

func (t *testUpstreamBizErr) BizExtra() map[string]string {
	return t.raw.BizExtra()
}

func (t *testUpstreamBizErr) Error() string {
	return t.raw.Error()
}
