// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable max-lines-per-function */
/* eslint-disable @coze-arch/max-line-per-function */
/* eslint-disable @typescript-eslint/no-magic-numbers */
/* eslint-disable complexity */

import { type ReactNode, useMemo } from 'react';

import { isEqual } from 'lodash-es';
import { useRequest } from 'ahooks';
import { PromptDiffEditor } from '@cozeloop/prompt-components';
import { useSpace } from '@cozeloop/biz-hooks-adapter';
import { type Prompt, ToolChoiceType } from '@cozeloop/api-schema/prompt';
import { Scenario } from '@cozeloop/api-schema/llm-manage';
import { LlmManageApi } from '@cozeloop/api-schema';
import {
  IconCozIllusEmpty,
  IconCozIllusEmptyDark,
} from '@coze-arch/coze-design/illustrations';
import { IconCozArrowRightFill } from '@coze-arch/coze-design/icons';
import { EmptyState, Tag, Typography } from '@coze-arch/coze-design';

import { objSortedKeys } from '@/utils/prompt';

import styles from './index.module.less';

export function DiffContent({
  base,
  current,
}: {
  base?: Prompt;
  current?: Prompt;
}) {
  const { spaceID } = useSpace();
  const { data } = useRequest(
    async () =>
      LlmManageApi.ListModels({
        workspace_id: spaceID,
        page_size: 100,
        page_token: '0',
        scenario: Scenario.scenario_prompt_debug,
      }),
    {
      manual: false,
    },
  );
  const { baseItem, currentItem, isSame } = useMemo(() => {
    const baseItemObj = {
      modelConfig: base?.prompt_commit?.detail?.model_config,
      variables: base?.prompt_commit?.detail?.prompt_template?.variable_defs,
      messageList: base?.prompt_commit?.detail?.prompt_template?.messages,
      tools: base?.prompt_commit?.detail?.tools,
      toolCallConfig: base?.prompt_commit?.detail?.tool_call_config,
    };
    const currentItemObj = {
      modelConfig: current?.prompt_draft?.detail?.model_config,
      variables: current?.prompt_draft?.detail?.prompt_template?.variable_defs,
      messageList: current?.prompt_draft?.detail?.prompt_template?.messages,
      tools: current?.prompt_draft?.detail?.tools,
      toolCallConfig: current?.prompt_draft?.detail?.tool_call_config,
    };

    return {
      baseItem: baseItemObj,
      currentItem: currentItemObj,
      isSame: isEqual(baseItemObj, currentItemObj),
    };
  }, [base, current]);

  const modelDiffData = useMemo(() => {
    if (isSame) {
      return [];
    }
    const array: { key: string; value: ReactNode }[] = [];
    const addDiffItem = (
      key: string,
      baseValue?: string | number,
      currentValue?: string | number,
    ) => {
      if (baseValue !== currentValue) {
        array.push({
          key,
          value: (
            <div className="flex items-center gap-4">
              <Tag color="primary">{key}</Tag>
              <Typography.Text className="flex gap-1 items-center !font-semibold">
                {baseValue ?? 'None'}
                <IconCozArrowRightFill />
                {currentValue ?? 'None'}
              </Typography.Text>
            </div>
          ),
        });
      }
    };

    addDiffItem(
      '模型 ID',
      baseItem.modelConfig?.model_id,
      currentItem.modelConfig?.model_id,
    );

    const baseModel = data?.models?.find(
      item => item.model_id === baseItem.modelConfig?.model_id,
    );
    const currentModel = data?.models?.find(
      item => item.model_id === currentItem.modelConfig?.model_id,
    );

    if (baseModel?.name !== currentModel?.name) {
      addDiffItem('模型名称', baseModel?.name || '', currentModel?.name || '');
    }

    addDiffItem(
      '生成随机性',
      baseItem.modelConfig?.temperature,
      currentItem.modelConfig?.temperature,
    );
    addDiffItem(
      '最大回复长度',
      baseItem.modelConfig?.max_tokens,
      currentItem.modelConfig?.max_tokens,
    );
    addDiffItem(
      'Top P',
      baseItem.modelConfig?.top_p,
      currentItem.modelConfig?.top_p,
    );
    addDiffItem(
      'JSON Mode',
      baseItem.modelConfig?.json_mode ? 'TRUE' : 'FALSE',
      currentItem.modelConfig?.json_mode ? 'TRUE' : 'FALSE',
    );

    return array;
  }, [isSame, baseItem, currentItem, data?.models?.length]);

  const templateIsSame = isEqual(baseItem.messageList, currentItem.messageList);

  const variabdlesDiffData = useMemo(() => {
    if (isSame) {
      return [];
    }
    const array: { key: string; value: ReactNode }[] = [];
    const deleteArray = baseItem.variables?.filter(
      item => !currentItem.variables?.find(it => it.key === item.key),
    );
    const addArray = currentItem.variables?.filter(
      item => !baseItem.variables?.find(it => it.key === item.key),
    );

    deleteArray?.forEach(item => {
      array.push({
        key: item.key || '',
        value: (
          <div className="flex items-center gap-4">
            <Tag color="primary">删除</Tag>
            <Typography.Text className="flex gap-1 items-center!font-semibold">
              {item.key}
            </Typography.Text>
          </div>
        ),
      });
    });

    addArray?.forEach(item => {
      array.push({
        key: item.key || '',
        value: (
          <div className="flex items-center gap-4">
            <Tag color="primary">新增</Tag>
            <Typography.Text className="flex gap-1 items-center!font-semibold">
              {item.key}
            </Typography.Text>
          </div>
        ),
      });
    });
    return array;
  }, [isSame, baseItem, currentItem]);

  const toolCallConfigIsSame = isEqual(
    baseItem.toolCallConfig?.tool_choice || ToolChoiceType.None,
    currentItem.toolCallConfig?.tool_choice || ToolChoiceType.None,
  );

  const toolsIsSame = isEqual(baseItem.tools, currentItem.tools);

  if (isSame) {
    return (
      <div className="w-full h-[433px] flex items-center justify-center">
        <EmptyState
          icon={<IconCozIllusEmpty width="160" height="160" />}
          darkModeIcon={<IconCozIllusEmptyDark width="160" height="160" />}
          title="本次提交无版本差异"
        />
      </div>
    );
  }
  return (
    <div className="w-full flex flex-col gap-5">
      {modelDiffData.length ? (
        <div className="flex flex-col gap-2">
          <Typography.Text className="!font-semibold">模型设置</Typography.Text>
          <div className={styles['diff-desc-table']}>
            {modelDiffData.map(it => (
              <div key={it.key} className={styles['diff-desc-table-row']}>
                {it.value}
              </div>
            ))}
          </div>
        </div>
      ) : null}
      {!templateIsSame ? (
        <div className="flex flex-col gap-2">
          <Typography.Text className="!font-semibold">
            Prompt 模板
          </Typography.Text>
          <div className={styles['diff-info-compare']}>
            <div className={styles['diff-info-compare-header']}>
              <Typography.Text
                strong
                size="small"
                className={styles['diff-info-compare-header-item']}
              >
                {base?.prompt_commit?.commit_info?.version}
              </Typography.Text>
              <Typography.Text
                className={styles['diff-info-compare-header-item']}
                size="small"
                strong
              >
                草稿
              </Typography.Text>
            </div>
            <div className="w-full h-[234px] overflow-auto styled-scrollbar !pr-[6px]">
              <PromptDiffEditor
                oldValue={JSON.stringify(
                  baseItem.messageList?.map(it =>
                    objSortedKeys({
                      ...it,
                      id: undefined,
                    }),
                  ) || [],
                  null,
                  2,
                )}
                newValue={JSON.stringify(
                  currentItem.messageList?.map(it =>
                    objSortedKeys({
                      ...it,
                      id: undefined,
                    }),
                  ) || [],
                  null,
                  2,
                )}
              />
            </div>
          </div>
        </div>
      ) : null}
      {variabdlesDiffData.length ? (
        <div className="flex flex-col gap-2">
          <Typography.Text className="!font-semibold">变量设置</Typography.Text>
          <div className={styles['diff-desc-table']}>
            {variabdlesDiffData.map(it => (
              <div key={it.key} className={styles['diff-desc-table-row']}>
                {it.value}
              </div>
            ))}
          </div>
        </div>
      ) : null}
      {!toolCallConfigIsSame || !toolsIsSame ? (
        <div className="flex flex-col gap-2">
          <Typography.Text className="!font-semibold">函数</Typography.Text>
          {toolCallConfigIsSame ? null : (
            <div className={styles['diff-desc-table']}>
              <div className={styles['diff-desc-table-row']}>
                <div className="flex items-center gap-4">
                  <Tag color="primary">函数</Tag>
                  <Typography.Text className="flex gap-1 items-center !font-semibold">
                    {baseItem.toolCallConfig?.tool_choice ===
                    ToolChoiceType.Auto
                      ? '打开 启用函数'
                      : '关闭 启用函数'}
                    <IconCozArrowRightFill />
                    {currentItem.toolCallConfig?.tool_choice ===
                    ToolChoiceType.Auto
                      ? '打开 启用函数'
                      : '关闭 启用函数'}
                  </Typography.Text>
                </div>
              </div>
            </div>
          )}
          {toolsIsSame ? null : (
            <div className={styles['diff-info-compare']}>
              <div className={styles['diff-info-compare-header']}>
                <Typography.Text
                  strong
                  size="small"
                  className={styles['diff-info-compare-header-item']}
                >
                  {base?.prompt_commit?.commit_info?.version}
                </Typography.Text>
                <Typography.Text
                  className={styles['diff-info-compare-header-item']}
                  size="small"
                  strong
                >
                  草稿
                </Typography.Text>
              </div>

              <div className="w-full h-[234px] overflow-auto styled-scrollbar !pr-[6px]">
                <PromptDiffEditor
                  oldValue={JSON.stringify(
                    baseItem.tools?.map(it => objSortedKeys(it)) || [],
                    null,
                    2,
                  )}
                  newValue={JSON.stringify(
                    currentItem.tools?.map(it => objSortedKeys(it)) || [],
                    null,
                    2,
                  )}
                />
              </div>
            </div>
          )}
        </div>
      ) : null}
    </div>
  );
}
