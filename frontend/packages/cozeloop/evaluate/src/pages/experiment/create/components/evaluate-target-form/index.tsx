// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import { useEvalTargetDefinition } from '@cozeloop/evaluate-components';
import { type EvalTargetType } from '@cozeloop/api-schema/evaluation';
import { type Form, FormSelect, useFormState } from '@coze-arch/coze-design';

import { type CreateExperimentValues } from '@/types/experiment/experiment-create';

import { evaluateTargetValidators } from '../validators/evaluate-target';

export interface EvaluateTargetFormProps {
  formRef: React.RefObject<Form<CreateExperimentValues>>;
  createExperimentValues: CreateExperimentValues;
  setCreateExperimentValues: React.Dispatch<
    React.SetStateAction<CreateExperimentValues>
  >;
}

export const EvaluateTargetForm = (props: EvaluateTargetFormProps) => {
  const { formRef, createExperimentValues, setCreateExperimentValues } = props;
  const formState = useFormState();

  const { values: formValues } = formState;

  const { getEvalTargetDefinitionList, getEvalTargetDefinition } =
    useEvalTargetDefinition();

  const pluginEvaluatorList = getEvalTargetDefinitionList();

  const evalTargetTypeOptions = pluginEvaluatorList
    .filter(e => e.selector)
    .map(eva => ({
      label: eva.name,
      value: eva.type,
    }));

  const currentEvaluator = getEvalTargetDefinition?.(
    formValues.evalTargetType as string,
  );

  const handleEvalTargetTypeChange = (v: EvalTargetType) => {
    // 评测类型修改, 清空相关字段
    formRef.current?.formApi?.setValues({
      ...formValues,
      evalTargetType: v as EvalTargetType,
      evalTarget: undefined,
      evalTargetVersion: undefined,
      evalTargetMapping: undefined,
    });
  };

  const targetType = formValues.evalTargetType;

  const TargetFormContent = currentEvaluator?.evalTargetFormSlotContent;

  const handleOnFieldChange = (
    key: keyof CreateExperimentValues,
    value: unknown,
  ) => {
    if (key) {
      formRef.current?.formApi?.setValue(key, value);
    }
  };

  return (
    <>
      <FormSelect
        className="w-full"
        field="evalTargetType"
        label={I18n.t('type')}
        optionList={evalTargetTypeOptions}
        placeholder={I18n.t('please_select', { field: I18n.t('type') })}
        rules={evaluateTargetValidators.evalTargetType}
        onChange={v => handleEvalTargetTypeChange(v as EvalTargetType)}
      />

      {targetType ? (
        <>
          {TargetFormContent ? (
            <TargetFormContent
              formValues={formState.values}
              createExperimentValues={createExperimentValues}
              onChange={handleOnFieldChange}
              setCreateExperimentValues={setCreateExperimentValues}
            />
          ) : null}
        </>
      ) : null}
    </>
  );
};
