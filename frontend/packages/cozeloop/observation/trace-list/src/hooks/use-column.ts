// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useMemo, useState } from 'react';

import { isEmpty, keys } from 'lodash-es';
import {
  CozeLoopStorage,
  type LocalStorageKeys,
  safeJsonParse,
} from '@cozeloop/toolkit';
import { type ColumnItem } from '@cozeloop/components';

import { type ConvertSpan } from '@/typings/span';

import { type SizedColumn } from '../typings/index';

const storage = new CozeLoopStorage({
  field: 'trace-list',
});

interface UseColumnsParams {
  columnsList: string[];
  columnConfig: Record<string, SizedColumn<ConvertSpan>>;
  storageOptions?: {
    enabled: boolean;
    key: LocalStorageKeys;
  };
}

export const useColumns: (params: UseColumnsParams) => {
  selectedColumns: SizedColumn<ConvertSpan>[];
  onColumnsChange: (newColumns: ColumnItem[]) => void;
  cols: SizedColumn<ConvertSpan>[];
  defaultColumns: SizedColumn<ConvertSpan>[];
  setSelectedColumns: (newColumns: SizedColumn<ConvertSpan>[]) => void;
} = (params: UseColumnsParams) => {
  const { columnsList, columnConfig, storageOptions } = params;
  const { enabled = false, key = '' } = storageOptions || {};

  const localValue = useMemo(() => {
    if (!enabled) {
      return [];
    }
    return safeJsonParse(storage.getItem(key as LocalStorageKeys));
  }, [enabled, key]);

  const defaultColumns = useMemo(
    () =>
      columnsList
        .map(item => {
          const column = columnConfig[item as keyof typeof columnConfig];
          return {
            ...column,
            key: column.dataIndex,
            value: column.displayName,
          };
        })
        .filter(Boolean) as SizedColumn<ConvertSpan>[],
    [columnsList, columnConfig],
  );

  const columns = useMemo(() => {
    if (!enabled || isEmpty(keys(localValue))) {
      return defaultColumns;
    }
    return defaultColumns
      .map(item => {
        const cacheItem = localValue[item.key ?? ''];
        return {
          ...item,
          checked: cacheItem?.checked ?? item.checked,
        };
      })
      .sort((a, b) => {
        const aIndex = localValue[a.key ?? '']?.index ?? Infinity;
        const bIndex = localValue[b.key ?? '']?.index ?? Infinity;
        return aIndex - bIndex;
      });
  }, [enabled, localValue, defaultColumns]);

  const [selectedColumns, setSelectedColumns] = useState(
    columns.filter(item => item.checked),
  );
  const [cols, setCols] = useState([...columns]);

  const onColumnsChange = (newColumns: ColumnItem[]) => {
    const newSelectedColumns = newColumns.filter(
      item => item.checked,
    ) as unknown as SizedColumn<ConvertSpan>[];
    setSelectedColumns([...(newSelectedColumns as SizedColumn<ConvertSpan>[])]);
    setCols([...(newColumns as SizedColumn<ConvertSpan>[])]);
    if (enabled) {
      storage.setItem(
        key as LocalStorageKeys,
        JSON.stringify(
          newColumns.reduce(
            (acc, item) => {
              acc[item.key] = {
                checked: item.checked,
                index: newColumns.findIndex(col => col.key === item.key),
              };
              return acc;
            },
            {} satisfies Record<string, { checked: boolean; index: number }>,
          ),
        ),
      );
    }
  };

  return {
    selectedColumns,
    onColumnsChange,
    cols,
    setSelectedColumns,
    defaultColumns,
  };
};
