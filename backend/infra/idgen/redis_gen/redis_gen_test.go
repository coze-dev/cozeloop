// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package redis_gen

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func newTestRedis(t *testing.T) (*redis.Client, error) {
	m := miniredis.NewMiniRedis()
	if err := m.Start(); err != nil {
		return nil, err
	}

	opts := &redis.Options{Addr: m.Addr()}
	p := redis.NewClient(opts)

	t.Cleanup(m.Close)
	return p, nil
}

func Test_generator_GenMultiIDs(t *testing.T) {
	ctx := context.Background()

	rcli, err := newTestRedis(t)
	assert.Nil(t, err)

	idgen, err := NewIDGenerator(rcli, []int64{0, 1, 2})
	assert.Nil(t, err)

	ids, err := idgen.GenMultiIDs(ctx, 10)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(ids))

	id, err := idgen.GenID(ctx)
	assert.Nil(t, err)
	assert.True(t, id >= time.Now().UnixNano()-(id>>32))
}
