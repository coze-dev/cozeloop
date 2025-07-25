// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import classNames from 'classnames';
import {
  type EvalTargetVersion,
  type EvalTarget,
} from '@cozeloop/api-schema/evaluation';
import { Typography } from '@coze-arch/coze-design';

const ellipsis = {
  showTooltip: true,
};

export const getPromptEvalTargetOption = (
  item: EvalTarget,
  onlyShowOptionName?: boolean,
): { value?: string; label?: React.ReactNode } => {
  const etc = item.eval_target_version?.eval_target_content;
  const avatar = '';
  const title = etc?.prompt?.prompt_key || '';
  const subTitle = etc?.prompt?.name || '';

  return {
    value: item.source_target_id,
    label: onlyShowOptionName ? (
      <Typography.Text ellipsis={ellipsis}>{subTitle}</Typography.Text>
    ) : (
      <div className="flex flex-row items-center w-full overflow-hidden">
        {avatar ? (
          <img
            className="w-5 h-5 rounded-[4px] mr-2 flex-shrink-0"
            src={avatar}
          />
        ) : null}
        <Typography.Text
          className={'flex-shrink !max-w-[600px] text-[13px]'}
          ellipsis={ellipsis}
        >
          {title}
        </Typography.Text>
        <Typography.Text
          className={classNames(
            'flex-1 w-0 ml-3 text-xs font-medium coz-fg-secondary',
          )}
          ellipsis={ellipsis}
        >
          {subTitle}
        </Typography.Text>
      </div>
    ),
    ...item,
  };
};

export function getPromptEvalTargetVersionOption(item: EvalTargetVersion): {
  value?: string;
  label?: React.ReactNode;
} {
  return {
    value: item.source_target_version,
    label: (
      <div className="flex flex-row items-center w-full pr-2">
        <div className="flex-shrink-0 text-[13px] coz-fg-plus">
          {item.source_target_version}
        </div>
        <Typography.Text
          className="flex-1 w-0 ml-3 text-xs font-medium coz-fg-secondary"
          ellipsis={{
            showTooltip: true,
            rows: 1,
          }}
        >
          {item.eval_target_content?.prompt?.description}
        </Typography.Text>
      </div>
    ),
    ...item,
  };
}
