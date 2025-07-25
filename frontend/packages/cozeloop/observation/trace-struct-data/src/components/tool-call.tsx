// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { isEmpty } from 'lodash-es';
import classNames from 'classnames';

import { ReactComponent as IconSpanPluginTool } from '@/icons/icon-plugin-tool.svg';

import { safeJsonParse } from '../utils/json';
import { type RawMessage } from '../types';
import { renderJsonContent, renderPlainText } from './plain-text';

import styles from './index.module.less';

interface ToolCallProps {
  raw: RawMessage;
}
export const ToolCall = (props: ToolCallProps) => {
  const { raw } = props;

  if (isEmpty(raw.tool_calls)) {
    return null;
  }
  return (
    <div className="flex gap-2 flex-col">
      {raw.tool_calls?.map((tool, ind) => {
        const query = safeJsonParse(tool?.function?.arguments ?? '');
        return (
          <div key={ind} className="flex gap-2 flex-col">
            <div className={classNames(styles['tool-title'], 'font-mono')}>
              <IconSpanPluginTool style={{ width: '16px', height: '16px' }} />
              {tool?.function?.name || '-'}
            </div>
            {raw.role === 'tool' && raw.content
              ? renderPlainText(raw.content)
              : renderJsonContent(query)}
          </div>
        );
      })}
    </div>
  );
};
