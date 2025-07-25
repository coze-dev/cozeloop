// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type FieldSchema } from '@cozeloop/api-schema/evaluation';

import {
  ContentType,
  DataType,
  type ConvertFieldSchema,
} from '../components/dataset-item/type';

export const getDataType = (fieldSchema?: FieldSchema) => {
  try {
    const json = JSON.parse(fieldSchema?.text_schema || '{}');
    return json.type || DataType.String;
  } catch (error) {
    console.error(error);
    return DataType.String;
  }
};

export const convertSchemaToDataType = (
  schema: FieldSchema,
): ConvertFieldSchema => {
  if (schema.content_type !== ContentType.Text) {
    return schema;
  }
  return {
    ...schema,
    type: getDataType(schema),
  };
};

export const convertDataTypeToSchema = (
  data: ConvertFieldSchema,
): FieldSchema => {
  if (data.content_type !== ContentType.Text) {
    return data;
  }

  return {
    ...data,
    ...(data.type ? { text_schema: TYPE_CONFIG[data.type] } : {}),
  };
};

export const TYPE_CONFIG: Record<DataType, string> = {
  [DataType.Float]:
    '{"type":"number","multipleOf":0.0001,"maximum":1.7976931348623157e+308,"minimum":-1.7976931348623157e+308}',
  [DataType.Integer]:
    '{"type":"integer","maximum":9223372036854775807,"minimum":-9223372036854775808}',
  [DataType.String]: '{"type":"string"}',
  [DataType.Boolean]: '{"type":"boolean"}',
};
