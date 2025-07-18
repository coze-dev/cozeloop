// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import {
  type Message,
  Role,
  type VariableDef,
  VariableType,
} from '@cozeloop/api-schema/prompt';

export const getPlaceholderErrorContent = (
  message?: Message,
  variables?: VariableDef[],
) => {
  if (message?.role === Role.Placeholder) {
    if (!message?.content) {
      return I18n.t('field_not_empty', {
        field: I18n.t('placeholder_var_name'),
      });
    }
    if (!/^[A-Za-z][A-Za-z0-9_]*$/.test(message?.content)) {
      return I18n.t('placeholder_format');
    }
    const normalVariables = variables?.filter(
      it => it.type !== VariableType.Placeholder,
    );
    const hasSameKey = normalVariables?.find(it => it.key === message?.content);
    if (hasSameKey) {
      return I18n.t('placeholder_name_exists');
    }
  }
  return '';
};
