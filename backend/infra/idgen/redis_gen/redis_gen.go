// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package redis_gen

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"

	"github.com/coze-dev/coze-loop/backend/infra/idgen"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
)

const (
	maxCounter = (1 << 8) - 1

	counterKeyExpiration = 10 * time.Minute
)

// NewIDGenerator 32b timestamp + 10b timestamp+ 8b counter + 14b serverid
func NewIDGenerator(client *redis.Client, srvIDs []int64) (idgen.IIDGenerator, error) {
	if len(srvIDs) == 0 {
		return nil, errors.New("idgen must init with valid server ids")
	}
	return &generator{
		cli:    client,
		srvIDs: srvIDs,
	}, nil
}

type generator struct {
	cli       *redis.Client
	srvIDs    []int64
	namespace string
}

func (i *generator) GenID(ctx context.Context) (int64, error) {
	ids, err := i.GenMultiIDs(ctx, 1)
	if err != nil {
		return 0, errorx.Wrapf(err, "failed to generate id") // todo 再细化一下
	}
	return ids[0], nil
}

func (i *generator) GenMultiIDs(ctx context.Context, counts int) ([]int64, error) {
	const maxTimeAddrTimes = 8

	leftNum := int64(counts)
	lastMs := int64(0)
	ids := make([]int64, 0, counts)
	svrID := i.pickSvrID()

	for idx := int64(0); leftNum > 0 && idx < maxTimeAddrTimes; idx++ {
		ms := lo.Ternary(i.timeMS() > lastMs, i.timeMS(), lastMs)
		if ms <= lastMs {
			ms++
		}

		lastMs = ms
		redisKey := i.counterKey(i.namespace, svrID, ms)

		counter, err := i.incrBy(ctx, redisKey, leftNum)
		if err != nil {
			return nil, err
		}

		var start, end int64

		start = counter - leftNum
		if start == 0 {
			i.expire(ctx, redisKey)
		}

		if start > maxCounter {
			continue
		} else if counter < leftNum {
			return nil, fmt.Errorf("recycling of counting space occurs, ms=%v", ms)
		}

		if counter > maxCounter {
			end = maxCounter + 1
			leftNum = counter - maxCounter - 1
		} else {
			end = counter
			leftNum = 0
		}

		seconds := ms / 1000
		millis := ms % 1000

		if seconds&0xFFFFFFFF != seconds {
			return nil, fmt.Errorf("seconds more than 32 bits, seconds=%v", seconds)
		}

		if svrID&0x3FFF != svrID {
			return nil, fmt.Errorf("server id more than 14 bits, serverID=%v", svrID)
		}

		for i := start; i < end; i++ {
			id := (seconds)<<32 + (millis)<<22 + i<<14 + svrID
			ids = append(ids, id)
		}
	}

	if len(ids) < counts || leftNum != 0 {
		return nil, fmt.Errorf("IDs num not enough, ns=%v, expect=%v, gotten=%v, lastMs=%v", i.namespace, counts, len(ids), lastMs)
	}

	return ids, nil
}

func (i *generator) incrBy(ctx context.Context, key string, num int64) (cntPos int64, err error) {
	return i.cli.IncrBy(ctx, key, num).Result()
}

func (i *generator) expire(ctx context.Context, key string) {
	_, _ = i.cli.Expire(ctx, key, counterKeyExpiration).Result()
}

func (i *generator) timeMS() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (i *generator) counterKey(space string, svrID int64, ms int64) string {
	return fmt.Sprintf("id_generator:%v:%v:%v", space, svrID, ms)
}

func (i *generator) pickSvrID() int64 {
	return i.srvIDs[rand.Intn(len(i.srvIDs))]
}
