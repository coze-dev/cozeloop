// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dist

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/coze-loop/backend/infra/limiter"
	"github.com/coze-dev/coze-loop/backend/infra/redis"
	"github.com/coze-dev/coze-loop/backend/pkg/conf/viper"
)

func Test_rateLimiter_AllowN(t *testing.T) {
	ctx := context.Background()

	f := NewRateLimiterFactory(redis.NewTestRedis(t))

	l := f.NewRateLimiter(limiter.WithRules([]limiter.Rule{
		{
			Match:   "itag==1 && stag==\"a\"",
			KeyExpr: "origin_key+string(itag)+stag",
			Limit: limiter.Limit{
				Rate:   5,
				Burst:  20,
				Period: time.Second,
			},
		},
		{
			Match:   "itag==2 && stag==\"b\"",
			KeyExpr: "origin_key+string(itag)+stag",
			Limit: limiter.Limit{
				Rate:   10,
				Burst:  20,
				Period: time.Second,
			},
		},
		{
			Match:   "object.ITag==3 && object.STag==\"c\"",
			KeyExpr: "origin_key+string(object.ITag)+object.STag",
			Limit: limiter.Limit{
				Rate:   10,
				Burst:  20,
				Period: time.Second,
			},
		},
		{
			Match:   "object.itag==4 && object.stag==\"d\"",
			KeyExpr: "origin_key+string(object.itag)+object.stag",
			Limit: limiter.Limit{
				Rate:   10,
				Burst:  20,
				Period: time.Second,
			},
		},
		{
			Limit: limiter.Limit{
				Rate:   5,
				Burst:  5,
				Period: time.Second,
			},
		},
	}...))

	t.Run("allowed", func(t *testing.T) {
		res, err := l.AllowN(ctx, "key_to_limit", 1, limiter.WithTags([]limiter.Tag{
			{K: "itag", V: 1},
			{K: "stag", V: "a"},
		}...))
		assert.Nil(t, err)
		assert.True(t, res.Allowed)
		assert.Equal(t, "key_to_limit", res.OriginKey)
		assert.Equal(t, "key_to_limit1a", res.LimitKey)

		res, err = l.AllowN(ctx, "key_to_limit", 1, limiter.WithTags([]limiter.Tag{
			{K: "itag", V: 2},
			{K: "stag", V: "b"},
		}...))
		assert.Nil(t, err)
		assert.True(t, res.Allowed)
		assert.Equal(t, "key_to_limit", res.OriginKey)
		assert.Equal(t, "key_to_limit2b", res.LimitKey)
	})

	t.Run("not allowed", func(t *testing.T) {
		res, err := l.AllowN(ctx, "key_to_limit", 21, limiter.WithTags([]limiter.Tag{
			{K: "itag", V: 1},
			{K: "stag", V: "a"},
		}...))
		assert.Nil(t, err)
		assert.False(t, res.Allowed)
		assert.Equal(t, "key_to_limit", res.OriginKey)

		res, err = l.AllowN(ctx, "key_to_limit", 21, limiter.WithTags([]limiter.Tag{
			{K: "itag", V: 2},
			{K: "stag", V: "b"},
		}...))
		assert.Nil(t, err)
		assert.False(t, res.Allowed)
		assert.Equal(t, "key_to_limit", res.OriginKey)
	})

	t.Run("allowed with burst", func(t *testing.T) {
		time.Sleep(3 * time.Second)

		res, err := l.AllowN(ctx, "key_to_limit", 10, limiter.WithTags([]limiter.Tag{
			{K: "itag", V: 1},
			{K: "stag", V: "a"},
		}...))
		assert.Nil(t, err)
		assert.True(t, res.Allowed)
		assert.Equal(t, "key_to_limit", res.OriginKey)

		res, err = l.AllowN(ctx, "key_to_limit", 15, limiter.WithTags([]limiter.Tag{
			{K: "itag", V: 2},
			{K: "stag", V: "b"},
		}...))
		assert.Nil(t, err)
		assert.True(t, res.Allowed)
		assert.Equal(t, "key_to_limit", res.OriginKey)
	})

	t.Run("not allowed with batch", func(t *testing.T) {
		time.Sleep(3 * time.Second)

		res, err := l.AllowN(ctx, "key_to_limit", 10, limiter.WithTags([]limiter.Tag{
			{K: "itag", V: 1},
			{K: "stag", V: "a"},
		}...))
		assert.Nil(t, err)
		assert.True(t, res.Allowed)
		assert.Equal(t, "key_to_limit", res.OriginKey)

		res, err = l.AllowN(ctx, "key_to_limit", 11, limiter.WithTags([]limiter.Tag{
			{K: "itag", V: 1},
			{K: "stag", V: "a"},
		}...))
		assert.Nil(t, err)
		assert.False(t, res.Allowed)
		assert.Equal(t, "key_to_limit", res.OriginKey)
	})

	t.Run("match default", func(t *testing.T) {
		res, err := l.AllowN(ctx, "key_to_limit", 2, limiter.WithTags([]limiter.Tag{
			{K: "itag", V: 3},
			{K: "stag", V: "c"},
		}...))
		assert.Nil(t, err)
		assert.True(t, res.Allowed)
		assert.Equal(t, "key_to_limit", res.OriginKey)
		assert.Equal(t, "key_to_limit", res.LimitKey)

		res, err = l.AllowN(ctx, "key_to_limit", 2, limiter.WithTags([]limiter.Tag{
			{K: "itag", V: 3},
			{K: "stag", V: "c"},
		}...))
		assert.Nil(t, err)
		assert.True(t, res.Allowed)
		assert.Equal(t, "key_to_limit", res.OriginKey)
		assert.Equal(t, "key_to_limit", res.LimitKey)

		res, err = l.AllowN(ctx, "key_to_limit", 2, limiter.WithTags([]limiter.Tag{
			{K: "itag", V: 3},
			{K: "stag", V: "c"},
		}...))
		assert.Nil(t, err)
		assert.False(t, res.Allowed)
		assert.Equal(t, "key_to_limit", res.OriginKey)
		assert.Equal(t, "key_to_limit", res.LimitKey)

		time.Sleep(time.Second)
		res, err = l.AllowN(ctx, "key_to_limit", 6, limiter.WithTags([]limiter.Tag{
			{K: "itag", V: 3},
			{K: "stag", V: "c"},
		}...))
		assert.Nil(t, err)
		assert.False(t, res.Allowed)
		assert.Equal(t, "key_to_limit", res.OriginKey)
		assert.Equal(t, "key_to_limit", res.LimitKey)
	})

	t.Run("object tag", func(t *testing.T) {
		res, err := l.AllowN(ctx, "key_to_limit", 21, limiter.WithTags([]limiter.Tag{
			{K: "object", V: struct {
				ITag int64
				STag string
			}{
				ITag: 3,
				STag: "c",
			}},
		}...))
		assert.Nil(t, err)
		assert.False(t, res.Allowed)
		assert.Equal(t, "key_to_limit", res.OriginKey)
		assert.Equal(t, "key_to_limit3c", res.LimitKey)
	})

	t.Run("map tag", func(t *testing.T) {
		res, err := l.AllowN(ctx, "key_to_limit", 21, limiter.WithTags([]limiter.Tag{
			{K: "object", V: map[string]any{"itag": int64(4), "stag": "d"}},
		}...))
		assert.Nil(t, err)
		assert.False(t, res.Allowed)
		assert.Equal(t, "key_to_limit", res.OriginKey)
		assert.Equal(t, "key_to_limit4d", res.LimitKey)
	})
}

func TestLoadConf(t *testing.T) {
	ctx := context.Background()

	loader, err := viper.NewFileConfLoader("rule_test_conf.yaml")
	assert.Nil(t, err)

	var rules []limiter.Rule
	err = loader.UnmarshalKey(ctx, "rules", &rules)
	assert.Nil(t, err)

	f := NewRateLimiterFactory(redis.NewTestRedis(t))

	l := f.NewRateLimiter(limiter.WithRules(rules...))
	res, err := l.AllowN(ctx, "key_to_limit", 5, limiter.WithTags([]limiter.Tag{
		{K: "itag", V: 1},
		{K: "stag", V: "a"},
	}...))
	assert.Nil(t, err)
	assert.True(t, res.Allowed)
	assert.Equal(t, "key_to_limit", res.OriginKey)
	assert.Equal(t, "key_to_limit1a", res.LimitKey)

	res, err = l.AllowN(ctx, "key_to_limit", 5, limiter.WithTags([]limiter.Tag{
		{K: "itag", V: 2},
		{K: "stag", V: "b"},
	}...))
	assert.Nil(t, err)
	assert.False(t, res.Allowed)
	assert.Equal(t, "key_to_limit", res.OriginKey)
	assert.Equal(t, "key_to_limit", res.LimitKey)
}

func TestConcurrentAllowN(t *testing.T) {
	ctx := context.Background()

	f := NewRateLimiterFactory(redis.NewTestRedis(t))

	l := f.NewRateLimiter(limiter.WithRules([]limiter.Rule{
		{
			Match:   "itag==1 && stag==\"a\"",
			KeyExpr: "origin_key+string(itag)+stag",
			Limit: limiter.Limit{
				Rate:   150,
				Burst:  150,
				Period: time.Second,
			},
		},
	}...))

	t.Run("concurrency allowed", func(t *testing.T) {
		const num = 50
		wg := sync.WaitGroup{}
		wg.Add(num)
		for i := 0; i < num; i++ {
			go func() {
				defer wg.Done()
				res, err := l.AllowN(ctx, "key_to_limit", 2, limiter.WithTags([]limiter.Tag{
					{K: "itag", V: 1},
					{K: "stag", V: "a"},
				}...))
				assert.Nil(t, err)
				assert.True(t, res.Allowed)
			}()
		}
	})

	t.Run("concurrency", func(t *testing.T) {
		time.Sleep(time.Second)

		const num = 100
		wg := sync.WaitGroup{}
		wg.Add(num)
		for i := 0; i < num; i++ {
			notAllowed := i == num
			go func() {
				defer wg.Done()
				res, err := l.AllowN(ctx, "key_to_limit", 2, limiter.WithTags([]limiter.Tag{
					{K: "itag", V: 1},
					{K: "stag", V: "a"},
				}...))
				if notAllowed {
					assert.Nil(t, err)
					assert.False(t, res.Allowed)
				}
			}()
		}
	})
}
