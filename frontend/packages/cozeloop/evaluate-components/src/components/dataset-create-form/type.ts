// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import { FieldDisplayFormat } from '@cozeloop/api-schema/data';

import {
  ContentType,
  type ConvertFieldSchema,
  DataType,
} from '../dataset-item/type';

export interface IDatasetCreateForm {
  name?: string;
  columns?: ConvertFieldSchema[];
  description?: string;
}

export const DEFAULT_DATASET_CREATE_FORM: IDatasetCreateForm = {
  name: '',
  columns: [
    {
      name: 'input',
      content_type: ContentType.Text,
      type: DataType.String,
      default_display_format: FieldDisplayFormat.PlainText,
      description: I18n.t('input_for_evaluation_object'),
    },
    {
      name: 'reference_output',
      content_type: ContentType.Text,
      type: DataType.String,
      default_display_format: FieldDisplayFormat.PlainText,
      description: I18n.t('expected_ideal_output'),
    },
  ],
  description: '',
};

export const DEFAULT_COLUMN_SCHEMA: ConvertFieldSchema = {
  name: '',
  content_type: ContentType.Text,
  type: DataType.String,
  default_display_format: FieldDisplayFormat.PlainText,
};
