// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useState } from 'react';

import { safeJsonParse } from './utils/json';
import { type SpanDefinition, type Span } from './types';
import {
  BUILT_IN_SPAN_DEFINITIONS,
  defaultSpanDefinition,
} from './span-definition';
import { structDataListKeys } from './consts';
import { TypeEnum } from './components/span-content-container';
import { RawContent, SpanContentContainer } from './components';

export interface TraceDetailLayoutProps {
  span: Span;
  customSpanDefinition?: SpanDefinition[];
}

function getRender(key: string, spanDefinition: SpanDefinition) {
  if (key === 'input') {
    return spanDefinition.renderInput;
  }
  if (key === 'output') {
    return spanDefinition.renderOutput;
  }
  if (key === 'tool') {
    return spanDefinition.renderTool;
  }
  if (key === 'reasoningContent') {
    return spanDefinition.renderReasoningContent;
  }
  if (key === 'error') {
    return spanDefinition.renderError;
  }
  return () => <></>;
}
const spanDefinitionMap = new Map<string, SpanDefinition>();

export const TraceStructData = (props: TraceDetailLayoutProps) => {
  const { span, customSpanDefinition } = props;
  const [spanDefinitionInitialized, setSpanDefinitionInitialized] =
    useState(false);

  useEffect(() => {
    (
      [
        ...BUILT_IN_SPAN_DEFINITIONS,
        ...(customSpanDefinition ?? []),
      ] as SpanDefinition[]
    ).forEach(spanDefinition => {
      if (spanDefinitionMap.has(spanDefinition.name)) {
        console.warn(
          `spanDefinition ${spanDefinition.name} already exists, it will be overwritten`,
        );
      }
      spanDefinitionMap.set(spanDefinition.name, spanDefinition);
    });
    setSpanDefinitionInitialized(true);
  }, [customSpanDefinition]);

  const { span_type } = span;
  if (!span_type || !spanDefinitionInitialized) {
    return null;
  }

  const targetSpanDefinition =
    spanDefinitionMap.get(span_type) ??
    (defaultSpanDefinition as SpanDefinition);

  const parseResults = targetSpanDefinition.parseSpanContent(span);

  return (
    <div>
      {structDataListKeys.map(key => {
        const result = parseResults[key];

        if (result.isEmpty) {
          return null;
        }
        return (
          <SpanContentContainer
            key={key}
            content={result.originalContent ?? ''}
            title={key}
            copyConfig={{
              moduleName: 'trace',
              point: 'span',
            }}
            children={(renderType, content) => {
              if (!result.isValidate || renderType === TypeEnum.JSON) {
                return (
                  <RawContent
                    structuredContent={content}
                    attrTos={span.attr_tos}
                    tagType={result.tagType}
                  />
                );
              }

              const structuredContent =
                typeof result.content === 'string'
                  ? safeJsonParse(result.content)
                  : result.content;
              return getRender(key, targetSpanDefinition)(
                span,
                structuredContent as any,
              );
            }}
          />
        );
      })}
    </div>
  );
};

export { SpanContentContainer, RawContent } from './components';
export { getSpanContentField } from './utils/span';
