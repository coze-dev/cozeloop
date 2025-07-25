// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type Span } from '@/trace-detail/typings/params';

import { type TreeProps } from '../tree/typing';

export type TraceTreeProps = {
  dataSource: SpanNode;
  selectedSpanId?: string;
  onCollapseChange: (id: string) => void;
} & Pick<
  TreeProps,
  | 'indentDisabled'
  | 'lineStyle'
  | 'globalStyle'
  | 'onSelect'
  | 'onClick'
  | 'onMouseMove'
  | 'onMouseEnter'
  | 'onMouseLeave'
  | 'className'
>;

export type SpanNode = Span & {
  children?: SpanNode[];
  isCollapsed: boolean;
  isLeaf: boolean;
};
