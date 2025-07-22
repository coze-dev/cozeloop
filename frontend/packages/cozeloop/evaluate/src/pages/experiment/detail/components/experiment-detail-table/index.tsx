// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useEffect, useMemo, useState } from 'react';

import { I18n } from '@cozeloop/i18n-adapter';
import {
  ColumnsManage,
  dealColumnsFromStorage,
  LogicEditor,
  getLogicFieldName,
  type LogicField,
  type SemiTableSort,
} from '@cozeloop/evaluate-components';
import {
  type Experiment,
  FieldType,
  type ColumnEvaluator,
  type FieldSchema,
  type TurnRunState,
  type BatchGetExperimentResultResponse,
  ExptStatus,
} from '@cozeloop/api-schema/evaluation';
import { IconCozIllusAdd } from '@coze-arch/coze-design/illustrations';
import { EmptyState, type ColumnProps } from '@coze-arch/coze-design';

import {
  getActualOutputColumn,
  getDatasetColumns,
  isTraceTargetExpr,
} from '@/utils/experiment';
import {
  type Filter,
  type Service,
} from '@/types/experiment/experiment-detail-table';
import { type ExperimentItem } from '@/types/experiment/experiment-detail';
import styles from '@/styles/table-row-hover-show-icon.module.less';
import { useExperimentDetailStore } from '@/hooks/use-experiment-detail-store';
import { useExperimentDetailActiveItem } from '@/hooks/use-experiment-detail-active-item';
import TableForExperiment, {
  TableHeader,
} from '@/components/table-for-experiment';
import {
  ExperimentItemRunStatusSelect,
  EvaluatorColumnPreview,
} from '@/components/experiment';
import TableCellExpand from '@/components/common/table-cell-expand';

import ExperimentItemDetail from '../experiment-item-detail';
import {
  experimentDataToRecordItems,
  getActionColumn,
  getEvaluatorColumns,
  getBaseColumn,
} from './utils';
import { TraceTargetTraceDetailPanel } from './trace-target-trace-detail-panel';

const filterFields: { key: keyof Filter; type: FieldType }[] = [
  {
    key: 'status',
    type: FieldType.TurnRunState,
  },
];

const experimentResultToRecordItems = (
  res: BatchGetExperimentResultResponse,
) => {
  const list = experimentDataToRecordItems(res.item_results ?? []);
  return list;
};

// eslint-disable-next-line @coze-arch/max-line-per-function, complexity
export default function ({
  spaceID,
  experimentID,
  refreshKey,
  experiment,
  onRefreshPage,
}: {
  spaceID: string;
  experimentID: string;
  refreshKey: string;
  experiment: Experiment | undefined;
  onRefreshPage: () => void;
}) {
  const [columns, setColumns] = useState<ColumnProps[]>([]);
  const [defaultColumns, setDefaultColumns] = useState<ColumnProps[]>([]);
  const [fieldSchemas, setFieldSchemas] = useState<FieldSchema[]>([]);
  const [columnEvaluators, setColumnEvaluators] = useState<ColumnEvaluator[]>(
    [],
  );
  const [itemTraceVisible, setTraceVisible] = useState(false);

  const logicFields: LogicField[] = useMemo(() => {
    const fields = columnEvaluators.map(evaluator => {
      const { evaluator_version_id: versionId = '' } = evaluator ?? {};
      const field: LogicField = {
        title: <EvaluatorColumnPreview evaluator={evaluator} />,
        name: getLogicFieldName(FieldType.EvaluatorScore, versionId),
        type: 'number',
        setterProps: { step: 0.1 },
      };
      return field;
    });
    return fields;
  }, [columnEvaluators]);

  const experimentIds = useMemo(() => [experimentID], [experimentID]);

  const {
    service,
    filter,
    setFilter,
    onFilterDebounceChange,
    logicFilter,
    setLogicFilter,
    onLogicFilterChange,
    onSortChange,
    expand,
    setExpand,
  } = useExperimentDetailStore<ExperimentItem, Filter>({
    experimentIds,
    experimentResultToRecordItems,
    pageSizeStorageKey: 'experiment_detail_page_size',
    filterFields,
    refreshKey,
  });

  const activeItemStore = useExperimentDetailActiveItem({
    experimentIds,
    filter,
    logicFilter,
    filterFields,
    experimentResultToRecordItems,
  });
  const {
    activeItem,
    setActiveItem,
    itemDetailVisible,
    setItemDetailVisible,
    onItemStepChange,
  } = activeItemStore;

  const columnManageStorageKey = `experiment_detail_column_manage_${experimentID}`;

  useEffect(() => {
    const res = service.data?.result;
    setFieldSchemas(res?.column_eval_set_fields ?? []);
    setColumnEvaluators(res?.column_evaluators ?? []);
  }, [service.data?.result]);

  const handleRefresh = () => {
    service.refresh();
  };

  useEffect(() => {
    const { column_evaluators = [], column_eval_set_fields = [] } =
      service.data?.result ?? {};
    const isTraceTarget = isTraceTargetExpr(experiment);
    const actualOutputColumns = isTraceTarget
      ? []
      : [
          getActualOutputColumn({
            expand,
            traceIdPath: 'evalTargetTraceID',
            experiment,
          }),
        ];
    // 评估器列
    const evaluatorColumns: ColumnProps<ExperimentItem>[] = getEvaluatorColumns(
      {
        columnEvaluators: column_evaluators,
        spaceID,
        experiment,
        handleRefresh,
      },
    );
    // 操作列
    const actionColumn: ColumnProps<ExperimentItem> = getActionColumn({
      onClick: (record: ExperimentItem) => {
        setActiveItem(record);
        if (isTraceTarget) {
          setTraceVisible(true);
        } else {
          setItemDetailVisible(true);
        }
      },
    });
    // 列配置
    const newColumns: ColumnProps<ExperimentItem>[] = [
      ...getBaseColumn(),
      ...getDatasetColumns(column_eval_set_fields, {
        prefix: 'datasetRow.',
        expand,
      }),
      ...actualOutputColumns,
      ...evaluatorColumns,
    ];
    setColumns([
      ...dealColumnsFromStorage(newColumns, columnManageStorageKey),
      actionColumn,
    ]);
    setDefaultColumns([...newColumns, actionColumn]);
  }, [service.data, spaceID, expand, experiment]);

  const filters = (
    <>
      <ExperimentItemRunStatusSelect
        style={{ minWidth: 170 }}
        value={filter?.status}
        onChange={val => {
          setFilter(oldState => ({
            ...oldState,
            status: val as TurnRunState[],
          }));
          onFilterDebounceChange();
        }}
      />
      <LogicEditor
        fields={logicFields}
        value={logicFilter}
        onConfirm={newVal => {
          setLogicFilter(newVal ?? {});
          onLogicFilterChange(newVal);
        }}
      />
    </>
  );
  // 操作
  const actions = (
    <>
      <TableCellExpand
        className="ml-auto"
        expand={expand}
        onChange={setExpand}
      />
      <ColumnsManage
        columns={columns}
        defaultColumns={defaultColumns}
        storageKey={columnManageStorageKey}
        onColumnsChange={setColumns}
      />
    </>
  );

  // 表格属性
  const tableProps = {
    className: styles['table-row-hover-show-icon'],
    rowKey: 'id',
    columns,
    onRow: record => ({
      onClick: () => {
        // 如果当前有选中的文本，不触发点击事件
        if (!window.getSelection()?.isCollapsed) {
          return;
        }
        setActiveItem(record);
        const isTraceTarget = isTraceTargetExpr(experiment);
        if (isTraceTarget) {
          setTraceVisible(true);
        } else {
          setItemDetailVisible(true);
        }
      },
    }),
    onChange(changeInfo) {
      if (changeInfo.extra?.changeType === 'sorter' && changeInfo.sorter?.key) {
        onSortChange(changeInfo.sorter as unknown as SemiTableSort);
      }
    },
  };

  // 表格空状态
  const tableEmpty =
    experiment?.status === ExptStatus.Pending ? (
      <EmptyState
        size="full_screen"
        icon={<IconCozIllusAdd />}
        title={I18n.t('experiment_initializing')}
        description={I18n.t('wait_and_refresh_page', {
          refresh: (
            <span
              className="text-[rgb(var(--coze-up-brand-9))] cursor-pointer"
              onClick={onRefreshPage}
            >
              {I18n.t('refresh')}
            </span>
          ),
        })}
      />
    ) : (
      <EmptyState
        size="full_screen"
        icon={<IconCozIllusAdd />}
        title={I18n.t('no_data')}
      />
    );

  return (
    <>
      <TableForExperiment<ExperimentItem>
        service={service as Service}
        heightFull={true}
        header={<TableHeader actions={actions} filters={filters} />}
        pageSizeStorageKey="experiment_detail_page_size"
        empty={tableEmpty}
        tableProps={tableProps}
      />
      {activeItem && itemDetailVisible ? (
        <ExperimentItemDetail
          spaceID={spaceID}
          activeItemStore={activeItemStore}
          fieldSchemas={fieldSchemas}
          columnEvaluators={columnEvaluators}
          onClose={() => setItemDetailVisible(false)}
          onStepChange={onItemStepChange}
        />
      ) : null}

      {itemTraceVisible ? (
        <TraceTargetTraceDetailPanel
          // 和服务端的约定，基于Trace的在线评测情况下，traceID和spanID的数据在评测集中存放
          traceID={activeItem?.datasetRow?.trace_id?.content?.text ?? ''}
          spanID={activeItem?.datasetRow?.span_id?.content?.text ?? ''}
          experiment={experiment}
          onClose={() => setTraceVisible(false)}
        />
      ) : null}
    </>
  );
}
