// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { keyBy, uniqBy } from 'lodash-es';
import { SpanStatus, SpanType } from '@cozeloop/api-schema/observation';

import { type Span } from '../typings/params';
import {
  BROKEN_ROOT_SPAN_ID,
  NODE_CONFIG_MAP,
  NORMAL_BROKEN_SPAN_ID,
  type NodeConfig,
} from '../consts/span';
import { type SpanNode } from '../components/graphs/trace-tree/type';
/** 数组转换成链式节点 */
export function spans2SpanNodes(spans: Span[]) {
  if (spans.length === 0) {
    return;
  }

  const roots: SpanNode[] = [];
  const map: Record<string, SpanNode> = {};

  // 排序 + 去重
  const sortedSpans = uniqBy(spans, span => span.span_id).sort((a, b) => {
    const startA = a.started_at ? Number(a.started_at) : Infinity;
    const startB = b.started_at ? Number(b.started_at) : Infinity;
    return startA - startB;
  });

  sortedSpans.forEach(span => {
    const currentSpan: SpanNode = {
      ...span,
      children: [],
      isLeaf: true,
      isCollapsed: false,
    };
    const { span_id } = span;
    if (span_id) {
      map[span_id] = currentSpan;
    }
  });

  sortedSpans.forEach(span => {
    const { span_id, parent_id } = span;
    if (span_id) {
      const spanNode = map[span_id];
      const parentSpanNode = parent_id ? map[parent_id] : undefined;
      if (parent_id === '0' || parentSpanNode === undefined) {
        roots.push(spanNode);
      } else {
        parentSpanNode.children = parentSpanNode.children ?? [];
        parentSpanNode.children.push(spanNode);
        parentSpanNode.isLeaf = false;
      }
    }
  });

  // const trueRoot = roots.find(root => root.parent_id === '0');
  // const brokenNodes = roots.filter(root => root.parent_id !== '0');

  // if (!trueRoot) {
  //   return appendBrokenToBrokenRoot(brokenNodes);
  // } else if (brokenNodes.length > 0) {
  //   return appendBrokenNodesToRoot(trueRoot, brokenNodes);
  // }
  // return trueRoot;
  return roots;
}

/** 把没有父节点的非根节点挂在到虚拟根节点上 */
export function appendBrokenToBrokenRoot(brokenNodes: SpanNode[]) {
  const vRoot: SpanNode = {
    span_id: BROKEN_ROOT_SPAN_ID,
    parent_id: '0',
    trace_id: '',
    span_name: '',
    type: SpanType.Unknown,
    status: SpanStatus.Success,
    span_type: '',
    status_code: 0,
    started_at: '',
    duration: '',
    input: '',
    output: '',
    custom_tags: {
      device_id: '',
      space_id: '',
      psm_env: '',
      err_msg: '',
      user_id: '',
      psm: '',
    },
    isCollapsed: false,
    isLeaf: false,
    children: brokenNodes,
  };
  return vRoot;
}

/** 把没有父节点的非根节点挂载root上 */
export function appendBrokenNodesToRoot(
  rootNode: SpanNode,
  brokenNodes: SpanNode[],
) {
  const brokenRoot: SpanNode = {
    span_id: NORMAL_BROKEN_SPAN_ID,
    parent_id: rootNode.span_id,
    trace_id: '',
    span_name: '',
    type: SpanType.Unknown,
    status: SpanStatus.Success,
    started_at: '',
    duration: '',
    span_type: '',
    status_code: 0,
    input: '',
    output: '',
    custom_tags: {
      device_id: '',
      space_id: '',
      psm_env: '',
      err_msg: '',
      user_id: '',
      psm: '',
    },
    isCollapsed: false,
    isLeaf: false,
    children: brokenNodes,
  };

  rootNode.children?.push(brokenRoot);
  return rootNode;
}

interface GetNodeConfigParameters {
  /** span type 枚举映射 */
  spanTypeEnum: SpanNode['type'];
  /** span type 字符串，用户真实上报的字段 */
  spanType: string;
}

export function getNodeConfig(params: GetNodeConfigParameters): NodeConfig {
  const { spanTypeEnum, spanType } = params;

  let nodeConfig = NODE_CONFIG_MAP[spanTypeEnum];
  if (!nodeConfig) {
    nodeConfig = NODE_CONFIG_MAP.unknown;
  }

  return {
    ...nodeConfig,
    typeName: spanType ?? '-',
  };
}

export function changeSpanNodeCollapseStatus(
  spanNodes: SpanNode[],
  targetId: string,
): SpanNode[] {
  return spanNodes.map(node => {
    if (node.span_id === targetId) {
      return {
        ...node,
        isCollapsed: !node.isCollapsed,
      };
    }

    if (node.children) {
      return {
        ...node,
        children: changeSpanNodeCollapseStatus(node.children, targetId),
      };
    } else {
      return node;
    }
  });
}

// eslint-disable-next-line  @typescript-eslint/no-explicit-any
type AnyObject = Record<string, any>;
export interface RowMessage {
  role: string;
  content: string | object;
  tool_calls?: AnyObject[];
  parts?: AnyObject[];
  reasoningContent?: string;
}

export enum SpanContentType {
  Model = 'model',
  Prompt = 'prompt',
}

export const getRootSpan = (spans: Span[]) => {
  if (spans.length === 0) {
    return;
  }

  // 排序 + 去重
  const sortedSpans = uniqBy(spans, span => span.span_id).sort((a, b) => {
    const startA = a.started_at ? Number(a.started_at) : Infinity;
    const startB = b.started_at ? Number(b.started_at) : Infinity;
    return startA - startB;
  });

  const map = keyBy(sortedSpans, 'span_id');

  for (const span of sortedSpans) {
    const { span_id, parent_id } = span;
    if (span_id) {
      const parentSpan = parent_id ? map[parent_id] : undefined;
      if (parent_id === '0' || parentSpan === undefined) {
        return span;
      }
    }
  }
};
