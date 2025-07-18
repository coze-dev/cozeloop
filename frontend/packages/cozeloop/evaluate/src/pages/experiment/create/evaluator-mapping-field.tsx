// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { type FC, useMemo } from 'react';

import { I18n } from '@cozeloop/i18n-adapter';
import { getTypeText } from '@cozeloop/evaluate-components';
import { type FieldSchema } from '@cozeloop/api-schema/evaluation';
import { IconCozEmpty } from '@coze-arch/coze-design/icons';
import {
  type CommonFieldProps,
  EmptyState,
  Loading,
  withField,
} from '@coze-arch/coze-design';

import {
  type OptionSchema,
  type OptionGroup,
} from '@/components/mapping-item-field/types';
import { MappingItemField } from '@/components/mapping-item-field';

import emptyStyles from './empty-state.module.less';

export interface EvaluatorMappingProps {
  loading?: boolean;
  keySchemas?: FieldSchema[];
  evaluationSetSchemas?: FieldSchema[];
  evaluateTargetSchemas?: FieldSchema[];
  prefixField: string;
  value?: Record<string, OptionSchema>;
  onChange?: (v?: Record<string, OptionSchema>) => void;
}

export const EvaluatorMappingField: FC<
  CommonFieldProps & EvaluatorMappingProps
> = withField(function ({
  loading,
  keySchemas,
  evaluationSetSchemas,
  evaluateTargetSchemas,
  prefixField,
  // value,
  // onChange,
  // ...props
}: EvaluatorMappingProps) {
  const optionGroups = useMemo(() => {
    const res: OptionGroup[] = [];
    if (evaluationSetSchemas) {
      res.push({
        schemaSourceType: 'set',
        children: evaluationSetSchemas?.map(s => ({
          ...s,
          schemaSourceType: 'set',
        })),
      });
    }
    if (evaluateTargetSchemas) {
      res.push({
        schemaSourceType: 'target',
        children: evaluateTargetSchemas?.map(s => ({
          ...s,
          schemaSourceType: 'target',
        })),
      });
    }
    return res;
  }, [evaluationSetSchemas]);

  if (loading) {
    return (
      <div className="h-[84px] w-full flex items-center justify-center">
        <Loading
          className="!w-full"
          size="large"
          label={I18n.t('loading_field_mapping')}
          loading={true}
        ></Loading>
      </div>
    );
  }

  if (!keySchemas) {
    return (
      <div className="h-[84px] w-full flex items-center justify-center">
        <EmptyState
          size="default"
          icon={<IconCozEmpty className="coz-fg-dim text-32px" />}
          title={I18n.t('no_data')}
          className={emptyStyles['empty-state']}
          // description="请选择评估器和版本号后再查看"
        />
      </div>
    );
  }

  return (
    <div>
      {keySchemas?.map(k => (
        <MappingItemField
          key={k.name}
          noLabel
          field={`${prefixField}.${k.name}`}
          fieldClassName="!pt-0"
          keyTitle={I18n.t('evaluator')}
          keySchema={k}
          optionGroups={optionGroups}
          rules={[
            {
              validator: (_rule, v) => {
                if (!v) {
                  return new Error(I18n.t('please_select', { field: '' }));
                }
                if (getTypeText(v) !== getTypeText(k)) {
                  return new Error(I18n.t('selected_fields_inconsistent'));
                }
                return true;
              },
            },
          ]}
        />
      ))}
    </div>
  );
});
