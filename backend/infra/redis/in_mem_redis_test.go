// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTestRedis(t *testing.T) {
	cli := NewTestRedis(t)
	ctx := context.TODO()

	key := `resource:1:name`
	require.NoError(t, cli.Set(ctx, key, `alice`, time.Hour).Err())

	got, err := cli.Get(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, `alice`, got)
}
