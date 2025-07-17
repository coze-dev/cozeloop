// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { isObject } from 'lodash-es';
import classNames from 'classnames';
import { JsonViewer, type JsonViewerProps } from '@textea/json-viewer';
import { I18n } from '@cozeloop/i18n-adapter';
import { Tag } from '@coze-arch/coze-design';

import { JSON_VIEW_CONFIG } from '../consts/json-view';

import styles from './index.module.less';

export const PlantText = ({ content }: { content: string }) => (
  <span className={classNames(styles['view-string'], {})}>
    {content || '-'}
  </span>
);

export const renderPlainText = (
  content: string | object,
  config?: Partial<JsonViewerProps>,
) =>
  isObject(content) ? (
    <JsonViewer {...JSON_VIEW_CONFIG} {...(config ?? {})} value={content} />
  ) : (
    <PlantText content={content} />
  );

export const renderJsonContent = (
  content: string | object,
  config?: Partial<JsonViewerProps>,
) =>
  isObject(content) ? (
    <JsonViewer {...JSON_VIEW_CONFIG} {...(config ?? {})} value={content} />
  ) : (
    <>
      <Tag color="red" size="small" className="inline-block !w-fit">
        {I18n.t('invalid_json')}
      </Tag>
      <PlantText content={content} />
    </>
  );
