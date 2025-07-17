// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { createPortal } from 'react-dom';
import React, {
  useDeferredValue,
  useTransition,
  useEffect,
  useState,
} from 'react';

import { isObject } from 'lodash-es';
import classNames from 'classnames';
import { JsonViewer } from '@textea/json-viewer';
import { Typography } from '@coze-arch/coze-design';

import { type Span, TagType } from '../types';
import { useFullScreen } from '../hooks/use-srcoll-view';
import { JSON_VIEW_CONFIG } from '../consts/json-view';
import { ViewAllModal } from './view-all';

import styles from './index.module.less';
import { I18n } from '@cozeloop/i18n-adapter';
interface RawContentProps {
  structuredContent: string | object;
  tagType?: TagType;
  className?: string;
  attrTos?: Span['attr_tos'];
}

export const RawContent: React.FC<RawContentProps> = ({
  structuredContent,
  tagType,
  className,
  attrTos,
}) => {
  const { isFullScreen } = useFullScreen();
  const [showModal, setShowModal] = useState(false);

  const handleViewAll = () => {
    setShowModal(true);
  };

  const showViewAllButton =
    (tagType === 'input' && attrTos?.input_data_url) ||
    (tagType === 'output' && attrTos?.output_data_url);
  return (
    <div
      className={classNames(
        styles['view-content'],
        'styled-scrollbar',
        className,
      )}
      style={
        isFullScreen
          ? {
              maxHeight: 'calc(100vh - 300px)',
            }
          : {}
      }
    >
      <div>
        {isObject(structuredContent) ? (
          <DeferredJSONViewer structuredContent={structuredContent} />
        ) : (
          <span
            className={classNames(styles['view-string'], {
              [styles.empty]: !structuredContent,
              '!text-[#ff441e]': tagType === TagType.Error,
            })}
          >
            {structuredContent || '-'}
          </span>
        )}
      </div>
      {showViewAllButton ? (
        <div className="inline-flex justify-end w-full pb-2">
          <Typography.Text
            className="!text-[rgb(var(--coze-up-brand-9))] text-xs leading-4 font-medium cursor-pointer"
            onClick={handleViewAll}
          >
            {I18n.t('view_all')}
          </Typography.Text>
        </div>
      ) : null}

      {showModal
        ? createPortal(
            <ViewAllModal
              onViewAllClick={setShowModal}
              tagType={tagType}
              attrTos={attrTos}
            />,
            document.getElementById(
              'trace-detail-side-sheet-panel',
            ) as HTMLDivElement,
          )
        : null}
    </div>
  );
};

function DeferredJSONViewer({ structuredContent }) {
  const deferredData = useDeferredValue(structuredContent);
  const [loading, setLoading] = useState(true);
  const [_, startTransition] = useTransition();

  useEffect(() => {
    startTransition(() => {
      setLoading(false);
    });
  }, []);
  if (loading) {
    return null;
  }
  return <JsonViewer value={deferredData} {...JSON_VIEW_CONFIG} />;
}
