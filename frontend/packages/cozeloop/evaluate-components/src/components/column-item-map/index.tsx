// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import classNames from 'classnames';
import { type FieldSchema } from '@cozeloop/api-schema/evaluation';
import { IconCozEqual } from '@coze-arch/coze-design/icons';
import { Tag, type TooltipProps, Typography } from '@coze-arch/coze-design';

import { getColumnType } from '../dataset-item/util';
import { dataTypeMap } from '../dataset-item/type';

export function ReadonlyItem({
  title,
  value,
  typeText,
  className,
  showType = true,
  tooltipProps,
}: {
  title?: string;
  value?: React.ReactNode;
  typeText?: string;
  className?: string;
  showType?: boolean;
  tooltipProps?: Omit<TooltipProps, 'showArrow'>;
}) {
  return (
    <div
      className={classNames(
        'flex flex-row items-center h-8 gap-[8px] border border-solid coz-stroke-plus rounded-[6px] text-sm font-normal',
        className,
      )}
    >
      <div className="flex-shrink-0 coz-fg-secondary ml-3">{title}</div>
      <Typography.Text
        className="flex-1 !coz-fg-primary overflow-hidden"
        ellipsis={{
          showTooltip: {
            opts: {
              theme: 'dark',
              ...(tooltipProps ?? {}),
            },
          },
        }}
      >
        {value}
      </Typography.Text>
      {showType ? (
        <Tag className="flex-shrink-0 mr-3" size="mini" color="primary">
          {typeText}
        </Tag>
      ) : null}
    </div>
  );
}

export function EqualItem() {
  return (
    <div className="w-8 h-8 border border-solid coz-stroke-plus rounded-[6px] coz-fg-primary flex items-center justify-center shrink-0">
      <IconCozEqual className="w-4 h-4 coz-fg-primary" />
    </div>
  );
}

export function getTypeText(item?: FieldSchema) {
  return dataTypeMap[getColumnType(item) as keyof typeof dataTypeMap];
}
