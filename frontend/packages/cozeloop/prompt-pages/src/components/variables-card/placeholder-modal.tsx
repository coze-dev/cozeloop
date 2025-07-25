// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @coze-arch/max-line-per-function */
import { useEffect, useRef, useState } from 'react';

import { useShallow } from 'zustand/react/shallow';
import Sortable from 'sortablejs';
import { type PromptMessage } from '@cozeloop/prompt-components';
import { I18n } from '@cozeloop/i18n-adapter';
import { Role } from '@cozeloop/api-schema/prompt';
import { IconCozPlus } from '@coze-arch/coze-design/icons';
import { Button, Modal } from '@coze-arch/coze-design';

import { messageId } from '@/utils/prompt';
import { useBasicStore } from '@/store/use-basic-store';

import { LoopPromptEditor } from '../loop-prompt-editor';

interface PlaceholderModalProps {
  variableKey: string;
  visible: boolean;
  data?: PromptMessage[];
  onOk: (v: PromptMessage[]) => void;
  onCancel: () => void;
}

export function PlaceholderModal({
  variableKey,
  visible,
  data,
  onOk,
  onCancel,
}: PlaceholderModalProps) {
  const sortableContainer = useRef<HTMLDivElement>(null);
  const { streaming, readonly } = useBasicStore(
    useShallow(state => ({
      streaming: state.streaming,
      readonly: state.readonly,
    })),
  );

  const [messageList, setMessageList] = useState<Array<PromptMessage>>(
    data || [],
  );
  const [isDrag, setIsDrag] = useState(false);

  const addMessage = () => {
    const { length } = messageList;
    const id = messageId();
    const chat = {
      id,
      content: '',
    };
    if (length) {
      const message = messageList[length - 1];
      if (message?.role === Role.User) {
        setMessageList(list => [
          ...list,
          {
            ...chat,
            role: Role.Assistant,
          },
        ]);
      } else {
        setMessageList(list => [
          ...list,
          {
            ...chat,
            role: Role.User,
          },
        ]);
      }
    } else {
      setMessageList(list => [
        ...list,
        {
          ...chat,
          role: Role.User,
        },
      ]);
    }
  };

  const handleOk = () => {
    onOk?.(messageList);
  };

  useEffect(() => {
    if (sortableContainer.current && visible) {
      new Sortable(sortableContainer.current, {
        animation: 150,
        handle: '.drag',
        onSort: evt => {
          setMessageList(list => {
            const draft = [...(list ?? [])];
            if (draft.length) {
              const { oldIndex = 0, newIndex = 0 } = evt;
              const [item] = draft.splice(oldIndex, 1);
              draft.splice(newIndex, 0, item);
            }
            return draft;
          });
        },
        onStart: () => setIsDrag(true),
        onEnd: () => setIsDrag(false),
      });
    }
  }, [visible]);

  useEffect(() => {
    if (visible) {
      setMessageList(data || []);
    } else {
      setMessageList([]);
    }
  }, [visible, data?.map(v => v?.content).join('')]);

  return (
    <Modal
      title={I18n.t('mock_message_group', { key: variableKey })}
      visible={visible}
      onCancel={onCancel}
      width={920}
      cancelText={I18n.t('Cancel')}
      okText={I18n.t('confirm')}
      onOk={handleOk}
    >
      <div className="flex flex-col gap-2 h-[500px] overflow-y-auto">
        <div className="flex flex-col gap-2 w-full" ref={sortableContainer}>
          {messageList?.map(message => (
            <LoopPromptEditor
              key={message.id}
              message={message}
              disabled={readonly || streaming}
              isDrag={isDrag}
              onDelete={key =>
                setMessageList(prev => {
                  const newInfo = prev?.filter(it => it.id !== key);
                  return newInfo;
                })
              }
              onMessageTypeChange={v =>
                setMessageList(prev => {
                  const newInfo = prev?.map(it =>
                    it.id === message.id ? { ...it, role: v } : it,
                  );
                  return newInfo;
                })
              }
              onMessageChange={v =>
                setMessageList(prev => {
                  const newInfo = prev?.map(it =>
                    it.id === v.id ? { ...it, ...v } : it,
                  );
                  return newInfo;
                })
              }
              messageTypeList={[
                { label: 'System', value: Role.System },
                { label: 'Assistant', value: Role.Assistant },
                { label: 'User', value: Role.User },
              ]}
              minHeight={26}
              maxHeight={240}
              forbidJinjaHighlight
              forbidVariables
              placeholder={I18n.t('please_input', {
                field: I18n.t('mock_message'),
              })}
            />
          ))}
        </div>
        <Button
          className="flex-shrink-0 w-[fit-content]"
          icon={<IconCozPlus />}
          onClick={addMessage}
          disabled={streaming || readonly}
          color="primary"
        >
          {I18n.t('add_message')}
        </Button>
      </div>
    </Modal>
  );
}
