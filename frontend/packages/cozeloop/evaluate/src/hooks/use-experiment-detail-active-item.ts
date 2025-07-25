// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useCallback, useEffect, useState } from 'react';

import {
  filterToFilters,
  type LogicFilter,
  type SemiTableSort,
} from '@cozeloop/evaluate-components';
import { useSpace } from '@cozeloop/biz-hooks-adapter';
import {
  type FieldType,
  type BatchGetExperimentResultResponse,
} from '@cozeloop/api-schema/evaluation';

import { DetailItemStepSwitch } from '@/types';
import { batchGetExperimentResult } from '@/request/experiment';

export interface ExperimentDetailActiveItemStore<
  RecordItem extends { groupIndex: number } = { groupIndex: number },
> {
  activeItem: RecordItem | undefined;
  setActiveItem: (item: RecordItem | undefined) => void;
  loading: boolean;
  setLoading: (loading: boolean) => void;
  itemDetailVisible: boolean;
  setItemDetailVisible: (visible: boolean) => void;
  isFirst: boolean;
  setIsFirst: (first: boolean) => void;
  isLast: boolean;
  setIsLast: (last: boolean) => void;
  onItemStepChange: (step: DetailItemStepSwitch) => void;
}

export function useExperimentDetailActiveItem<
  RecordItem extends { groupIndex: number },
  Filter extends object,
>({
  experimentResultToRecordItems,
  experimentIds,
  filter,
  logicFilter,
  filterFields,
  // sort,
}: {
  experimentIds: string[] | undefined;
  logicFilter?: LogicFilter;
  filter?: Filter;
  sort?: SemiTableSort;
  filterFields?: { key: keyof Filter; type: FieldType }[];
  experimentResultToRecordItems: (
    result: BatchGetExperimentResultResponse,
  ) => RecordItem[];
}): ExperimentDetailActiveItemStore<RecordItem> {
  const { spaceID } = useSpace();
  const [activeItem, setActiveItem] = useState<RecordItem | undefined>();
  const [loading, setLoading] = useState(false);
  const [itemDetailVisible, setItemDetailVisible] = useState(false);
  const [isFirst, setIsFirst] = useState(false);
  const [isLast, setIsLast] = useState(false);
  const [total, setTotal] = useState(0);

  const fetchRecordItemByIndex = useCallback(
    async (pageIndex: number) => {
      const filters = filterToFilters<Filter>({
        filter,
        logicFilter,
        filterFields,
      });
      const res = await batchGetExperimentResult({
        workspace_id: spaceID,
        experiment_ids: experimentIds ?? [],
        baseline_experiment_id: experimentIds?.[0] ?? '',
        filters: { [experimentIds?.[0] ?? '']: { filters } },
        page_number: pageIndex + 1,
        page_size: 1,
      });
      const list = experimentResultToRecordItems(res);
      return {
        item: list[0],
        total: Number(res.total) || 0,
      };
    },
    [spaceID, experimentIds, filter, logicFilter],
  );

  const onItemStepChange = async (step: DetailItemStepSwitch) => {
    if (step === DetailItemStepSwitch.Prev && activeItem?.groupIndex === 0) {
      return;
    }
    if (step === DetailItemStepSwitch.Next && isLast) {
      return;
    }
    setLoading(true);
    try {
      const newPageIndex = (activeItem?.groupIndex ?? 0) + step;
      const { item, total: totalCount } =
        await fetchRecordItemByIndex(newPageIndex);
      setActiveItem(item);
      setTotal(totalCount);
    } catch (e) {
      console.error(e);
    }
    setLoading(false);
  };

  useEffect(() => {
    setIsFirst(activeItem?.groupIndex === 0);
    setIsLast(activeItem?.groupIndex === total - 1);
  }, [activeItem, total]);

  return {
    activeItem,
    setActiveItem,
    loading,
    setLoading,
    itemDetailVisible,
    setItemDetailVisible,
    isFirst,
    setIsFirst,
    isLast,
    setIsLast,
    onItemStepChange,
  };
}
