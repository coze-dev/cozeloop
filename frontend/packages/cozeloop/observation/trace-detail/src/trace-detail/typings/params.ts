// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type GetTraceResponse } from '@cozeloop/api-schema/observation';

export type Span = GetTraceResponse['spans'][0];
export type TraceAdvanceInfo = GetTraceResponse['traces_advance_info'];

export interface DataSource {
  spans: Span[];
  advanceInfo?: GetTraceResponse['traces_advance_info'];
}

export interface UrlParams {
  spaceID: string;
  id: string;
  searchType: string;
  platformType?: string;
  startTime?: number | string;
  endTime?: number | string;
  defaultSpanID?: string;
  logId?: string;
}

export const enum SpanStatus {
  Success = 'success',
  Error = 'error',
  Broken = 'broken',
  Unknown = 'unknown',
}

export const enum SpanType {
  Unknown = 'unknown',
  Model = 'model',
  Prompt = 'prompt',
  Parser = 'parser',
  Embedding = 'embedding',
  Memory = 'memory',
  Plugin = 'plugin',
  Function = 'function',
  Graph = 'graph',
  Remote = 'remote',
  Loader = 'loader',
  Transformer = 'transformer',
  VectorStore = 'vector_store',
  VectorRetriever = 'vector_retriever',
  Agent = 'agent',
  CozeBot = 'bot',
  LLMCall = 'llm_call',
}
