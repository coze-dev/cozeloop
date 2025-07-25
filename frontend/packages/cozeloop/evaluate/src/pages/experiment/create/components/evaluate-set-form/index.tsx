// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @coze-arch/max-line-per-function */
/* eslint-disable max-len */
/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, type RefObject } from 'react';

import { useRequest } from 'ahooks';
import {
  EvaluateSetSelect,
  EvaluateSetVersionSelect,
  OpenDetailText,
} from '@cozeloop/evaluate-components';
import { useBaseURL, useSpace } from '@cozeloop/biz-hooks-adapter';
import {
  type EvaluationSet,
  type EvaluationSetVersion,
  type FieldSchema,
} from '@cozeloop/api-schema/evaluation';
import { IconCozLoading } from '@coze-arch/coze-design/icons';
import { Form, useFormState, withField } from '@coze-arch/coze-design';

import { type CreateExperimentValues } from '@/types/experiment/experiment-create';
import { getEvaluationSetVersion } from '@/request/evaluation-set';

import { evaluateSetValidators } from '../validators/evaluate-set';
import { EvaluateSetColList } from '../../evaluate-set-col-list';
import { I18n } from '@cozeloop/i18n-adapter';

export interface EvaluateSetFormProps {
  formRef: RefObject<Form<CreateExperimentValues>>;
  createExperimentValues: CreateExperimentValues;
  setCreateExperimentValues: React.Dispatch<
    React.SetStateAction<CreateExperimentValues>
  >;
  setNextStepLoading: (loading: boolean) => void;
}

const FormEvaluateSetSelect = withField(EvaluateSetSelect);

export const EvaluateSetForm = (props: EvaluateSetFormProps) => {
  const {
    formRef,
    setNextStepLoading,
    setCreateExperimentValues,
    createExperimentValues,
  } = props;
  const { spaceID } = useSpace();
  const { baseURL } = useBaseURL();
  const formState = useFormState();

  const { values: formValues } = formState;

  const formSetVersionId = formValues?.evaluationSetVersion;

  const formSetId = formValues?.evaluationSet;

  const versionDetail = createExperimentValues?.evaluationSetVersionDetail;

  const versionDetailService = useRequest(
    async (params: { evaluation_set_id: string; version_id: string }) => {
      const evaluationSetVersionDetail = await getEvaluationSetVersion({
        workspace_id: spaceID,
        ...params,
      });
      // 新版本的 field_schemas
      const newFieldSchemas =
        evaluationSetVersionDetail.version?.evaluation_set_schema
          ?.field_schemas;

      const mappingData =
        formRef.current?.formApi?.getValue('evalTargetMapping');
      // 挨个清空 evalTargetMapping 中的 key
      if (mappingData) {
        const mappingKeys = Object.entries(mappingData) || [];
        mappingKeys.forEach(([key, value]) => {
          // 如果当前字段, 在评测集中存在, 就替换, 不存在就清空
          const findItem = newFieldSchemas?.find(
            item => item.name === value?.key,
          );
          if (findItem) {
            formRef.current?.formApi?.setValue(
              `evalTargetMapping.${key}` as any,
              findItem,
            );
          } else {
            formRef.current?.formApi?.setValue(
              `evalTargetMapping.${key}` as any,
              undefined,
            );
          }
        });
      }
      setCreateExperimentValues(prev => ({
        ...prev,
        // 用于渲染的数据, 不在表单上面, 与表单数据有隔离
        evaluationSetVersionDetail:
          evaluationSetVersionDetail.version as EvaluationSetVersion,
        evaluationSetDetail:
          evaluationSetVersionDetail.evaluation_set as EvaluationSet,
      }));
    },
    {
      manual: true,
    },
  );

  useEffect(() => {
    if (formSetVersionId && formSetId) {
      setNextStepLoading(true);
      versionDetailService.runAsync({
        version_id: formSetVersionId,
        evaluation_set_id: formSetId,
      });
      setNextStepLoading(false);
    }
  }, [formSetId, formSetVersionId]);

  const renderColumns = (fieldSchemas?: FieldSchema[]) => {
    if (versionDetailService.loading) {
      return (
        <div className="flex flex-row items-center">
          <IconCozLoading className="w-4 h-4 animate-spin coz-fg-secondary" />
          <div className="ml-[6px] text-sm coz-fg-secondary">正在加载</div>
        </div>
      );
    }

    return <EvaluateSetColList fieldSchemas={fieldSchemas} />;
  };

  const handleOnEvaluateSetSelectChange = (v: any) => {
    formRef.current?.formApi?.setValue('evaluationSetVersion', undefined);
  };

  return (
    <>
      <div className="flex flex-row gap-5 relative">
        <div className="flex-1 w-0">
          <FormEvaluateSetSelect
            className="w-full"
            field="evaluationSet"
            label={I18n.t('evaluation_set')}
            placeholder={I18n.t('please_select', {
              field: I18n.t('evaluation_set'),
            })}
            rules={evaluateSetValidators.evaluationSet}
            onChange={handleOnEvaluateSetSelectChange}
            onChangeWithObject={false}
          />
        </div>
        <div className="flex-1 flex flex-row items-end">
          <div className="flex-1 w-0">
            <EvaluateSetVersionSelect
              evaluationSetId={formState?.values?.evaluationSet}
              className="w-full"
              field="evaluationSetVersion"
              label={{
                text: I18n.t('version'),
                className: 'justify-between pr-0',
                extra: (
                  <>
                    {formSetVersionId ? (
                      <OpenDetailText
                        className="absolute top-2.5 right-0"
                        url={`${baseURL}/evaluation/datasets/${formState.values.evaluationSet}?version=${formState.values.evaluationSetVersion}`}
                      />
                    ) : null}
                  </>
                ),
              }}
              placeholder={I18n.t('please_select', {
                field: I18n.t('version_number'),
              })}
              rules={evaluateSetValidators.evaluationSetVersion}
            />
          </div>
        </div>
      </div>
      <Form.Slot label={I18n.t('description')}>
        <div className="text-sm coz-fg-primary font-normal">
          {versionDetail?.description || '-'}
        </div>
      </Form.Slot>
      <Form.Slot label={I18n.t('column_name')}>
        {formSetVersionId && formSetId
          ? renderColumns(versionDetail?.evaluation_set_schema?.field_schemas)
          : null}
      </Form.Slot>
      <Form.Slot label={I18n.t('data_total_count')}>
        <div className="text-sm coz-fg-primary font-normal">
          {versionDetail?.item_count ?? '-'}
        </div>
      </Form.Slot>
    </>
  );
};
