// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useMemo } from 'react';

import { I18n } from '@cozeloop/i18n-adapter';
import { OutputInfo, useGlobalEvalConfig } from '@cozeloop/evaluate-components';
import { EvaluatorType } from '@cozeloop/api-schema/evaluation';
import { FormSelect, withField } from '@coze-arch/coze-design';

import { PromptField } from './prompt-field';

export function ConfigContent({
  refreshEditorModelKey,
  disabled,
}: {
  refreshEditorModelKey?: number;
  disabled?: boolean;
}) {
  const { modelConfigEditor } = useGlobalEvalConfig();

  const FormModelConfig = useMemo(
    () => withField(modelConfigEditor),
    [modelConfigEditor],
  );

  return (
    <>
      <FormSelect
        label={I18n.t('evaluator_type')}
        field="evaluator_type"
        initValue={EvaluatorType.Prompt}
        fieldClassName="hidden"
      />
      <FormModelConfig
        refreshModelKey={refreshEditorModelKey}
        label={I18n.t('model_selection')}
        disabled={disabled}
        field="current_version.evaluator_content.prompt_evaluator.model_config"
        rules={[{ required: true, message: I18n.t('choose_model') }]}
      />
      <PromptField
        disabled={disabled}
        refreshEditorKey={refreshEditorModelKey}
      />

      <OutputInfo />
    </>
  );
}
