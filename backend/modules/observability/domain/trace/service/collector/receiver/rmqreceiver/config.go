// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package rmqreceiver

import "fmt"

type Config struct {
	Addr          []string `mapstructure:"addr" json:"addr"`
	ConsumerGroup string   `mapstructure:"consumer_group" json:"consumer_group"`
	Topic         string   `mapstructure:"topic" json:"topic"`
	Timeout       int64    `mapstructure:"timeout" json:"timeout"`
}

func (cfg *Config) Validate() error {
	if len(cfg.Addr) == 0 {
		return fmt.Errorf("addr is empty")
	}
	if cfg.ConsumerGroup == "" {
		return fmt.Errorf("rmq receiver consumer_group is empty")
	}
	if cfg.Topic == "" {
		return fmt.Errorf("rmq receiver topic is empty")
	}
	return nil
}
