// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useEffect, useState } from 'react';

import { isObject } from 'lodash-es';
import classNames from 'classnames';
import { type JsonViewerProps } from '@textea/json-viewer';
import { I18n } from '@cozeloop/i18n-adapter';
import { handleCopy as copy } from '@cozeloop/components';
import { IconCozCopy, IconCozCorner } from '@coze-arch/coze-design/icons';
import { Button, SegmentTab, Tooltip } from '@coze-arch/coze-design';

import { beautifyJson, safeJsonParse } from '../utils/json';
import { type Span, type TagType } from '../types';
import { FullscreenContext, useScrollView } from '../hooks/use-srcoll-view';

import styles from './index.module.less';

export enum TypeEnum {
  TEXT = 'TEXT',
  JSON = 'JSON',
}
interface SpanContentContainerProps {
  content: string | object;
  title: string;
  config?: Partial<JsonViewerProps>;
  hasBottomLine?: boolean;
  canSwitchRawType?: boolean;
  spanType?: string;
  tagType?: TagType;
  copyConfig: {
    moduleName?: string;
    point: string;
  };
  isEncryptionData?: boolean;
  spanID?: string;
  attrTos?: Span['attr_tos'];
  children: (
    renderType: TypeEnum,
    structuredContent: string | object,
  ) => React.ReactNode;
  hideSwitchRawType?: boolean;
}
export const SpanContentContainer = (props: SpanContentContainerProps) => {
  const [showType, setShowType] = useState<TypeEnum>(TypeEnum.TEXT);
  const {
    content,
    title,
    hasBottomLine = true,
    spanID,
    children,
    hideSwitchRawType = false,
  } = props;

  const { containerRef, isFullScreen, onFullScreenStateChange } =
    useScrollView();
  const handleCopy = (data: object | string) => {
    let str = '';
    if (isObject(data)) {
      str = beautifyJson(data);
    } else {
      str = data;
    }
    copy(str);
  };
  const structuredContent =
    typeof content === 'string' ? safeJsonParse(content) : content;
  useEffect(() => {
    setShowType(TypeEnum.TEXT);
    onFullScreenStateChange(false);
  }, [spanID]);

  return (
    <FullscreenContext.Provider value={{ isFullScreen }}>
      <div
        ref={containerRef}
        style={{
          borderBottom: hasBottomLine ? '1px solid #1D1C2314' : 'none',
        }}
        className={classNames('flex flex-col items-stretch px-[20px] py-3')}
      >
        <div className="flex items-center align-self-stretch justify-between h-8 mb-2">
          <div className="flex gap-1 items-center text-[16px] font-medium leading-[20px] text-[#000000]">
            <span className="mr-1">{title}</span>
            {structuredContent ? (
              <Tooltip content={I18n.t('Copy')} theme="dark">
                <Button
                  className="!w-[24px] !h-[24px] box-border mr-1"
                  size="small"
                  color="secondary"
                  icon={
                    <IconCozCopy className="flex items-center justify-center w-[14px] h-[14px] text-[var(--coz-fg-secondary)]" />
                  }
                  onClick={() => {
                    handleCopy(structuredContent);
                  }}
                />
              </Tooltip>
            ) : null}
            {showType === TypeEnum.JSON ? (
              <Tooltip
                content={isFullScreen ? I18n.t('collapse') : I18n.t('expand')}
                theme="dark"
              >
                <Button
                  size="small"
                  className="!w-[24px] !h-[24px] box-border"
                  color={isFullScreen ? 'primary' : 'secondary'}
                  icon={<IconCozCorner className={styles['copy-icon']} />}
                  onClick={() => {
                    onFullScreenStateChange(!isFullScreen);
                  }}
                />
              </Tooltip>
            ) : null}
          </div>
          <div>
            {isObject(structuredContent) && !hideSwitchRawType && (
              <SegmentTab
                className={styles['segment-tab']}
                value={showType}
                size="small"
                onChange={event => {
                  setShowType(event.target.value as unknown as TypeEnum);
                }}
                options={[
                  { label: 'TEXT', value: TypeEnum.TEXT },
                  { label: 'JSON', value: TypeEnum.JSON },
                ]}
              />
            )}
          </div>
        </div>
        {children?.(showType, structuredContent)}
      </div>
    </FullscreenContext.Provider>
  );
};
