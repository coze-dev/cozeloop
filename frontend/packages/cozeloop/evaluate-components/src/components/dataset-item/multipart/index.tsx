// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { Tag } from '@coze-arch/coze-design';

import { ContentType, type DatasetItemProps } from '../type';
import { StringDatasetItem } from '../text/string';
import { ImageDatasetItem } from '../image';
import { AudioDatasetItem } from '../audio';

const MultipartItemComponentMap = {
  [ContentType.Image]: ImageDatasetItem,
  [ContentType.Audio]: AudioDatasetItem,
  [ContentType.Text]: StringDatasetItem,
};

export const MultipartDatasetItem: React.FC<DatasetItemProps> = props => {
  const { fieldContent, expand } = props;
  const { multi_part } = fieldContent || {};

  return !expand ? (
    <SlimMultipartDatasetItem {...props} />
  ) : (
    <div className="flex flex-wrap gap-1 max-h-[292px] overflow-y-auto">
      {multi_part?.map((item, index) => {
        if (!item.content_type) {
          return;
        }
        const className =
          item.content_type === ContentType.Text
            ? 'w-full max-h-[auto] !border-0 !p-0'
            : '';
        const Component =
          MultipartItemComponentMap[item.content_type] || StringDatasetItem;
        return (
          <Component
            key={index}
            fieldContent={item}
            expand={true}
            displayFormat={true}
            className={className}
          />
        );
      })}
    </div>
  );
};
export const MAX_SHOW_ITEM = 4;

export const SlimMultipartDatasetItem: React.FC<DatasetItemProps> = props => {
  const { fieldContent } = props;
  const { multi_part } = fieldContent || {};
  if (multi_part?.length === 0) {
    return '';
  }
  if (multi_part?.[0]?.content_type === ContentType.Text) {
    return (
      <div className="flex items-center gap-1">
        <StringDatasetItem
          fieldContent={multi_part?.[0]}
          expand={false}
          isEdit={false}
          displayFormat={false}
        />
        {multi_part?.length > 1 && (
          <Tag size="mini" color="primary" className="min-w-[30px]">
            + {multi_part?.length - 1}
          </Tag>
        )}
      </div>
    );
  } else {
    let leftCount = 0;
    return (
      <div
        className="flex gap-1"
        onClick={e => {
          e.stopPropagation();
        }}
      >
        {multi_part?.map((item, index) => {
          if (!item.content_type) {
            return;
          }
          if (
            item.content_type === ContentType.Text ||
            index >= MAX_SHOW_ITEM
          ) {
            leftCount = multi_part.length - index;
            return;
          }
          const Component =
            MultipartItemComponentMap[item.content_type] || StringDatasetItem;
          return (
            <Component
              key={index}
              fieldContent={item}
              expand={true}
              displayFormat={true}
            />
          );
        })}

        {leftCount > 0 && (
          <Tag size="mini" color="primary">
            + {leftCount}
          </Tag>
        )}
      </div>
    );
  }
};
