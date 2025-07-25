// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useShallow } from 'zustand/react/shallow';
import { nanoid } from 'nanoid';
import { type PromptMessage } from '@cozeloop/prompt-components';
import { I18n } from '@cozeloop/i18n-adapter';
import { CollapseCard } from '@cozeloop/components';
import { Role, VariableType } from '@cozeloop/api-schema/prompt';
import { Typography } from '@coze-arch/coze-design';

import { useBasicStore } from '@/store/use-basic-store';
import { useCompare } from '@/hooks/use-compare';

import { VariableInput } from '../variable-input';

interface VariablesCardProps {
  uid?: number;
  defaultVisible?: boolean;
}

export function VariablesCard({ uid, defaultVisible }: VariablesCardProps) {
  const {
    streaming,
    variables,
    setMessageList,
    mockVariables,
    setMockVariables,
  } = useCompare(uid);
  const { readonly: basicReadonly } = useBasicStore(
    useShallow(state => ({
      readonly: state.readonly,
    })),
  );

  const readonly = basicReadonly || streaming;

  const onDeleteVariable = (key?: string) => {
    if (key) {
      const variableItem = variables?.find(it => it.key === key);
      if (variableItem?.type === VariableType.Placeholder) {
        setMessageList(list => {
          if (!Array.isArray(list)) {
            return [];
          }
          const newList = list?.filter(
            it => !(it.role === Role.Placeholder && it.content === key),
          );
          return newList;
        });
      } else {
        const rep = new RegExp(`{{${key}}}`, 'g');
        setMessageList(list => {
          if (!Array.isArray(list)) {
            return [];
          }
          const newList = list?.map(it => {
            if (it.content) {
              return {
                ...it,
                key: rep.test(it.content) ? nanoid() : it.key,
                content: it.content.replace(rep, ''),
              };
            }
            return it;
          });
          return newList;
        });
      }
    }
  };

  const changeInputVariableValue = ({
    key,
    value,
    messageList,
  }: {
    key?: string;
    value?: string;
    messageList?: PromptMessage[];
  }) => {
    setMockVariables(list => {
      if (!Array.isArray(list)) {
        return [];
      }
      const newList = list?.map(it => {
        if (it.key === key) {
          return {
            ...it,
            value,
            placeholder_messages: messageList,
          };
        }
        return it;
      });
      return newList;
    });
  };

  return (
    <CollapseCard
      title={
        <Typography.Text strong>{I18n.t('prompt_variable')}</Typography.Text>
      }
      defaultVisible={defaultVisible}
    >
      <div className="flex flex-col gap-2 pt-4 pb-3">
        {mockVariables?.map(item => {
          const variable = variables?.find(it => it.key === item.key);
          return (
            <VariableInput
              key={`${item.key}`}
              variableKey={item.key}
              variableType={variable?.type}
              placeholderMessages={item.placeholder_messages}
              readonly={readonly}
              variableValue={item.value}
              onDelete={key => onDeleteVariable(key)}
              onValueChange={value => changeInputVariableValue({ ...value })}
            />
          );
        })}

        {variables?.some(it => it.key) ? null : (
          <Typography.Text
            type="tertiary"
            style={{ color: 'var(--coz-fg-dim)' }}
          >
            {I18n.t('no_variable')}
          </Typography.Text>
        )}
      </div>
    </CollapseCard>
  );
}
