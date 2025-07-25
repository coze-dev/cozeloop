// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable complexity */
import React, { type ReactText } from 'react';

import { isEmpty } from 'lodash-es';
import classNames from 'classnames';
import { type FieldMeta } from '@cozeloop/api-schema/observation';
import {
  Input,
  Select,
  CozInputNumber,
  type SelectProps,
} from '@coze-arch/coze-design';

import {
  getOptionsWithKind,
  getOptionCopywriting,
  getLabelUnit,
} from './utils';
import { type CustomRightRenderMap } from './logic-expr';
import {
  EMPTY_RENDER_CMP_OP_LIST,
  type FilterFields,
  NUMBER_RENDER_CMP_OP_LIST,
  SELECT_MULTIPLE_RENDER_CMP_OP_LIST,
  SELECT_RENDER_CMP_OP_LIST,
} from './consts';

import styles from './index.module.less';

export const checkValueIsEmpty = (
  value?: number | number[] | string | string[] | null,
) =>
  isEmpty(value) ||
  (typeof value === 'string' && value.trim() === '') ||
  value === undefined ||
  value === null;

export interface RightRenderProps {
  left?: string;
  operator?: number | string;
  right?: string | number | string[] | number[];
  disabled?: boolean;
  defaultImmutableKeys?: string[];
  isInvalidateExpr?: boolean;
  valueChanged?: boolean;
  tagFilterRecord: Record<string, FieldMeta>;
  onRightValueChange?: (value: string | number | string[] | number[]) => void;
  onValueChangeStatus?: (fieldName: string, changed: boolean) => void;
  customRightRenderMap?: CustomRightRenderMap;
}

export const RightRender: React.FC<RightRenderProps> = props => {
  const {
    left = '',
    operator,
    right,
    disabled,
    defaultImmutableKeys,
    isInvalidateExpr,
    valueChanged,
    tagFilterRecord,
    onRightValueChange,
    onValueChangeStatus,
    customRightRenderMap,
  } = props;

  const { field_options, value_type, support_customizable_option } =
    tagFilterRecord[left] || {};

  const options = getOptionsWithKind({
    fieldOptions: field_options,
    valueKind: value_type,
  });

  if (
    !left ||
    !operator ||
    EMPTY_RENDER_CMP_OP_LIST.includes(String(operator))
  ) {
    return <div className={styles['expr-value-item-content']} />;
  }

  const multipleSelectProps: Partial<SelectProps> = {
    allowCreate: support_customizable_option,
    filter: support_customizable_option,
    multiple: true,
    maxTagCount: 4,
    ellipsisTrigger: true,
    showRestTagsPopover: true,
    restTagsPopoverProps: {
      position: 'top',
      stopPropagation: true,
    },
  };

  const showSelect =
    SELECT_MULTIPLE_RENDER_CMP_OP_LIST.includes(String(operator)) ||
    SELECT_RENDER_CMP_OP_LIST.includes(String(operator));

  const isMultiple = SELECT_MULTIPLE_RENDER_CMP_OP_LIST.includes(
    String(operator),
  );

  const isNumberInput = NUMBER_RENDER_CMP_OP_LIST.includes(
    left as FilterFields,
  );

  const customRightRender = customRightRenderMap?.[left];

  if (customRightRender) {
    return (
      <div
        className={classNames(styles['expr-value-item-content'], {
          [styles['expr-value-item-content-invalidate']]:
            isInvalidateExpr || (checkValueIsEmpty(right) && valueChanged),
        })}
      >
        {customRightRender?.({
          disabled: disabled || defaultImmutableKeys?.includes(left),
          style: { width: '100%' },
          value: right,
          onChange: v => {
            onRightValueChange?.(v as string[] | number[] | string | number);
            onValueChangeStatus?.(left, true);
          },
          optionList: options?.map(item => ({
            label: getOptionCopywriting(left, item),
            value: item,
          })),
          ...(isMultiple ? multipleSelectProps : {}),
        })}
      </div>
    );
  }

  return (
    <div
      className={classNames(styles['expr-value-item-content'], {
        [styles['expr-value-item-content-invalidate']]:
          isInvalidateExpr || (checkValueIsEmpty(right) && valueChanged),
      })}
    >
      {operator && showSelect ? (
        <Select
          dropdownClassName={styles['render-select']}
          disabled={disabled || defaultImmutableKeys?.includes(left)}
          style={{ width: '100%' }}
          value={right}
          onChange={v => {
            onRightValueChange?.(v as string[] | number[] | string | number);
            onValueChangeStatus?.(left, true);
          }}
          optionList={options?.map(item => ({
            label: getOptionCopywriting(left, item),
            value: item,
          }))}
          {...(isMultiple ? multipleSelectProps : {})}
        />
      ) : isNumberInput ? (
        <CozInputNumber
          formatter={v => `${v}`.replace(/\D/g, '')}
          disabled={disabled}
          hideButtons
          value={right as ReactText}
          max={Number.MAX_SAFE_INTEGER}
          min={Number.MIN_SAFE_INTEGER}
          onChange={v => {
            onRightValueChange?.(`${v}`.replace(/\D/g, ''));
            onValueChangeStatus?.(left, true);
          }}
          suffix={getLabelUnit(left)}
        />
      ) : (
        <Input
          disabled={disabled}
          value={right as ReactText}
          onChange={v => {
            onRightValueChange?.(v);
            onValueChangeStatus?.(left, true);
          }}
        />
      )}
    </div>
  );
};
