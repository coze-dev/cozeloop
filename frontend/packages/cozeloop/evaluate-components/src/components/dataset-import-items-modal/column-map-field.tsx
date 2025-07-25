// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import { TooltipWhenDisabled } from '@cozeloop/components';
import { type FieldSchema } from '@cozeloop/api-schema/evaluation';
import { type FieldMapping } from '@cozeloop/api-schema/data';
import { Select, Typography } from '@coze-arch/coze-design';

import { EqualItem, getTypeText, ReadonlyItem } from '../column-item-map';

interface FieldMappingConvert extends FieldMapping {
  description?: string;
  fieldSchema?: FieldSchema;
}
interface ColumnMapFieldProps {
  sourceColumns: string[];
  value?: FieldMappingConvert[];
  onChange?: (value: FieldMappingConvert[]) => void;
}

export const ColumnMapField = ({
  sourceColumns,
  onChange,
  value,
}: ColumnMapFieldProps) => (
  <div className="flex flex-col gap-3 mt-3">
    {value?.map((item, index) => (
      <div key={index} className="flex gap-2">
        <TooltipWhenDisabled
          content={item?.description}
          disabled={!!item?.description}
          theme="dark"
        >
          <div>
            <ReadonlyItem
              className="w-[276px] overflow-hidden"
              title={I18n.t('evaluation_set_column')}
              typeText={getTypeText(item?.fieldSchema)}
              value={item.target}
            />
          </div>
        </TooltipWhenDisabled>
        <EqualItem />
        <Select
          prefix={
            <Typography.Text className="!coz-fg-secondary ml-3">
              {I18n.t('import_data_column')}
            </Typography.Text>
          }
          className="!w-[276px]"
          optionList={sourceColumns.map(column => ({
            label: column,
            value: column,
          }))}
          showClear
          value={item.source}
          onChange={newTarget => {
            const newValues = [...value];
            newValues[index].source = newTarget as string;
            onChange?.(newValues);
          }}
        ></Select>
      </div>
    ))}
  </div>
);
