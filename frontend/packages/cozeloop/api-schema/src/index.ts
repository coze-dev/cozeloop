// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @typescript-eslint/naming-convention -- skip */
export { ApiOption } from './api/config';
export * from './api/idl';

export { $notification } from './notification';

import {
  evaluationEvalSet,
  evaluationEvalTarget,
  evaluationEvaluator,
  evaluationExpt,
  dataDataset,
  llmManage,
  foundationUpload,
  promptManage,
  promptDebug,
  foundationAuthn,
  foundationUser,
  foundationSpace,
} from './api/idl';

export const StoneEvaluationApi = {
  ...evaluationEvalSet,
  ...evaluationEvalTarget,
  ...evaluationEvaluator,
  ...evaluationExpt,
};

export const DataApi = {
  ...dataDataset,
};

export const LlmManageApi = {
  ...llmManage,
};
export const FoundationApi = {
  ...foundationUpload,
  ...foundationAuthn,
  ...foundationUser,
  ...foundationSpace,
};

export const StonePromptApi = {
  ...promptManage,
  ...promptDebug,
};
