// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @typescript-eslint/no-explicit-any */
import { useRequest } from 'ahooks';
import { BaseSearchSelect } from '@cozeloop/components';
import { useBaseURL, useSpace } from '@cozeloop/biz-hooks-adapter';
import {
  EvalTargetType,
  type EvalTargetVersion,
} from '@cozeloop/api-schema/evaluation';
import { StoneEvaluationApi } from '@cozeloop/api-schema';
import { type FormSelect } from '@coze-arch/coze-design';

import { NoVersionJumper } from '../../common';
import { getPromptEvalTargetVersionOption } from './utils';
import { I18n } from '@cozeloop/i18n-adapter';

const PromptEvalTargetVersionSelect = ({
  promptId,
  ...props
}: React.ComponentProps<typeof FormSelect> & {
  promptId?: string;
}) => {
  const { spaceID } = useSpace();
  const { baseURL } = useBaseURL();

  const service = useRequest(
    async () => {
      if (!promptId) {
        return [];
      }
      const res = await StoneEvaluationApi.ListSourceEvalTargetVersions({
        workspace_id: spaceID,
        source_target_id: promptId,
        target_type: EvalTargetType.CozeLoopPrompt,
        page_size: 200,
      });

      const result: any[] =
        res.versions?.map(item => getPromptEvalTargetVersionOption(item)) || [];

      // 如果是 prompt 类型, 如果没有版本, 也需要提示去提交
      if (!res?.versions?.length) {
        result?.unshift({
          value: '__UNCOMMITTED__',
          label: (
            <NoVersionJumper targetUrl={`${baseURL}/pe/prompts/${promptId}`} />
          ),
          disabled: true,
        });
      }

      return result;
    },
    {
      refreshDeps: [promptId],
    },
  );

  const renderSelectedItem = (optionNode: any) => {
    const item: EvalTargetVersion = optionNode;
    return item.source_target_version;
  };

  return (
    <BaseSearchSelect
      loading={service.loading}
      emptyContent={I18n.t('no_data')}
      placeholder={I18n.t('please_select', { field: I18n.t('version') })}
      showRefreshBtn={true}
      onClickRefresh={() => service.run()}
      optionList={service.data}
      renderSelectedItem={renderSelectedItem}
      {...props}
    />
  );
};

export default PromptEvalTargetVersionSelect;
