// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useBaseURL } from '@cozeloop/biz-hooks-adapter';
import { I18n } from '@cozeloop/i18n-adapter';
import {
  type CommonFieldProps,
  type SelectProps,
  withField,
} from '@coze-arch/coze-design';

import { type CreateExperimentValues } from '../../../types/evaluate-target';
import { PromptEvalTargetVersionSelect } from '../../../components/selectors/evaluate-target';
import { OpenDetailText } from '../../../components/common';

const FormSelectInner = withField(PromptEvalTargetVersionSelect);

interface EvalTargetVersionProps {
  promptId: string;
  sourceTargetVersion: string;
  onChange: (key: keyof CreateExperimentValues, value: unknown) => void;
}

const PromptEvalTargetVersionFormSelect: React.FC<
  SelectProps & CommonFieldProps & EvalTargetVersionProps
> = props => {
  const { promptId, sourceTargetVersion } = props;
  const { baseURL } = useBaseURL();

  return (
    <FormSelectInner
      remote
      onChangeWithObject
      rules={[
        {
          required: true,
          message: I18n.t('please_select', { field: I18n.t('version') }),
        },
      ]}
      label={{
        text: I18n.t('version'),
        className: 'justify-between pr-0',
        extra: (
          <>
            {promptId && sourceTargetVersion ? (
              <OpenDetailText
                url={`${baseURL}/pe/prompts/${
                  promptId
                }?version=${sourceTargetVersion}`}
              />
            ) : null}
          </>
        ),
      }}
      {...props}
    />
  );
};

export default PromptEvalTargetVersionFormSelect;
