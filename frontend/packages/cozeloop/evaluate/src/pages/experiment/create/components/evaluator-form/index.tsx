// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useState } from 'react';

import { isEqual } from 'lodash-es';
import { useDebounceFn } from 'ahooks';
import { I18n } from '@cozeloop/i18n-adapter';
import {
  COMMON_OUTPUT_FIELD_NAME,
  DEFAULT_TEXT_STRING_SCHEMA,
} from '@cozeloop/evaluate-components';
import { type FieldSchema } from '@cozeloop/api-schema/evaluation';
import { useFormState, ArrayField } from '@coze-arch/coze-design';

import {
  type EvaluatorPro,
  type CreateExperimentValues,
} from '@/types/experiment/experiment-create';

import { EvaluatorFieldItem } from '../../evaluator-field-item';

export interface EvaluatorFormProps {
  initValue: CreateExperimentValues['evaluatorProList'];
  evaluationSetVersionDetail: CreateExperimentValues['evaluationSetVersionDetail'];
}

// 评测对象的输出字段定义
const evaluateTargetSchemas: FieldSchema[] = [
  {
    name: COMMON_OUTPUT_FIELD_NAME,
    description: I18n.t('actual_output'),
    ...DEFAULT_TEXT_STRING_SCHEMA,
  },
];

// 获取已选择的评估器版本ID列表
const getSelectedVersionIds = (evaluatorProList: EvaluatorPro[]) => {
  const list: string[] = [];
  evaluatorProList?.forEach(ep => {
    const versionId = ep?.evaluatorVersion?.id;
    if (versionId) {
      list.push(String(versionId));
    }
  });
  return list;
};

export const EvaluatorForm = (props: EvaluatorFormProps) => {
  const { initValue, evaluationSetVersionDetail } = props;
  const formState = useFormState();
  const formValues = formState.values as CreateExperimentValues;
  const evaluatorProList = formValues?.evaluatorProList || [];

  const [selectedVersionIds, setSelectedVersionIds] = useState(() =>
    getSelectedVersionIds(evaluatorProList),
  );

  // 计算并更新已选择的评估器版本ID列表
  const calcSelectedVersionIds = useDebounceFn(
    (ls: EvaluatorPro[]) => {
      const newList = getSelectedVersionIds(ls);

      setSelectedVersionIds(pre => {
        if (isEqual(pre, newList)) {
          return pre;
        }
        return newList;
      });
    },
    {
      wait: 200,
    },
  );
  // TODO: FIXME: @武文琦 这里的evaluatorProList副作用更新不及时 更新选中的版本ID
  calcSelectedVersionIds.run(evaluatorProList);
  // useEffect(() => {
  //   // 更新选中的版本ID
  //   calcSelectedVersionIds.run(evaluatorProList);
  // }, [evaluatorProList]);

  return (
    <div className="flex flex-col gap-5 mt-3">
      <ArrayField field="evaluatorProList" initValue={initValue || [{}]}>
        {({ add, arrayFields }) =>
          arrayFields?.map((arrayField, index) => (
            <EvaluatorFieldItem
              key={arrayField.key + arrayField.field}
              arrayField={arrayField}
              index={index}
              disableDelete={arrayFields.length <= 1}
              evaluationSetSchemas={
                evaluationSetVersionDetail?.evaluation_set_schema?.field_schemas
              }
              evaluateTargetSchemas={evaluateTargetSchemas}
              selectedVersionIds={selectedVersionIds}
              disableAdd={arrayFields.length >= 10}
              isLast={index === arrayFields.length - 1}
              onAdd={add}
            />
          ))
        }
      </ArrayField>
    </div>
  );
};
