// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useEffect, useRef, useState } from 'react';

import { useShallow } from 'zustand/react/shallow';
import Sortable from 'sortablejs';
import { nanoid } from 'nanoid';
import classNames from 'classnames';
import { CollapseCard } from '@cozeloop/components';
import { Role } from '@cozeloop/api-schema/prompt';
import { IconCozPlus } from '@coze-arch/coze-design/icons';
import { Button, Typography } from '@coze-arch/coze-design';

import { useBasicStore } from '@/store/use-basic-store';
import { useCompare } from '@/hooks/use-compare';

import { LoopPromptEditor } from '../loop-prompt-editor';
import { I18n } from '@cozeloop/i18n-adapter';

interface PromptEditorCardProps {
  uid?: number;
  canCollapse?: boolean;
  defaultVisible?: boolean;
}

export function PromptEditorCard({
  canCollapse,
  defaultVisible,
  uid,
}: PromptEditorCardProps) {
  const sortableContainer = useRef<HTMLDivElement>(null);
  const { streaming, messageList, setMessageList, variables } = useCompare(uid);

  const { readonly: basicReadonly } = useBasicStore(
    useShallow(state => ({
      readonly: state.readonly,
    })),
  );
  const [isDrag, setIsDrag] = useState(false);
  const readonly = basicReadonly || streaming;

  const handleAddMessage = () => {
    let messageType = Role.User;
    setMessageList(prev => {
      if (!prev?.length) {
        messageType = Role.System;
      } else if (prev?.[prev.length - 1]?.role === Role.User) {
        messageType = Role.Assistant;
      }
      const newInfo = (prev || [])?.concat({
        key: nanoid(),
        role: messageType,
        content: '',
      });
      return newInfo;
    });
  };

  useEffect(() => {
    if (sortableContainer.current) {
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
  }, []);

  return (
    <CollapseCard
      title={
        <Typography.Text strong>{I18n.t('prompt_template')}</Typography.Text>
      }
      defaultVisible={defaultVisible}
      disableCollapse={!canCollapse}
    >
      <div
        className={classNames('flex flex-col gap-2', {
          'pt-3': canCollapse,
        })}
      >
        <div className={'flex flex-col gap-2'} ref={sortableContainer}>
          {messageList
            ?.filter(it => Boolean(it))
            ?.map(message => (
              <LoopPromptEditor
                key={message.key}
                message={message}
                variables={variables}
                disabled={readonly}
                isDrag={isDrag}
                onMessageTypeChange={v =>
                  setMessageList(prev => {
                    const newInfo = prev?.map(it => {
                      if (it.key === message.key) {
                        if (
                          it.role === Role.Placeholder ||
                          v === Role.Placeholder
                        ) {
                          return {
                            ...it,
                            role: v,
                            content: '',
                            key: nanoid(),
                          };
                        }
                        return { ...it, role: v };
                      }
                      return it;
                    });
                    return newInfo;
                  })
                }
                onMessageChange={v =>
                  setMessageList(prev => {
                    const newInfo = prev?.map(it =>
                      it.key === v.key ? { ...it, ...v } : it,
                    );
                    return newInfo;
                  })
                }
                minHeight={26}
                showDragBtn
                onDelete={delKey =>
                  setMessageList(prev => {
                    const newInfo = prev?.filter(it => it.key !== delKey);
                    return newInfo;
                  })
                }
              />
            ))}
        </div>

        <Button
          color="primary"
          icon={<IconCozPlus />}
          onClick={handleAddMessage}
          disabled={readonly}
        >
          {I18n.t('add_message')}
        </Button>
      </div>
    </CollapseCard>
  );
}
