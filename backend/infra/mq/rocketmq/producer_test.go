// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package rocketmq

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToDelayLevel(t *testing.T) {
	testCases := []struct {
		name           string
		inputDuration  time.Duration
		expectedResult int
	}{
		{
			name:           "小于最小延迟级别",
			inputDuration:  time.Second,
			expectedResult: 1,
		},
		{
			name:           "在延迟级别范围内",
			inputDuration:  35 * time.Second,
			expectedResult: 5,
		},
		{
			name:           "大于最大延迟级别",
			inputDuration:  8000 * time.Second,
			expectedResult: len(delayLevels),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &Producer{}
			result := p.toDelayLevel(tc.inputDuration)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestCustomNsResolver(t *testing.T) {
	testCases := []struct {
		InputAddrs  []string
		OutputAddrs []string
	}{
		{
			InputAddrs:  []string{"1.2.3.4:1000"},
			OutputAddrs: []string{"1.2.3.4:1000"},
		},
		{
			InputAddrs:  []string{"1.2.3.4"},
			OutputAddrs: []string{"1.2.3.4"},
		},
		{
			InputAddrs:  []string{"a.b.c"},
			OutputAddrs: []string{"a.b.c"},
		},
		{
			InputAddrs:  []string{"a.b.c:8000"},
			OutputAddrs: []string{"a.b.c:8000"},
		},
	}
	for _, tc := range testCases {
		ns := NewCustomResolver(tc.InputAddrs)
		out := ns.Resolve()
		assert.Equal(t, tc.OutputAddrs, out)
	}
}
