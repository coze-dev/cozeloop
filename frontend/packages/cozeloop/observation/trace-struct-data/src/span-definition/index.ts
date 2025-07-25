// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { PromptSpanDefinition } from './prompt';
import { ModelSpanDefinition } from './model';
import { DefaultSpanDefinition } from './default';

const modelSpanDefinition = new ModelSpanDefinition();
const promptSpanDefinition = new PromptSpanDefinition();
export const defaultSpanDefinition = new DefaultSpanDefinition();

export const BUILT_IN_SPAN_DEFINITIONS = [
  modelSpanDefinition,
  promptSpanDefinition,
  defaultSpanDefinition,
];
