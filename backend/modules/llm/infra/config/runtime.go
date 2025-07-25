// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/qianfan"

	llm_conf "github.com/coze-dev/cozeloop/backend/modules/llm/domain/component/conf"
	"github.com/coze-dev/cozeloop/backend/modules/llm/domain/entity"
	"github.com/coze-dev/cozeloop/backend/pkg/conf"
)

type RuntimeImpl struct {
	cfg *entity.RuntimeConfig
}

func NewRuntime(ctx context.Context, factory conf.IConfigLoaderFactory) (llm_conf.IConfigRuntime, error) {
	loader, err := factory.NewConfigLoader("model_runtime_config.yaml")
	if err != nil {
		return nil, err
	}
	var cfg entity.RuntimeConfig
	if err = loader.Unmarshal(ctx, &cfg); err != nil {
		return nil, err
	}
	qianfanCfg := qianfan.GetQianfanSingletonConfig()
	qianfanCfg.AccessKey = cfg.QianfanAk
	qianfanCfg.SecretKey = cfg.QianfanSk
	return &RuntimeImpl{
		cfg: &cfg,
	}, nil
}

func (r *RuntimeImpl) NeedCvtURLToBase64() bool {
	if r == nil || r.cfg == nil {
		return false
	}
	return r.cfg.NeedCvtURLToBase64
}
