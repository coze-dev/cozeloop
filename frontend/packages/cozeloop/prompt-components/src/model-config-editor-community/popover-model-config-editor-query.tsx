// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useRequest } from 'ahooks';
import { useSpace } from '@cozeloop/biz-hooks-adapter';
import { Scenario } from '@cozeloop/api-schema/llm-manage';
import { LlmManageApi } from '@cozeloop/api-schema';

import {
  type ModelConfigPopoverProps,
  PopoverModelConfigEditor,
} from './popover-model-config-editor';

export function PopoverModelConfigEditorQuery(
  props: Omit<ModelConfigPopoverProps, 'models'>,
) {
  const { spaceID } = useSpace();

  const service = useRequest(async () => {
    const res = await LlmManageApi.ListModels({
      workspace_id: spaceID,
      page_size: 100,
      page_token: '0',
      scenario: Scenario.scenario_evaluator,
    });
    return { models: res?.models };
  });

  return <PopoverModelConfigEditor {...props} models={service.data?.models} />;
}
