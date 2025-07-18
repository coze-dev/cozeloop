// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useMemo } from 'react';

import { I18n } from '@cozeloop/i18n-adapter';
import { type Model } from '@cozeloop/api-schema/llm-manage';
import { Select, type SelectProps } from '@coze-arch/coze-design';

export interface ModelSelectOption {
  label: React.ReactNode;
  value: string | number;
  model: Model;
}

function getOption(model: Model) {
  const option: ModelSelectOption = {
    label: model.name ?? '',
    value: model.model_id ?? '',
    model,
  };
  return option;
}

export function ModelSelectWithObject(
  props: Omit<SelectProps, 'value' | 'onChange'> & {
    value?: Model;
    onChange?: (model: Model | undefined) => void;
    modelList?: Model[];
  },
) {
  const { value, onChange, modelList = [] } = props;

  const optionList = useMemo(() => modelList?.map(getOption), [modelList]);

  const val = useMemo(() => (value ? getOption(value) : undefined), [value]);

  return (
    <Select
      placeholder={I18n.t('please_select', { field: I18n.t('model') })}
      {...props}
      optionList={optionList}
      // 使value为option对象，不能去掉
      onChangeWithObject={true}
      value={val}
      onChange={newVal => {
        const option = newVal as ModelSelectOption;
        onChange?.(option.model);
      }}
    />
  );
}
