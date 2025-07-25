// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @typescript-eslint/no-explicit-any */
import { useMemo, useState } from 'react';

import { LogicExpr, type LogicOperator } from '@cozeloop/components';
import { type FieldMeta } from '@cozeloop/api-schema/observation';
import { type OptionProps, type SelectProps } from '@coze-arch/coze-design';

import {
  formatExprValue,
  formatSpanFilterValue,
  getKeyCopywriting,
} from './utils';
import { RightRender } from './right-render';
import { OperatorRenderer } from './operator-render';
import { LeftRenderer } from './left-render';
import { ErrorMsgRender } from './error-msg-render';

import styles from './index.module.less';

type LogicExprTypes = [
  string | undefined,
  number | undefined | string,
  string | number | string[] | number[] | undefined,
];

// 自定义右侧渲染器的映射类型
export type CustomRightRenderMap = Record<
  string,
  (props?: SelectProps & Record<string, any>) => React.ReactNode
>;

export interface LogicValue {
  filter_fields?: LogicItem[];
  query_and_or?: string;
  sub_filter?: Array<LogicValue>;
}

export interface LogicItem {
  field_name: string;
  query_type: string;
  values: string[];
}

export interface AnalyticsLogicExprProps {
  disabled?: boolean;
  value?: LogicValue;
  disableDuplicateSelect?: boolean;
  defaultImmutableKeys?: string[];
  tagFilterRecord: Record<string, FieldMeta>;
  onChange?: (value?: any) => void;
  allowLogicOperators?: LogicOperator[];
  invalidateExpr?: Set<string>;
  // 新增的自定义渲染器
  customRightRenderMap?: CustomRightRenderMap;
}

// Helper function to sort strings by first character (letters first)
const sortByFirstChar = (a: string, b: string): number => {
  const isLetter = (char: string) => /^[A-Za-z]$/.test(char);
  const aIsLetter = isLetter(a.charAt(0));
  const bIsLetter = isLetter(b.charAt(0));

  if (aIsLetter && bIsLetter) {
    return a.localeCompare(b, undefined, { sensitivity: 'base' });
  }
  return aIsLetter ? -1 : bIsLetter ? 1 : 0;
};

export const AnalyticsLogicExpr = (props: AnalyticsLogicExprProps) => {
  const {
    value,
    tagFilterRecord,
    disableDuplicateSelect,
    defaultImmutableKeys,
    onChange,
    disabled,
    allowLogicOperators = ['and'],
    invalidateExpr = new Set(),
    customRightRenderMap = {},
  } = props;

  const exprValue = useMemo(
    () =>
      formatExprValue<LogicExprTypes[0], LogicExprTypes[1], LogicExprTypes[2]>(
        value,
        tagFilterRecord,
        defaultImmutableKeys,
      ),
    [value, defaultImmutableKeys, tagFilterRecord],
  );

  const checkIsInvalidateExpr = (expr: string) => invalidateExpr.has(expr);
  const [valueChangeMap, setValueChangeMap] = useState<Record<string, boolean>>(
    {},
  );

  const { tagLeftOption } = useMemo<{
    tagLeftOption: OptionProps[];
  }>(() => {
    const selectedItemKeyList = exprValue?.exprs?.map((item: any) => item.left);
    return {
      tagLeftOption: Object.keys(tagFilterRecord)
        .sort((a, b) => sortByFirstChar(a, b))
        .map(key => ({
          label: getKeyCopywriting(key),
          value: key,
          disabled:
            disableDuplicateSelect && selectedItemKeyList?.includes(key),
        })),
    };
  }, [exprValue, disableDuplicateSelect]);

  const handleValueChangeStatus = (fieldName: string, changed: boolean) => {
    setValueChangeMap(prev => ({
      ...prev,
      [fieldName]: changed,
    }));
  };

  return (
    <LogicExpr<LogicExprTypes[0], LogicExprTypes[1], LogicExprTypes[2]>
      value={exprValue}
      readonly={disabled}
      allowLogicOperators={allowLogicOperators}
      onDeleteExpr={key => {
        setValueChangeMap(prev => ({
          ...prev,
          [key as string]: false,
        }));
      }}
      exprGroupRenderContentItemsClassName={
        styles['expr-group-render-content-items']
      }
      leftRender={leftRenderProps => (
        <LeftRenderer
          expr={leftRenderProps.expr}
          onExprChange={leftRenderProps.onExprChange}
          tagLeftOption={tagLeftOption}
          disabled={disabled}
          defaultImmutableKeys={defaultImmutableKeys}
          checkIsInvalidateExpr={checkIsInvalidateExpr}
        />
      )}
      operatorRender={operatorRenderProps => (
        <OperatorRenderer
          expr={operatorRenderProps.expr}
          onExprChange={operatorRenderProps.onExprChange}
          tagFilterRecord={tagFilterRecord}
          disabled={disabled}
          defaultImmutableKeys={defaultImmutableKeys}
          checkIsInvalidateExpr={checkIsInvalidateExpr}
        />
      )}
      rightRender={rightRenderProps => {
        const {
          expr: { left = '', operator, right },
          onChange: onRightValueChange,
        } = rightRenderProps;

        const isInvalidateExpr = checkIsInvalidateExpr(left);

        return (
          <RightRender
            left={left}
            operator={operator}
            right={right}
            disabled={disabled}
            defaultImmutableKeys={defaultImmutableKeys}
            isInvalidateExpr={isInvalidateExpr}
            valueChanged={valueChangeMap[left]}
            tagFilterRecord={tagFilterRecord}
            onRightValueChange={onRightValueChange}
            onValueChangeStatus={handleValueChangeStatus}
            customRightRenderMap={customRightRenderMap}
          />
        );
      }}
      errorMsgRender={errorMsgRenderProps => (
        <ErrorMsgRender
          expr={errorMsgRenderProps.expr}
          tagLeftOption={tagLeftOption}
          checkIsInvalidateExpr={checkIsInvalidateExpr}
          valueChangeMap={valueChangeMap}
        />
      )}
      maxNestingDepth={1}
      defaultExpr={{
        left: undefined,
        operator: undefined,
        right: undefined,
      }}
      onChange={expr => {
        onChange?.(formatSpanFilterValue(expr, tagFilterRecord));
      }}
    />
  );
};
