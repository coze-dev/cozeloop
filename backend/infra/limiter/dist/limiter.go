// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dist

import (
	"context"
	"fmt"
	"strconv"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/go-redis/redis_rate/v10"
	"github.com/samber/lo"

	"github.com/coze-dev/coze-loop/backend/infra/limiter"
	"github.com/coze-dev/coze-loop/backend/infra/redis"
	"github.com/coze-dev/coze-loop/backend/pkg/json"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/conv"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/slices"
	"github.com/coze-dev/coze-loop/backend/pkg/logs"
	"github.com/coze-dev/coze-loop/backend/pkg/mcache"
	"github.com/coze-dev/coze-loop/backend/pkg/mcache/byted"
)

type factory struct {
	cmdable   redis.Cmdable
	cacheSize int // cache size in byte
}

func (f *factory) NewRateLimiter(opts ...limiter.FactoryOptionFn) limiter.IRateLimiter {
	opt := &limiter.FactoryOption{}
	for _, fn := range opts {
		fn(opt)
	}

	rawRedis, ok := redis.Unwrap(f.cmdable)
	if !ok {
		panic(fmt.Errorf("redis cmdable must be unwrappable"))
	}

	rl := &rateLimiter{
		rules:   make([]*rule, 0, len(opt.Rules)),
		limiter: redis_rate.NewLimiter(rawRedis),
		vmCache: byted.NewLRUCache(lo.Ternary(f.cacheSize > 0, f.cacheSize, 5*1024*1024)),
	}

	for _, r := range opt.Rules {
		if rr, err := rl.newRule(r); err != nil {
			logs.Error("rate limiter set rule failed, rule: %v, err: %v", r, err)
		} else {
			rl.addRule(rr)
		}
	}

	return rl
}

const (
	originKeyPlaceholderKey = "origin_key"
)

type rateLimiter struct {
	limiter *redis_rate.Limiter

	rules   []*rule
	vmCache mcache.IByteCache
}

type rule struct {
	limiter.Rule

	matchVM *vm.Program
	keyVM   *vm.Program
}

func (r rule) match(env map[string]any) bool {
	if r.matchVM != nil {
		res, err := expr.Run(r.matchVM, env)
		if err != nil {
			return false
		}
		return toBool(res)
	}
	return true
}

func (r rule) getKey(env map[string]any) string {
	if r.keyVM != nil {
		res, err := expr.Run(r.keyVM, env)
		if err == nil {
			return conv.ToString(res)
		}
	}
	okey, _ := env[originKeyPlaceholderKey].(string)
	return okey
}

func (rl *rateLimiter) AllowN(ctx context.Context, key string, n int, opts ...limiter.LimitOptionFn) (*limiter.Result, error) {
	opt := &limiter.LimitOption{}
	for _, fn := range opts {
		fn(opt)
	}

	limitKey, limit := rl.getLimitKey(key, opt.Tags)
	limit = lo.Ternary(opt.Limit != nil, opt.Limit, limit)

	if len(limitKey) == 0 || limit == nil {
		logs.Warn("AllowN with invalid limit key: %v, tags: %v", limitKey, opt.Tags)
		return &limiter.Result{
			Allowed:   true,
			OriginKey: key,
			LimitKey:  limitKey,
		}, nil
	}

	result, err := rl.limiter.AllowN(ctx, limitKey, redis_rate.Limit{
		Rate:   limit.Rate,
		Burst:  limit.Burst,
		Period: limit.Period,
	}, n)
	if err != nil {
		return nil, fmt.Errorf("redis limiter allown fail, origin_key=%s, limit_key=%s, err=%w", key, limitKey, err)
	}

	return &limiter.Result{
		Allowed:   result.Allowed > 0,
		OriginKey: key,
		LimitKey:  limitKey,
	}, nil
}

func (rl *rateLimiter) getLimitKey(key string, tags []limiter.Tag) (LimitKey string, limit *limiter.Limit) {
	env := slices.ToMap(tags, func(e limiter.Tag) (string, any) {
		return e.K, e.V
	})
	env[originKeyPlaceholderKey] = key

	for _, r := range rl.rules {
		if !r.match(env) {
			continue
		}
		if k := r.getKey(env); len(k) > 0 {
			return k, ptr.Of(r.Limit)
		}
	}

	return key, nil
}

func (rl *rateLimiter) getExprVM(prefix, exprStr string) (*vm.Program, error) {
	k := prefix + exprStr
	kb := conv.UnsafeStringToBytes(k)

	got, err := rl.vmCache.Get(kb)
	if err == nil {
		p := new(vm.Program)
		if err := json.Unmarshal(got, p); err == nil {
			return p, nil
		}
	}

	p, err := expr.Compile(exprStr)
	if err != nil {
		return nil, fmt.Errorf("expr compile with invalid str, str=%s, err=%w", exprStr, err)
	}

	_ = rl.vmCache.Set(kb, conv.UnsafeStringToBytes(json.Jsonify(p)), 0)

	return p, nil
}

func toBool(v any) bool {
	if v == nil {
		return false
	}
	switch t := v.(type) {
	case bool:
		return t
	case string:
		if vv, err := strconv.ParseBool(t); err == nil {
			return vv
		}
		return false
	default:
		return false
	}
}

func (rl *rateLimiter) newRule(r limiter.Rule) (*rule, error) {
	rr := &rule{Rule: r}

	if len(r.Match) > 0 {
		match, err := rl.getExprVM("match", r.Match)
		if err != nil {
			return nil, err
		}
		rr.matchVM = match
	}

	if len(r.KeyExpr) > 0 {
		key, err := rl.getExprVM("key", r.KeyExpr)
		if err != nil {
			return nil, err
		}
		rr.keyVM = key
	}

	return rr, nil
}

func (rl *rateLimiter) addRule(r *rule) {
	rl.rules = append(rl.rules, r)
}
