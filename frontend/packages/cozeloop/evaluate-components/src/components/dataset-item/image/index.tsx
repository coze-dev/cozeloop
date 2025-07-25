// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { Image } from '@coze-arch/coze-design';

import { type DatasetItemProps } from '../type';
export const ImageDatasetItem: React.FC<DatasetItemProps> = ({
  fieldContent,
  expand,
  onChange,
}) => {
  const { image } = fieldContent || {};

  return (
    <Image
      className="inline-block"
      src={image?.url}
      alt={image?.name}
      width={36}
      height={36}
    />
  );
};
