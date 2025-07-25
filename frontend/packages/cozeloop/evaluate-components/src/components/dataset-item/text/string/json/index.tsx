// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { isObject } from 'lodash-es';
import { JsonViewer } from '@textea/json-viewer';
import { safeJsonParse } from '@cozeloop/toolkit';
import { CodeEditor } from '@cozeloop/components';

import { PlainTextDatasetItemReadOnly } from '../plain-text/readonly';
import styles from '../index.module.less';
import { useEditorLoading } from '../code/use-editor-loading';
import { codeOptionsConfig } from '../code/config';
import { type DatasetItemProps } from '../../../type';
import { jsonViewerConfig } from './config';
export const JSONDatasetItem = (props: DatasetItemProps) => {
  const { fieldContent, onChange, isEdit } = props;
  const { LoadingNode, onMount } = useEditorLoading();
  const jsonObject = safeJsonParse(fieldContent?.text || '');
  return isEdit ? (
    <div className={styles['code-container']}>
      {LoadingNode}
      <CodeEditor
        language={'json'}
        value={fieldContent?.text || ''}
        options={{
          readOnly: !isEdit,
          ...codeOptionsConfig,
        }}
        theme="vs-dark"
        onMount={onMount}
        onChange={value => {
          onChange?.({
            ...fieldContent,
            text: value,
          });
        }}
      />
    </div>
  ) : isObject(jsonObject) ? (
    <div className={styles['code-container-readonly']}>
      <JsonViewer {...jsonViewerConfig} value={jsonObject} />
    </div>
  ) : (
    <PlainTextDatasetItemReadOnly {...props} />
  );
};
