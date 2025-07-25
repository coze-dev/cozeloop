// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useMemo } from 'react';

import {
  type EvaluationSetItem,
  type FieldSchema,
  type Turn,
} from '@cozeloop/api-schema/evaluation';
import { withField } from '@coze-arch/coze-design';

import { validarDatasetItem } from '../../dataset-item/util';
import { DatasetFieldItemRender } from '../../dataset-item/dataset-field-render';

interface DatasetItemRenderListProps {
  datasetItem?: EvaluationSetItem;
  turn?: Turn;
  fieldSchemas?: FieldSchema[];
  isEdit: boolean;
  fieldKey?: string;
}
const FormFieldItemRender = withField(DatasetFieldItemRender);

export const DatasetItemRenderList = ({
  fieldSchemas,
  isEdit,
  turn,
  fieldKey,
}: DatasetItemRenderListProps) => {
  const fieldSchemaMap = useMemo(() => {
    const map = new Map<string, FieldSchema>();
    fieldSchemas?.forEach(item => {
      if (item.key) {
        map.set(item.key, item);
      }
    });
    return map;
  }, [fieldSchemas]);
  return (
    <div className="flex flex-col">
      {turn?.field_data_list?.map((fieldData, index) => {
        const fieldSchema = fieldSchemaMap.get(fieldData.key || '');
        return (
          <FormFieldItemRender
            noLabel
            field={`${fieldKey}.field_data_list[${index}]`}
            key={fieldData?.key}
            fieldSchema={fieldSchema}
            fieldData={fieldData}
            showColumnKey
            expand={true}
            showEmpty={true}
            displayFormat={true}
            isEdit={isEdit}
            rules={[
              {
                validator: (_, value, callback) => {
                  if (
                    value?.content?.text === '' ||
                    value?.content?.text === undefined
                  ) {
                    return true;
                  }
                  const res = validarDatasetItem(
                    value?.content?.text,
                    callback,
                    fieldSchema,
                  );
                  return res;
                },
              },
            ]}
          />
        );
      })}
    </div>
  );
};
