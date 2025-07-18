// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import React from 'react';

import { useShallow } from 'zustand/react/shallow';
import { useRequest } from 'ahooks';
import { BasicModelConfigEditor } from '@cozeloop/prompt-components';
import { I18n } from '@cozeloop/i18n-adapter';
import { CollapseCard } from '@cozeloop/components';
import { useSpace } from '@cozeloop/biz-hooks-adapter';
import { Scenario } from '@cozeloop/api-schema/llm-manage';
import { LlmManageApi } from '@cozeloop/api-schema';
import { Typography } from '@coze-arch/coze-design';

import { usePromptStore } from '@/store/use-prompt-store';
import { useBasicStore } from '@/store/use-basic-store';
import { useCompare } from '@/hooks/use-compare';

export function ModelConfigCard() {
  const { spaceID } = useSpace();
  const { modelConfig, setModelConfig, setCurrentModel, promptInfo } =
    usePromptStore(
      useShallow(state => ({
        modelConfig: state.modelConfig,
        setModelConfig: state.setModelConfig,
        setCurrentModel: state.setCurrentModel,
        promptInfo: state.promptInfo,
        currentModel: state.currentModel,
      })),
    );
  const { readonly } = useBasicStore(
    useShallow(state => ({ readonly: state.readonly })),
  );
  const { streaming } = useCompare();

  const service = useRequest(async () =>
    LlmManageApi.ListModels({
      workspace_id: spaceID,
      page_size: 100,
      page_token: '0',
      scenario: Scenario.scenario_prompt_debug,
    }),
  );

  return (
    <CollapseCard
      title={<Typography.Text strong>{I18n.t('model_config')}</Typography.Text>}
      defaultVisible
      key={`${modelConfig?.model_id}-${promptInfo?.prompt_commit?.commit_info?.version}`}
    >
      <BasicModelConfigEditor
        value={modelConfig}
        onChange={config => {
          setModelConfig({ ...config });
        }}
        disabled={streaming || readonly}
        models={service.data?.models}
        onModelChange={setCurrentModel}
        modelSelectProps={{ className: 'w-full', loading: service.loading }}
      />
    </CollapseCard>
  );
}
