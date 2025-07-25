// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { getDataType } from '../util';
import { DataType, type DatasetItemProps } from '../type';
import { StringDatasetItem } from './string';
import { IntegerDatasetItem } from './integer';
import { FloatDatasetItem } from './float';
import { BoolDatasetItem } from './bool';

const TextColumnComponentMap = {
  [DataType.String]: StringDatasetItem,
  [DataType.Integer]: IntegerDatasetItem,
  [DataType.Boolean]: BoolDatasetItem,
  [DataType.Float]: FloatDatasetItem,
};

export const TextDatasetItem = (props: DatasetItemProps) => {
  const type = getDataType(props.fieldSchema);
  const Component = TextColumnComponentMap[type];

  return <Component {...props} />;
};
