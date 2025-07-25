// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	"github.com/coze-dev/cozeloop/backend/modules/llm/domain/component/conf"
	"github.com/coze-dev/cozeloop/backend/modules/llm/domain/repo"
	"github.com/coze-dev/cozeloop/backend/modules/llm/domain/service/llmfactory"
)

func NewRuntime(
	llmFact llmfactory.IFactory,
	idGen idgen.IIDGenerator,
	runtimeRepo repo.IRuntimeRepo,
	cfg conf.IConfigRuntime,
) IRuntime {
	return &RuntimeImpl{
		llmFact:     llmFact,
		idGen:       idGen,
		runtimeRepo: runtimeRepo,
		runtimeCfg:  cfg,
	}
}

func NewManage(cfg conf.IConfigManage) IManage {
	return &ManageImpl{
		conf: cfg,
	}
}
