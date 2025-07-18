// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useState } from 'react';

import { I18n } from '@cozeloop/i18n-adapter';
import { CodeEditor, handleCopy } from '@cozeloop/components';
import { IconCozCopy } from '@coze-arch/coze-design/icons';
import { Button, SemiSelect } from '@coze-arch/coze-design';

import styles from '../index.module.less';
import { type DatasetItemProps } from '../../../type';
import { useEditorLoading } from './use-editor-loading';
import { codeOptionsConfig, languageList } from './config';

export const CodeDatasetItem = ({
  fieldContent,
  onChange,
  isEdit,
}: DatasetItemProps) => {
  const [language, setLanguage] = useState('java');
  const { LoadingNode, onMount } = useEditorLoading();
  return (
    <div className={styles['code-container']}>
      {LoadingNode}
      <div className="flex items-center justify-between ">
        <SemiSelect
          zIndex={16000}
          size="small"
          optionList={languageList.map(item => ({
            label: item,
            value: item,
          }))}
          value={language}
          onChange={value => {
            setLanguage(value as string);
          }}
        />
        <Button
          icon={<IconCozCopy />}
          onClick={() => {
            handleCopy(fieldContent?.text || '');
          }}
          color="primary"
          size="small"
        >
          {I18n.t('copy')}
        </Button>
      </div>
      <div className="flex-1">
        <CodeEditor
          language={language}
          value={fieldContent?.text || ''}
          options={{
            readOnly: !isEdit,
            ...codeOptionsConfig,
          }}
          onMount={onMount}
          theme="vs-dark"
          onChange={value => {
            onChange?.({
              ...fieldContent,
              text: value,
            });
          }}
        />
      </div>
    </div>
  );
};
