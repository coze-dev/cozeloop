// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import type { FieldMeta } from '@cozeloop/api-schema/observation';

import { type LogicValue } from '@/components/logic-expr/logic-expr';

import type { PersistentFilter } from '../typings/index';
import { useTraceContext } from '../contexts/trace-context';
import type { PresetRange } from '../consts/time';

export interface TraceStoreState {
  timestamps: [number, number];
  fieldMetas?: Record<string, FieldMeta | undefined>;
  presetTimeRange: PresetRange;
  refreshFlag: number;
  selectedSpanType: string | number;
  filters?: LogicValue;
  persistentFilters: PersistentFilter[];
  relation: string;
  selectedPlatform: string | number;
  applyFilters?: LogicValue;
  filterPopupVisible: boolean;
  lastUserRecord: {
    filters?: LogicValue;
    selectedPlatform?: string | number;
    selectedSpanType?: string | number;
  };
}

export interface TraceStoreAction {
  setTimestamps: (e: [number, number]) => void;
  setFilters: (val?: LogicValue) => void;
  setFieldMetas: (e?: Record<string, FieldMeta>) => void;
  setPresetTimeRange: (e: PresetRange) => void;
  updateRefreshFlag: () => void;
  setSelectedSpanType: (e: string | number) => void;
  setSelectedPlatform: (e: string | number) => void;
  setApplyFilters: (e: LogicValue) => void;
  setFilterPopupVisible: (e: boolean) => void;
  setLastUserRecord: (e: TraceStoreState['lastUserRecord']) => void;
}

export type TraceStoreType = TraceStoreState & TraceStoreAction;

export function useTraceStore(
  selector?: (state: TraceStoreType) => TraceStoreType,
) {
  const context = useTraceContext();
  if (!selector) {
    return context as unknown as TraceStoreType;
  }
  return selector(context as unknown as TraceStoreType);
}

export function useStoreState<T>(selector: (state: TraceStoreType) => T): T {
  const context = useTraceContext();
  return selector(context as unknown as TraceStoreType);
}
