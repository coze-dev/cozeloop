// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import {
  ContentType,
  EvalTargetType,
  type FieldSchema,
} from '@cozeloop/api-schema/evaluation';

export const evalTargetTypeMap = {
  [EvalTargetType.CozeBot]: I18n.t('coze_agent'),
  [EvalTargetType.CozeLoopPrompt]: 'Prompt',
};

export const evalTargetTypeOptions = [
  {
    label: evalTargetTypeMap[EvalTargetType.CozeBot],
    value: EvalTargetType.CozeBot,
  },
  {
    label: evalTargetTypeMap[EvalTargetType.CozeLoopPrompt],
    value: EvalTargetType.CozeLoopPrompt,
  },
];

export const COZE_BOT_INPUT_FIELD_NAME = 'input';
export const COMMON_OUTPUT_FIELD_NAME = 'actual_output';

export const DEFAULT_TEXT_STRING_SCHEMA: FieldSchema = {
  content_type: ContentType.Text,
  text_schema: '{"type": "string"}',
};
