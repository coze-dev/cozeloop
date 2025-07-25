// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type CSSProperties } from 'react';

import classNames from 'classnames';
import { I18n } from '@cozeloop/i18n-adapter';
import { IconCozCopy } from '@coze-arch/coze-design/icons';
import { IconButton, Tooltip, Typography } from '@coze-arch/coze-design';

import { handleCopy } from '../utils/basic';

interface TextWithCopyProps {
  content?: string;
  displayText?: string;
  copyTooltipText?: string;
  maxWidth?: number | string;
  className?: string;
  style?: CSSProperties;
  textClassName?: string;
  textType?:
    | 'success'
    | 'secondary'
    | 'primary'
    | 'danger'
    | 'warning'
    | 'tertiary'
    | 'quaternary';
}

export function TextWithCopy({
  displayText,
  copyTooltipText,
  content,
  className,
  maxWidth,
  style,
  textClassName,
  textType = 'secondary',
}: TextWithCopyProps) {
  return (
    <div
      className={classNames('flex items-center justify-start gap-1', className)}
      style={style}
    >
      <Typography.Text
        className={classNames('max-w-full', textClassName)}
        type={textType}
        style={{ maxWidth }}
        ellipsis={{
          showTooltip: { opts: { theme: 'dark', content } },
        }}
        onClick={e => {
          content && handleCopy(content);
          e?.stopPropagation();
        }}
      >
        {displayText || content || ''}
      </Typography.Text>
      {content ? (
        <Tooltip
          content={copyTooltipText || I18n.t('copy_content')}
          theme="dark"
        >
          <IconButton
            size="mini"
            color="secondary"
            className="flex-shrink-0 !w-[20px] !h-[20px]"
            icon={
              <IconCozCopy
                className=""
                onClick={e => {
                  content && handleCopy(content);
                  e?.stopPropagation();
                }}
                fontSize={14}
                fill="var(--semi-color-text-2)"
              />
            }
          />
        </Tooltip>
      ) : null}
    </div>
  );
}
