// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package queueprocessor

import "fmt"

type Config struct {
	PoolName        string `mapstructure:"pool_name"`
	MaxPoolSize     int32  `mapstructure:"max_pool_size"`
	QueueSize       int64  `mapstructure:"queue_size"`
	MaxBatchSize    int    `mapstructure:"max_batch_size"`
	TickIntervalsMs int64  `mapstructure:"tick_intervals_ms"`
	ShardCount      int    `mapstructure:"shard_count"`
}

func (cfg *Config) Validate() error {
	if cfg.QueueSize <= 0 {
		return fmt.Errorf("empty queue size")
	}
	if cfg.MaxPoolSize <= 0 {
		return fmt.Errorf("queue processor empty pool size")
	}
	if cfg.PoolName == "" {
		return fmt.Errorf("queue processor empty pool name")
	}
	if cfg.MaxBatchSize <= 0 {
		return fmt.Errorf("queue processor empty max batch size")
	}
	if cfg.TickIntervalsMs <= 0 {
		return fmt.Errorf("queue processor empty tick intervals ms")
	}
	if cfg.ShardCount <= 0 {
		return fmt.Errorf("queue processor empty shard count")
	}
	return nil
}
