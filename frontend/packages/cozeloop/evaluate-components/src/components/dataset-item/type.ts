// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import {
  type FieldData,
  type Content,
  type FieldSchema,
  ContentType,
} from '@cozeloop/api-schema/evaluation';
import { FieldDisplayFormat } from '@cozeloop/api-schema/data';
export interface DatasetItemBaseProps {
  fieldSchema?: FieldSchema;
  isEdit?: boolean;
  /**
   * 是否显示列名
   */
  showColumnKey?: boolean;
  /**
   * 是否展开
   */
  expand?: boolean;
  /**
   * 是否显示md,code,json的渲染格式
   */
  displayFormat?: boolean;

  /**
   * className
   */
  className?: string;

  /**
   * 是否显示空状态图标，默认不展示。
   */
  showEmpty?: boolean;
}

export interface DatasetFieldItemRenderProps extends DatasetItemBaseProps {
  fieldData?: FieldData;
  onChange?: (fieldData: FieldData) => void;
}
export interface DatasetItemProps extends DatasetItemBaseProps {
  className?: string;
  fieldContent?: Content;
  onChange?: (content: Content) => void;
}
export enum DataType {
  String = 'string',
  Integer = 'integer',
  Float = 'number',
  Boolean = 'boolean',
}
export const dataTypeMap = {
  [DataType.String]: 'String',
  [DataType.Integer]: 'Integer',
  [DataType.Float]: 'Float',
  [DataType.Boolean]: 'Boolean',
};
export const displayFormatType = {
  [FieldDisplayFormat.PlainText]: 'PlainText',
  [FieldDisplayFormat.Code]: 'Code',
  [FieldDisplayFormat.JSON]: 'JSON',
  [FieldDisplayFormat.Markdown]: 'Markdown',
};

export { ContentType };

export const COLUMN_TYPE_MAP = {
  [ContentType.Text]: 'Text',
  [ContentType.Audio]: 'Audio',
  [ContentType.Image]: 'Image',
  [ContentType.MultiPart]: 'MultiPart',
};

export const DATA_TYPE_LIST = [
  {
    label: 'String',
    value: DataType.String,
  },
  {
    label: 'Integer',
    value: DataType.Integer,
  },
  {
    label: 'Float',
    value: DataType.Float,
  },
  {
    label: 'Boolean',
    value: DataType.Boolean,
  },
];

export const DISPLAY_TYPE_MAP = {
  [DataType.String]: [
    FieldDisplayFormat.PlainText,
    FieldDisplayFormat.Code,
    FieldDisplayFormat.JSON,
    FieldDisplayFormat.Markdown,
  ],
  [DataType.Integer]: [FieldDisplayFormat.PlainText],
  [DataType.Float]: [FieldDisplayFormat.PlainText],
  [DataType.Boolean]: [FieldDisplayFormat.PlainText],
};
export const DISPLAY_FORMAT_MAP = {
  [FieldDisplayFormat.PlainText]: 'PlainText',
  [FieldDisplayFormat.Code]: 'Code',
  [FieldDisplayFormat.JSON]: 'JSON',
  [FieldDisplayFormat.Markdown]: 'Markdown',
};

export interface ConvertFieldSchema extends FieldSchema {
  type?: DataType;
}
