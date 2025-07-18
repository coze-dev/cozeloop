// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { arrayToMap } from '@cozeloop/evaluate-components';
import { IDRender, TableColActions } from '@cozeloop/components';
import {
  type ColumnEvaluator,
  type Experiment,
  type EvaluatorRecord,
  type FieldData,
  type ItemResult,
} from '@cozeloop/api-schema/evaluation';
import { Tooltip, type ColumnProps } from '@coze-arch/coze-design';

import { type ExperimentItem } from '@/types/experiment/experiment-detail';
import { ExperimentItemRunStatus } from '@/components/experiment';

import EvaluatorScore from '../evaluator-score';
import { EvaluatorColumnHeader } from './evaluator-column-header';
import { I18n } from '@cozeloop/i18n-adapter';

const experimentDataToRecordItems = (data: ItemResult[]) => {
  const recordItems: ExperimentItem[] = [];
  data.forEach(group => {
    group.turn_results?.forEach(turn => {
      // eslint-disable-next-line complexity
      turn.experiment_results?.forEach(experiment => {
        const { eval_set, evaluator_output, target_output, system_info } =
          experiment.payload ?? {};
        const evaluatorsResult: Record<string, EvaluatorRecord | undefined> =
          {};
        Object.entries(evaluator_output?.evaluator_records ?? {}).forEach(
          ([evaluatorVersionId, record]) => {
            evaluatorsResult[evaluatorVersionId ?? ''] = record;
          },
        );

        const actualOutput =
          target_output?.eval_target_record?.eval_target_output_data
            ?.output_fields?.actual_output;
        const evalTargetTraceID = target_output?.eval_target_record?.trace_id;

        recordItems.push({
          id: `${group.item_id}_${turn.turn_id}`,
          groupID: group.item_id,
          turnID: turn.turn_id ?? '',
          groupIndex: Number(group.item_index) || 0,
          turnIndex: Number(turn.turn_index) || 0,
          datasetRow: arrayToMap<FieldData, FieldData>(
            eval_set?.turn?.field_data_list ?? [],
            'key',
          ),
          actualOutput,
          targetErrorMsg:
            target_output?.eval_target_record?.eval_target_output_data
              ?.eval_target_run_error?.message,
          evaluatorsResult,
          runState: system_info?.turn_run_state,
          itemErrorMsg: system_info?.error?.detail,
          logID: system_info?.log_id,
          evalTargetTraceID,
        });
      });
    });
  });
  return recordItems;
};

const getEvaluatorColumns = (params: {
  columnEvaluators: ColumnEvaluator[];
  spaceID: Int64;
  experiment: Experiment | undefined;
  handleRefresh: () => void;
}) => {
  const { columnEvaluators, spaceID, experiment, handleRefresh } = params;
  const evaluatorColumns: ColumnProps<ExperimentItem>[] = columnEvaluators.map(
    evaluator => ({
      title: (
        <EvaluatorColumnHeader
          evaluator={evaluator}
          tagProps={{ className: 'font-normal' }}
        />
      ),
      // 用来在列管理里面使用的title
      displayName: evaluator.name ?? '-',
      dataIndex: `evaluatorsResult.${evaluator.evaluator_version_id}_${evaluator.name}`,
      key: `${evaluator.evaluator_version_id}_${evaluator.name}`,
      align: 'right',
      width: 180,
      // 本期不支持排序
      // sorter: true,
      // sortIcon: LoopTableSortIcon,
      render(_: unknown, record: ExperimentItem) {
        const evaluatorRecord =
          record.evaluatorsResult?.[evaluator?.evaluator_version_id];
        return (
          <EvaluatorScore
            evaluatorRecord={evaluatorRecord}
            spaceID={spaceID}
            traceID={evaluatorRecord?.trace_id ?? ''}
            evaluatorRecordID={evaluatorRecord?.id ?? ''}
            align="right"
            experiment={experiment}
            onRefresh={handleRefresh}
          />
        );
      },
    }),
  );
  return evaluatorColumns;
};

const getBaseColumn = (
  idColumnTitle?: string,
): ColumnProps<ExperimentItem>[] => [
  {
    title: '',
    // 用来在列管理里面使用的title
    displayName: '状态',
    // 不支持列管理
    disableColumnManage: true,
    dataIndex: 'status',
    key: 'status',
    width: 60,
    render: (_, record: ExperimentItem) => (
      <ExperimentItemRunStatus status={record.runState} onlyIcon={true} />
    ),
  },
  {
    title: idColumnTitle || 'ID',
    disableColumnManage: true,
    dataIndex: 'groupID',
    key: 'id',
    width: 110,
    render: val => <IDRender id={val} useTag={true} />,
  },
];

const getActionColumn = (params: {
  onClick: (record: ExperimentItem) => void;
}): ColumnProps<ExperimentItem> => {
  const { onClick } = params;

  return {
    title: I18n.t('operation'),
    disableColumnManage: true,
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
            onClick: () => onClick(record),
          },
        ]}
      />
    ),
  };
};

export {
  experimentDataToRecordItems,
  getEvaluatorColumns,
  getActionColumn,
  getBaseColumn,
};
