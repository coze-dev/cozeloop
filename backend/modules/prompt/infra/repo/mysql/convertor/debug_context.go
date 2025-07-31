// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"github.com/coze-dev/coze-loop/backend/modules/prompt/domain/entity"
	"github.com/coze-dev/coze-loop/backend/modules/prompt/infra/repo/mysql/gorm_gen/model"
	"github.com/coze-dev/coze-loop/backend/pkg/json"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
)

func DebugContextDO2PO(do *entity.DebugContext) (*model.PromptDebugContext, error) {
	if do == nil {
		return nil, nil
	}
	var mockContexts, mockVariables, mockTools, debugConfig, compareConfig *string
	if do.DebugCore != nil {
		mockContextsStr, err := json.MarshalString(do.DebugCore.MockContexts)
		if err != nil {
			return nil, err
		}
		mockContexts = ptr.Of(mockContextsStr)
		mockVariablesStr, err := json.MarshalString(do.DebugCore.MockVariables)
		if err != nil {
			return nil, err
		}
		mockVariables = ptr.Of(mockVariablesStr)
		mockToolsStr, err := json.MarshalString(do.DebugCore.MockTools)
		if err != nil {
			return nil, err
		}
		mockTools = ptr.Of(mockToolsStr)
	}
	if do.DebugConfig != nil {
		debugConfigStr, err := json.MarshalString(do.DebugConfig)
		if err != nil {
			return nil, err
		}
		debugConfig = ptr.Of(debugConfigStr)
	}
	if do.CompareConfig != nil {
		compareConfigStr, err := json.MarshalString(do.CompareConfig)
		if err != nil {
			return nil, err
		}
		compareConfig = ptr.Of(compareConfigStr)
	}
	return &model.PromptDebugContext{
		PromptID:      do.PromptID,
		UserID:        do.UserID,
		MockContexts:  mockContexts,
		MockVariables: mockVariables,
		MockTools:     mockTools,
		DebugConfig:   debugConfig,
		CompareConfig: compareConfig,
	}, nil
}

//===============================================================

func DebugContextPO2DO(po *model.PromptDebugContext) (*entity.DebugContext, error) {
	if po == nil {
		return nil, nil
	}
	var mockContexts []*entity.DebugMessage
	if po.MockContexts != nil {
		err := json.Unmarshal([]byte(ptr.From(po.MockContexts)), &mockContexts)
		if err != nil {
			return nil, err
		}
	}
	var mockVariables []*entity.VariableVal
	if po.MockVariables != nil {
		err := json.Unmarshal([]byte(ptr.From(po.MockVariables)), &mockVariables)
		if err != nil {
			return nil, err
		}
	}
	var mockTools []*entity.MockTool
	if po.MockTools != nil {
		err := json.Unmarshal([]byte(ptr.From(po.MockTools)), &mockTools)
		if err != nil {
			return nil, err
		}
	}
	var debugConfig *entity.DebugConfig
	if po.DebugConfig != nil {
		err := json.Unmarshal([]byte(ptr.From(po.DebugConfig)), &debugConfig)
		if err != nil {
			return nil, err
		}
	}
	var compareConfig *entity.CompareConfig
	if po.CompareConfig != nil {
		err := json.Unmarshal([]byte(ptr.From(po.CompareConfig)), &compareConfig)
		if err != nil {
			return nil, err
		}
	}
	return &entity.DebugContext{
		PromptID: po.PromptID,
		UserID:   po.UserID,
		DebugCore: &entity.DebugCore{
			MockContexts:  mockContexts,
			MockVariables: mockVariables,
			MockTools:     mockTools,
		},
		DebugConfig:   debugConfig,
		CompareConfig: compareConfig,
	}, nil
}
