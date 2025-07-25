// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable complexity */
import { forwardRef, type ReactNode, useState } from 'react';

import cn from 'classnames';
import { I18n } from '@cozeloop/i18n-adapter';
import { type Message, Role, VariableType } from '@cozeloop/api-schema/prompt';
import { IconCozHandle } from '@coze-arch/coze-design/icons';
import { IconButton, Input, Space } from '@coze-arch/coze-design';

import { VARIABLE_MAX_LEN } from '@/consts';

import {
  PromptBasicEditor,
  type PromptBasicEditorRef,
  type PromptBasicEditorProps,
} from '../basic-editor';
import { MessageTypeSelect } from './message-type-select';

import styles from './index.module.less';

export type PromptMessage<R = Role> = Omit<Message, 'role'> & {
  role?: R;
  id?: string;
  key?: string;
  optimize_key?: string;
};

type BasicEditorProps = Pick<
  PromptBasicEditorProps,
  | 'variables'
  | 'height'
  | 'minHeight'
  | 'maxHeight'
  | 'forbidJinjaHighlight'
  | 'forbidVariables'
  | 'linePlaceholder'
  | 'isGoTemplate'
>;

export interface PromptEditorProps<R extends string | number = Role>
  extends BasicEditorProps {
  className?: string;
  message?: PromptMessage<R>;
  dragBtnHidden?: boolean;
  messageTypeDisabled?: boolean;
  disabled?: boolean;
  isDrag?: boolean;
  placeholder?: string;
  messageTypeList?: Array<{ label: string; value: R }>;
  leftActionBtns?: ReactNode;
  rightActionBtns?: ReactNode;
  placeholderRoleValue?: R;
  onMessageChange?: (v: PromptMessage<R>) => void;
  onMessageTypeChange?: (v: R) => void;
  children?: ReactNode;
}

type PromptEditorType = <R extends string | number = Role>(
  props: PromptEditorProps<R> & {
    ref?: React.ForwardedRef<PromptBasicEditorRef>;
  },
) => JSX.Element;

export const PromptEditor = forwardRef(
  <R extends string | number = Role>(
    props: PromptEditorProps<R>,
    ref: React.ForwardedRef<PromptBasicEditorRef>,
  ) => {
    const {
      className,
      message,
      dragBtnHidden,
      messageTypeDisabled,
      variables,
      disabled,
      isDrag,
      onMessageChange,
      onMessageTypeChange,
      placeholder,
      messageTypeList,
      leftActionBtns,
      rightActionBtns,
      placeholderRoleValue = Role.Placeholder as R,
      children,
      ...rest
    } = props;
    const [editorActive, setEditorActive] = useState(false);
    const handleMessageContentChange = (v: string) => {
      onMessageChange?.({ ...message, content: v });
    };

    const readonly = disabled || isDrag;

    return (
      <div
        className={cn(
          styles['prompt-editor-container'],
          {
            [styles['prompt-editor-container-active']]: editorActive,
            [styles['prompt-editor-container-disabled']]: disabled,
          },
          className,
        )}
      >
        <div className={styles.header}>
          <Space spacing={4}>
            {dragBtnHidden ? null : (
              <IconButton
                color="secondary"
                size="mini"
                icon={<IconCozHandle fontSize={14} />}
                className={cn('drag !w-[14px]', styles.drag)}
              />
            )}
            {message?.role ? (
              <MessageTypeSelect<R>
                value={message.role}
                onChange={onMessageTypeChange}
                disabled={messageTypeDisabled || readonly}
                messageTypeList={messageTypeList}
              />
            ) : null}
            {leftActionBtns}
          </Space>
          {rightActionBtns}
        </div>
        <div
          className={cn('w-full', {
            'py-1': message?.role !== placeholderRoleValue,
          })}
        >
          {message?.role === placeholderRoleValue ? (
            <Input
              key={message.key || message.id}
              value={message.content}
              onChange={handleMessageContentChange}
              borderless
              disabled={readonly}
              style={{ border: 0, borderRadius: 0 }}
              onInput={event => {
                // 获取当前输入的值
                const target = event.target as HTMLInputElement;
                if (target) {
                  let { value } = target;
                  // 如果输入为空，不做处理
                  if (value === '') {
                    return;
                  }
                  // 确保首字母是字母
                  if (!/^[A-Za-z]/.test(value)) {
                    // 如果首字母不是字母，去掉首字母
                    value = value.slice(1);
                  }

                  // 确保其余部分只包含字母、数字和下划线
                  value = value.replace(/[^A-Za-z0-9_]/g, '');

                  // 更新输入框的值
                  target.value = value;
                }
              }}
              maxLength={VARIABLE_MAX_LEN}
              max={50}
              className="!pl-3 font-sm"
              inputStyle={{
                fontSize: 13,
                color: 'var(--Green-COZColorGreen7, #00A136)',
                fontFamily: 'JetBrainsMonoRegular',
              }}
              onFocus={() => setEditorActive(true)}
              onBlur={() => setEditorActive(false)}
              placeholder={I18n.t('prompt_var_format')}
            />
          ) : (
            <PromptBasicEditor
              key={message?.key || message?.id}
              {...rest}
              defaultValue={message?.content}
              onChange={handleMessageContentChange}
              variables={variables?.filter(
                it => it.type !== VariableType.Placeholder,
              )}
              readOnly={readonly}
              linePlaceholder={placeholder}
              onFocus={() => setEditorActive(true)}
              onBlur={() => setEditorActive(false)}
              ref={ref}
            >
              {children}
            </PromptBasicEditor>
          )}
        </div>
      </div>
    );
  },
) as PromptEditorType;
