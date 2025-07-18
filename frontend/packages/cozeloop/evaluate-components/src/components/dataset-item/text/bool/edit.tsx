// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import { Radio } from '@coze-arch/coze-design';

import { saftJsonParse } from '../../util';
import { type DatasetItemProps } from '../../type';

export const BoolDatasetItemEdit = ({
  fieldContent,
  onChange,
}: DatasetItemProps) => {
  const value = saftJsonParse(fieldContent?.text);
  const isTrue = value === true;
  const isFalse = value === false;
  const handleChange = (newValue: boolean) => {
    onChange?.({
      ...fieldContent,
      text: JSON.stringify(newValue),
    });
  };
  return (
    <div className="flex items-center gap-6">
      <Radio checked={isTrue} onChange={() => handleChange(true)}>
        {I18n.t('yes')}
      </Radio>
      <Radio checked={isFalse} onChange={() => handleChange(false)}>
        {I18n.t('no')}
      </Radio>
    </div>
  );
};
