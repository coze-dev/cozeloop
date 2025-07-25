// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import {
  type FieldData,
  type Content,
  type EvaluatorRecord,
  type TurnRunState,
} from '@cozeloop/api-schema/evaluation';

export interface ExperimentItem {
  id: Int64;
  groupID: Int64;
  turnID: Int64;
  groupIndex: number;
  turnIndex: number;
  datasetRow: Record<string, FieldData>;
  actualOutput: Content | undefined;
  targetErrorMsg: string | undefined;
  evaluatorsResult: Record<string, EvaluatorRecord | undefined>;
  runState: TurnRunState | undefined;
  itemErrorMsg: string | undefined;
  logID: Int64 | undefined;
  evalTargetTraceID: Int64 | undefined;
}
