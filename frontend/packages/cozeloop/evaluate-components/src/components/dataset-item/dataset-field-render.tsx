// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type Content } from '@cozeloop/api-schema/evaluation';

import { type DatasetFieldItemRenderProps } from './type';
import { DatasetItem } from './dataset-item-render';

export const DatasetFieldItemRender = ({
  fieldData,
  onChange,
  ...props
}: DatasetFieldItemRenderProps) => {
  const onContentChange = (content: Content) => {
    onChange?.({
      key: fieldData?.key,
      name: fieldData?.name,
      content: {
        ...fieldData?.content,
        ...content,
      },
    });
  };
  return (
    <DatasetItem
      {...props}
      onChange={onContentChange}
      fieldContent={fieldData?.content}
    />
  );
};
