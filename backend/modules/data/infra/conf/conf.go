// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"context"

	"github.com/samber/lo"

	dataset_conf "github.com/coze-dev/coze-loop/backend/modules/data/domain/component/conf"
	"github.com/coze-dev/coze-loop/backend/modules/data/pkg/consts"
	"github.com/coze-dev/coze-loop/backend/pkg/conf"
)

func NewConfiger(configFactory conf.IConfigLoaderFactory) (dataset_conf.IConfig, error) {
	loader, err := configFactory.NewConfigLoader(consts.DataConfigFileName)
	if err != nil {
		return nil, err
	}
	return &configer{
		loader: loader,
	}, nil
}

type configer struct {
	loader conf.IConfigLoader
}

func (c *configer) GetConsumerConfigs() *dataset_conf.ConsumerConfig {
	const key = "consumer_configs"
	var conf *dataset_conf.ConsumerConfig
	return lo.Ternary(c.loader.UnmarshalKey(context.Background(), key, &conf) == nil, conf, &dataset_conf.ConsumerConfig{})
}

func (c *configer) GetSnapshotRetry() *dataset_conf.SnapshotRetry {
	const key = "snapshot_retry"
	var conf *dataset_conf.SnapshotRetry
	return lo.Ternary(c.loader.UnmarshalKey(context.Background(), key, &conf) == nil, conf, &dataset_conf.SnapshotRetry{})
}

func (c *configer) GetProducerConfig() *dataset_conf.ProducerConfig {
	const key = "job_mq_producer"
	var conf *dataset_conf.ProducerConfig
	return lo.Ternary(c.loader.UnmarshalKey(context.Background(), key, &conf) == nil, conf, &dataset_conf.ProducerConfig{})
}

func (c *configer) GetDatasetFeature() *dataset_conf.DatasetFeature {
	const key = "default_dataset_feature"
	var conf *dataset_conf.DatasetFeature
	return lo.Ternary(c.loader.UnmarshalKey(context.Background(), key, &conf) == nil, conf, &dataset_conf.DatasetFeature{})
}

func (c *configer) GetDatasetItemStorage() *dataset_conf.DatasetItemStorage {
	const key = "dataset_item_storage"
	var conf *dataset_conf.DatasetItemStorage
	return lo.Ternary(c.loader.UnmarshalKey(context.Background(), key, &conf) == nil, conf, &dataset_conf.DatasetItemStorage{})
}

func (c *configer) GetDatasetSpec() *dataset_conf.DatasetSpec {
	const key = "default_dataset_spec"
	var conf *dataset_conf.DatasetSpec
	return lo.Ternary(c.loader.UnmarshalKey(context.Background(), key, &conf) != nil, conf, &dataset_conf.DatasetSpec{})
}
