// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useMemo } from 'react';

import { type EvaluatorContent } from '@cozeloop/api-schema/evaluation';

import { extractDoubleBraceFields } from '../../utils/double-brace';
import { PromptVariablesList } from './prompt-variables-list';
import { PromptMessage } from './prompt-message';
import { OutputInfo } from './output-info';

export function TemplateInfo({
  data,
  notTemplate,
}: {
  data?: EvaluatorContent;
  notTemplate?: boolean;
}) {
  const variables = useMemo(() => {
    if (data?.prompt_evaluator?.message_list) {
      const strSet = new Set<string>();
      data.prompt_evaluator.message_list.forEach(message => {
        const str = message?.content?.text;
        if (str) {
          extractDoubleBraceFields(str).forEach(item => strSet.add(item));
        }
      });
      return Array.from(strSet);
    }
  }, [data]);

  return (
    <>
      {notTemplate ? null : (
        <div className="text-[16px] leading-8 font-medium coz-fg-plus mb-5">
          {data?.prompt_evaluator?.prompt_template_name}
        </div>
      )}

      <div className="text-sm font-medium coz-fg-primary mb-2">{'Prompt'}</div>
      {data?.prompt_evaluator?.message_list?.map((m, idx) => (
        <PromptMessage className="mb-2" key={idx} message={m} />
      ))}

      {variables?.length ? (
        <PromptVariablesList className="mb-3" variables={variables} />
      ) : null}
      <div className="h-2" />
      <OutputInfo />
    </>
  );
}
