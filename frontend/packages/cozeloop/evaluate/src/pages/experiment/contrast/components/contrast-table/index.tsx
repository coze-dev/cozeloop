// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useEffect, useState } from 'react';

import { cloneDeep } from 'lodash-es';
import { type PaginationResult } from 'ahooks/lib/usePagination/types';
import { safeJsonParse } from '@cozeloop/toolkit';
import { I18n } from '@cozeloop/i18n-adapter';
import { type LogicFilter } from '@cozeloop/evaluate-components';
import { TableColActions, IDRender } from '@cozeloop/components';
import {
  type BatchGetExperimentResultResponse,
  type Experiment,
  type FieldSchema,
} from '@cozeloop/api-schema/evaluation';
import { Tooltip, type ColumnProps } from '@coze-arch/coze-design';

import { getDatasetColumns } from '@/utils/experiment';
import styles from '@/styles/table-row-hover-show-icon.module.less';
import { useExperimentDetailStore } from '@/hooks/use-experiment-detail-store';
import { useExperimentDetailActiveItem } from '@/hooks/use-experiment-detail-active-item';
import TableForExperiment from '@/components/table-for-experiment';

import ExperimentContrastItemDetail from '../item-detail';
import {
  type ExperimentContrastItem,
  experimentContrastToRecordItems,
} from '../../utils/tools';
import ExperimentContrastTableHeader from './table-header';
import { getExperimentContrastColumns } from './columns';

interface RequestParams {
  current: number;
  pageSize: number;
  logicFilter?: LogicFilter;
}

type Service = PaginationResult<
  { total: number; list: ExperimentContrastItem[] },
  [RequestParams]
>;

function experimentResultToRecordItems(res: BatchGetExperimentResultResponse) {
  const list = experimentContrastToRecordItems(res.item_results ?? []);
  return list;
}

function readLocalColumn(columnManageStorageKey: string) {
  const str = localStorage.getItem(columnManageStorageKey);
  const data:
    | {
        hiddenFieldMap: Record<Int64, boolean>;
        hiddenColumnMap: Record<Int64, boolean>;
      }
    | undefined = safeJsonParse(str) ?? undefined;
  return data;
}

// eslint-disable-next-line @coze-arch/max-line-per-function
export default function ExperimentContrastTable({
  spaceID,
  experiments = [],
  experimentIds,
  onExperimentChange,
}: {
  spaceID: string;
  experiments: Experiment[] | undefined;
  experimentIds: string[] | undefined;
  onExperimentChange?: (experiments: Experiment[]) => void;
}) {
  const [columns, setColumns] = useState<ColumnProps[]>([]);
  const [fieldSchemas, setFieldSchemas] = useState<FieldSchema[]>([]);

  const columnManageStorageKey = `experiment_contrast_detail_column_manage_${experimentIds?.[0]}`;

  const [hiddenColumnMap, setHiddenColumnMap] = useState<
    Record<string, boolean>
  >({});
  const [hiddenExperimentFieldMap, setHiddenExperimentFieldMap] = useState<
    Record<Int64, boolean>
  >({});

  const { service, logicFilter, onLogicFilterChange, expand, setExpand } =
    useExperimentDetailStore<ExperimentContrastItem>({
      experimentIds,
      experimentResultToRecordItems,
      pageSizeStorageKey: 'experiment_contrast_page_size',
    });

  const activeItemStore = useExperimentDetailActiveItem({
    experimentIds,
    logicFilter,
    experimentResultToRecordItems,
  });
  const { activeItem, setActiveItem, onItemStepChange } = activeItemStore;

  useEffect(() => {
    const res = service.data?.result;
    setFieldSchemas(res?.column_eval_set_fields ?? []);
  }, [service.data?.result]);

  useEffect(() => {
    const newExperiments = cloneDeep(experiments);
    newExperiments.forEach(experiment => {
      experiment.evaluators = experiment.evaluators?.filter(
        evaluator =>
          !hiddenExperimentFieldMap[evaluator?.current_version?.id ?? ''],
      );
    });
    const newColumns: ColumnProps<ExperimentContrastItem>[] = [
      {
        title: 'ID',
        disableColumnManage: true,
        dataIndex: 'groupID',
        key: 'id',
        width: 110,
        hidden: hiddenColumnMap.turnID ?? false,
        canManage: true,
        render: val => <IDRender id={val} useTag={true} />,
      },
      ...getDatasetColumns(fieldSchemas, { prefix: 'datasetRow.', expand }).map(
        column => ({
          ...column,
          hidden: hiddenColumnMap[column.key ?? ''] ?? false,
          canManage: true,
        }),
      ),
      ...getExperimentContrastColumns(newExperiments, {
        expand,
        spaceID,
        enableDelete: true,
        onExperimentChange,
        hiddenFieldMap: hiddenExperimentFieldMap,
        onRefresh: service.refresh,
      }).map(column => ({
        ...column,
        canManage: false,
      })),
      {
        title: I18n.t('operation'),
        dataIndex: 'action',
        key: 'action',
        fixed: 'right',
        align: 'left',
        width: 68,
        render: (_, record) => (
          <TableColActions
            actions={[
              {
                label: (
                  <Tooltip content={I18n.t('view_detail')} theme="dark">
                    {I18n.t('detail')}
                  </Tooltip>
                ),
                onClick: () => {
                  setActiveItem(record);
                },
              },
            ]}
          />
        ),
      },
    ];
    setColumns(newColumns);
  }, [
    spaceID,
    experiments,
    expand,
    fieldSchemas,
    hiddenColumnMap,
    hiddenExperimentFieldMap,
  ]);

  useEffect(() => {
    const key = `experiment_contrast_detail_column_manage_${experimentIds?.[0]}`;
    const data = readLocalColumn(key);
    if (data?.hiddenColumnMap) {
      setHiddenColumnMap(data.hiddenColumnMap);
    }
    if (data?.hiddenFieldMap) {
      setHiddenExperimentFieldMap(data.hiddenFieldMap);
    }
  }, [experimentIds?.[0]]);

  const handleRefresh = () => service.refresh();

  return (
    <div className="h-full flex flex-col gap-3 overflow-hidden">
      <TableForExperiment
        service={service as Service}
        heightFull={true}
        header={
          <ExperimentContrastTableHeader
            experiments={experiments}
            columnManageStorageKey={columnManageStorageKey}
            columns={columns}
            logicFilter={logicFilter}
            setLogicFilter={onLogicFilterChange}
            expand={expand}
            setExpand={setExpand}
            hiddenFieldMap={hiddenExperimentFieldMap}
            setHiddenFieldMap={setHiddenExperimentFieldMap}
            setHiddenColumnMap={setHiddenColumnMap}
            onRefresh={handleRefresh}
          />
        }
        pageSizeStorageKey="experiment_contrast_page_size"
        tableProps={{
          className: styles['table-row-hover-show-icon'],
          rowKey: 'id',
          columns,
          bordered: true,
          onRow: record => ({
            onClick: () => {
              // 如果当前有选中的文本，不触发点击事件
              if (!window.getSelection()?.isCollapsed) {
                return;
              }
              setActiveItem(record);
            },
          }),
        }}
      />
      {activeItem ? (
        <ExperimentContrastItemDetail
          experiments={experiments}
          datasetFieldSchemas={fieldSchemas}
          activeItemStore={activeItemStore}
          spaceID={spaceID}
          onStepChange={onItemStepChange}
          onClose={() => setActiveItem(undefined)}
        />
      ) : null}
    </div>
  );
}
