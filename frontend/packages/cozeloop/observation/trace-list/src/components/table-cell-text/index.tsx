// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import {
  type PropsWithChildren,
  type ReactNode,
  type CSSProperties,
  type MouseEvent,
} from 'react';

import { I18n } from '@cozeloop/i18n-adapter';
import { handleCopy as copy } from '@cozeloop/components';
import { IconCozCopy } from '@coze-arch/coze-design/icons';
import { type PopoverProps, type TooltipProps } from '@coze-arch/coze-design';
import { Typography, IconButton, Tooltip } from '@coze-arch/coze-design';

export const handleCopy = (e: MouseEvent, text: string) => {
  e.stopPropagation();
  copy(text);
};

export const CustomTableTooltip = ({
  enableCopy,
  copyText,
  content,
  opts,
  style,
  children,
  textClassName,
  textAlign = 'left',
}: PropsWithChildren<{
  enableCopy?: boolean;
  copyText?: string;
  content?: ReactNode;
  style?: CSSProperties;
  opts?: Partial<PopoverProps> & Partial<TooltipProps>;
  textClassName?: string;
  textAlign?: 'left' | 'right';
}>) => {
  const enableTextCopy =
    enableCopy && children !== undefined && children !== '-';
  return (
    <div
      className={`flex items-center ${textAlign === 'left' ? 'justify-start' : 'justify-end'} gap-x-2 w-full`}
    >
      <Typography.Text
        ellipsis={{
          rows: 1,
          showTooltip: {
            type: 'popover',
            opts: {
              showArrow: false,
              stopPropagation: true,
              ...opts,
              style: {
                maxWidth: 500,
                maxHeight: 400,
                fontSize: 12,
                padding: 8,
                overflowY: 'auto',
                wordBreak: 'break-word',
                ...opts?.style,
              },
              content: content ?? children,
            },
          },
        }}
        style={{ fontSize: 13, ...style }}
        className={`text-[var(--coz-fg-plus)] max-w-full ${textClassName}`}
      >
        {children !== undefined ? children : '-'}
      </Typography.Text>
      {enableCopy ? (
        <Tooltip content={I18n.t('Copy')} position="top" theme="dark">
          <IconButton
            size="small"
            color="secondary"
            className="text-[var(--coz-fg-secondary)] !w-[24px] !h-[24px]"
            onClick={e => {
              if (enableTextCopy) {
                handleCopy(e, copyText || '-');
              }
            }}
            icon={
              <IconCozCopy className="w-[14px] h-[14px] text-[var(--coz-fg-secondary)]" />
            }
          />
        </Tooltip>
      ) : null}
    </div>
  );
};
