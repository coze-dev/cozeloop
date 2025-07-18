// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useEffect, useMemo, useState } from 'react';

import { I18n } from '@cozeloop/i18n-adapter';
import {
  LogicEditor,
  EvaluatorPreview,
  type LogicFilter,
  ColumnsManage,
  RefreshButton,
  uniqueExperimentsEvaluators,
} from '@cozeloop/evaluate-components';
import { type Experiment } from '@cozeloop/api-schema/evaluation';
import { type ColumnProps } from '@coze-arch/coze-design';

import { TableHeader } from '@/components/table-for-experiment';
import TableCellExpand from '@/components/common/table-cell-expand';

import { type ExperimentContrastItem } from '../../utils/tools';
import { getExperimentContrastLogicFields } from '../../utils/logic-filter-tools';

// eslint-disable-next-line @coze-arch/max-line-per-function
export default function ExperimentContrastTableHeader({
  experiments = [],
  columnManageStorageKey,
  columns,
  logicFilter,
  setLogicFilter,
  onFilterClose,
  expand,
  setExpand,
  hiddenFieldMap,
  setHiddenFieldMap,
  setHiddenColumnMap,
  onRefresh,
}: {
  experiments: Experiment[] | undefined;
  columnManageStorageKey: string;
  columns: ColumnProps[];
  logicFilter: LogicFilter | undefined;
  setLogicFilter: (logicFilter?: LogicFilter) => void;
  onFilterClose?: () => void;
  expand?: boolean;
  setExpand?: (expand: boolean) => void;
  hiddenFieldMap: Record<string, boolean>;
  setHiddenFieldMap: (map: Record<string, boolean>) => void;
  setHiddenColumnMap: (map: Record<string, boolean>) => void;
  onRefresh?: () => void;
}) {
  const [defaultColumns, setDefaultColumns] = useState<ColumnProps[]>([]);
  const [manageableColumns, setManageableColumns] = useState<ColumnProps[]>([]);

  const logicFields = useMemo(
    () => getExperimentContrastLogicFields(experiments),
    [experiments],
  );

  const handleColumnManageChange: React.Dispatch<
    React.SetStateAction<ColumnProps[]>
  > = newManageColumnsData => {
    setManageableColumns(oldManageableColumns => {
      // 先获得新的管理列数组
      let newManageColumns: ColumnProps[] = [];
      if (typeof newManageColumnsData === 'function') {
        newManageColumns = [...newManageColumnsData(oldManageableColumns)];
      } else {
        newManageColumns = newManageColumnsData;
      }
      // 处理列变化
      const newHiddenColumnMap: Record<string, boolean> = {};
      newManageColumns.forEach(column => {
        if (column.canManage) {
          newHiddenColumnMap[column.key ?? ''] = column.hidden;
        }
      });
      setHiddenColumnMap(newHiddenColumnMap);

      // 处理评估器指标变化
      const newHiddenEvalutorMap: Record<Int64, boolean> = {};
      newManageColumns.forEach(column => {
        if (column.isFieldColumn) {
          newHiddenEvalutorMap[column.key ?? ''] = column.hidden;
        }
      });
      setHiddenFieldMap(newHiddenEvalutorMap);
      // 存储列管理数据到本地
      localStorage.setItem(
        columnManageStorageKey,
        JSON.stringify({
          hiddenFieldMap: newHiddenEvalutorMap,
          hiddenColumnMap: newHiddenColumnMap,
        }),
      );
      return newManageColumns;
    });
  };

  useEffect(() => {
    const newColumns = columns.filter(column => column.canManage);
    const evaluators = uniqueExperimentsEvaluators(experiments);
    const evaluatorColumns: ColumnProps<ExperimentContrastItem>[] =
      evaluators.map(evaluator => ({
        title: <EvaluatorPreview evaluator={evaluator} />,
        displayName: `${evaluator?.name} v${evaluator?.current_version?.version}`,
        // title: `${evaluator?.name} v${evaluator?.current_version?.version}`,
        dataIndex: `${evaluator?.current_version?.id ?? ''}`,
        key: `${evaluator?.current_version?.id ?? ''}`,
        // 标记为评估器列
        isFieldColumn: true,
        hidden: hiddenFieldMap[evaluator?.current_version?.id ?? ''] ?? false,
      }));

    newColumns.push(...evaluatorColumns);

    newColumns.push(
      ...[
        // {
        //   title: 'Latency',
        //   dataIndex: 'latency',
        //   key: 'latency',
        //   isFieldColumn: true,
        //   hidden: hiddenFieldMap.latency ?? false,
        // },
        // {
        //   title: 'Token',
        //   dataIndex: 'token',
        //   key: 'token',
        //   isFieldColumn: true,
        //   hidden: hiddenFieldMap.token ?? false,
        // },
        {
          title: I18n.t('status'),
          dataIndex: 'status',
          key: 'status',
          isFieldColumn: true,
          hidden: hiddenFieldMap.status ?? false,
        },
      ],
    );
    setManageableColumns(newColumns);
    setDefaultColumns(
      newColumns.map(column => ({
        ...column,
        hidden: false,
      })),
    );
  }, [columns, experiments, hiddenFieldMap]);

  const filters = (
    <>
      <LogicEditor
        fields={logicFields}
        value={logicFilter}
        onConfirm={newVal => setLogicFilter(newVal)}
        onClose={onFilterClose}
      />
    </>
  );

  const actions = (
    <>
      <TableCellExpand
        className="ml-auto"
        expand={expand}
        onChange={setExpand}
      />
      <ColumnsManage
        columns={manageableColumns}
        defaultColumns={defaultColumns}
        onColumnsChange={handleColumnManageChange}
        sortable={false}
      />
      <RefreshButton onRefresh={onRefresh} />
    </>
  );
  return <TableHeader filters={filters} actions={actions} />;
}
