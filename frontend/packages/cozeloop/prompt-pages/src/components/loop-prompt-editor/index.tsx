// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { forwardRef, useImperativeHandle, useRef } from 'react';

import {
  type PromptBasicEditorRef,
  PromptEditor,
  type PromptEditorProps,
} from '@cozeloop/prompt-components';
import { I18n } from '@cozeloop/i18n-adapter';
import { handleCopy } from '@cozeloop/components';
import { Role } from '@cozeloop/api-schema/prompt';
import { IconCozCopy, IconCozTrashCan } from '@coze-arch/coze-design/icons';
import { IconButton, Popconfirm, Space } from '@coze-arch/coze-design';

interface LoopPromptEditorProps extends PromptEditorProps {
  onDelete?: (id?: Int64) => void;
  showDragBtn?: boolean;
  children?: React.ReactNode;
}

export const LoopPromptEditor = forwardRef<
  PromptBasicEditorRef,
  LoopPromptEditorProps
>(({ showDragBtn, messageTypeList, ...restProps }, ref) => {
  const editorRef = useRef<PromptBasicEditorRef>(null);

  useImperativeHandle(ref, () => ({
    setEditorValue: (value?: string) => {
      editorRef.current?.setEditorValue?.(value);
    },
    insertText: (text: string) => {
      editorRef.current?.insertText?.(text);
    },
  }));

  return (
    <>
      <PromptEditor
        ref={editorRef}
        rightActionBtns={
          <Space>
            <IconButton
              icon={<IconCozCopy />}
              color="secondary"
              size="mini"
              onClick={() => handleCopy(restProps.message?.content || '')}
            />
            {!restProps.onDelete ? null : (
              <Popconfirm
                title={I18n.t('delete_prompt_template')}
                content={I18n.t('confirm_delete_current_prompt_template')}
                cancelText={I18n.t('Cancel')}
                okText={I18n.t('delete')}
                okButtonProps={{ color: 'red' }}
                onConfirm={() =>
                  restProps.onDelete?.(
                    `${restProps.message?.key || restProps.message?.id || ''}`,
                  )
                }
              >
                <IconButton
                  icon={<IconCozTrashCan />}
                  color="secondary"
                  size="mini"
                  disabled={restProps.disabled}
                />
              </Popconfirm>
            )}
          </Space>
        }
        dragBtnHidden={!showDragBtn}
        messageTypeList={
          messageTypeList ?? [
            { label: 'System', value: Role.System },
            { label: 'Assistant', value: Role.Assistant },
            { label: 'User', value: Role.User },
            { label: 'Placeholder', value: Role.Placeholder },
          ]
        }
        {...restProps}
      >
        {restProps.children}
      </PromptEditor>
    </>
  );
});
