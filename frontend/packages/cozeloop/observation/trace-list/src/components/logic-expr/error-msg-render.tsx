// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { type OptionProps } from '@coze-arch/coze-design';

import { checkValueIsEmpty } from './right-render';
import { I18n } from '@cozeloop/i18n-adapter';

interface ErrorMsgRenderProps {
  expr: {
    left?: string;
    right?: string | number | string[] | number[];
  };
  tagLeftOption: OptionProps[];
  checkIsInvalidateExpr: (expr: string) => boolean;
  valueChangeMap: Record<string, boolean>;
}

export const ErrorMsgRender = ({
  expr,
  tagLeftOption,
  checkIsInvalidateExpr,
  valueChangeMap,
}: ErrorMsgRenderProps) => {
  const { left = '', right } = expr;

  const isInvalidateExpr = checkIsInvalidateExpr(left ?? '');
  const leftname = tagLeftOption.find(item => item.value === left)?.label;

  if (isInvalidateExpr) {
    return (
      <div className="text-[#D0292F] text-[12px] whitespace-nowrap mt-1">
        {leftname ?? left} {I18n.t('filter_item_conflict')}
      </div>
    );
  }

  if (checkValueIsEmpty(right) && left && valueChangeMap[left]) {
    return (
      <div className="text-[#D0292F] text-[12px] whitespace-nowrap mt-1">
        {I18n.t('not_allowed_to_be_empty')}
      </div>
    );
  }

  return null;
};
