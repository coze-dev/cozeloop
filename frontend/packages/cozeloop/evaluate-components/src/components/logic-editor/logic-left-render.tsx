// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import { type LeftRenderProps } from '@cozeloop/components';
import { Select } from '@coze-arch/coze-design';

import { dataTypeList, type RenderProps } from './logic-types';

export default function LeftRender(
  props: LeftRenderProps<string, string, string | number | undefined> &
    RenderProps,
) {
  const { expr, onExprChange, fields, disabled } = props;
  return (
    <div className="w-40">
      <Select
        placeholder={I18n.t('please_select', { field: '' })}
        value={expr.left}
        className="w-full"
        disabled={disabled}
        filter={true}
        optionList={fields.map(field => ({
          label: field.title,
          value: field.name,
        }))}
        onChange={val => {
          const field = fields.find(item => item.name === val);
          const { disabledOperations = [] } = field ?? {};
          const dataType = dataTypeList.find(item => item.type === field?.type);
          let operations = dataType?.operations ?? [];

          if (disabledOperations.length > 0) {
            operations = operations.filter(
              item => !disabledOperations.includes(item.value),
            );
          }
          onExprChange?.({
            left: val as string | undefined,
            operator: operations?.find(e => e.value === expr.operator)
              ? expr.operator
              : operations?.[0]?.value,
            right: undefined,
          });
        }}
      />
    </div>
  );
}
