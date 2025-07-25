// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { getSpanContentField } from '../utils/span';
import { type Span } from '../types';
import { SpanContentContainer } from './span-content-container';
import { RawContent } from './raw-content';

interface SpanContentDetailProps {
  span: Span;
}
export const SpanContentDetail = (props: SpanContentDetailProps) => {
  const { span } = props;
  const spanDetailList = getSpanContentField(span);
  return (
    <>
      {spanDetailList.map((spanDetail, index) => (
        <SpanContentContainer
          key={spanDetail.title}
          content={spanDetail.content}
          title={spanDetail.title}
          hasBottomLine={index !== spanDetailList.length - 1}
          copyConfig={{
            moduleName: 'trace',
            point: 'span',
          }}
          children={(_renderType, content) => (
            <RawContent
              structuredContent={content}
              tagType={spanDetail.tagType}
              attrTos={span.attr_tos}
            />
          )}
        />
      ))}
    </>
  );
};
