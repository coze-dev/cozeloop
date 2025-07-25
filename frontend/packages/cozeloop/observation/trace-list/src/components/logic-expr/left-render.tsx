// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import classNames from 'classnames';
import { type OptionProps, Select } from '@coze-arch/coze-design';

import styles from './index.module.less';

interface LeftRendererProps {
  expr: {
    left?: string;
  };
  onExprChange?: (value: {
    left?: string;
    operator?: number;
    right?: string | number | string[] | number[];
  }) => void;
  tagLeftOption: OptionProps[];
  disabled?: boolean;
  defaultImmutableKeys?: string[];
  checkIsInvalidateExpr: (expr: string) => boolean;
}

export const LeftRenderer = ({
  expr,
  onExprChange,
  tagLeftOption,
  disabled,
  defaultImmutableKeys,
  checkIsInvalidateExpr,
}: LeftRendererProps) => {
  const { left } = expr;
  const isInvalidateExpr = checkIsInvalidateExpr(left ?? '');

  return (
    <div
      className={classNames(styles['expr-value-item-content'], {
        [styles['expr-value-item-content-invalidate']]: isInvalidateExpr,
      })}
    >
      <Select
        dropdownClassName={styles['render-select']}
        filter
        style={{ width: '100%' }}
        defaultOpen={!left}
        disabled={disabled || defaultImmutableKeys?.includes(left ?? '')}
        value={left}
        onChange={v => {
          const typedValue = v as string;
          onExprChange?.({
            left: typedValue,
            operator: undefined,
            right: undefined,
          });
        }}
        optionList={tagLeftOption}
      />
    </div>
  );
};
