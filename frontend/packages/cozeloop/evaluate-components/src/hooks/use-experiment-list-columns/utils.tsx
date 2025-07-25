// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import { UserProfile } from '@cozeloop/components';
import {
  ExptRetryMode,
  type UserInfo,
  type Experiment,
} from '@cozeloop/api-schema/evaluation';
import { StoneEvaluationApi } from '@cozeloop/api-schema';
import { Modal, type ColumnProps } from '@coze-arch/coze-design';

import { formateTime } from '../../utils';
import { TypographyText } from '../../components/text-ellipsis';
import { EvaluateTargetTypePreview } from '../../components/previews/evaluate-target-type-preview';
import { EvalTargetPreview } from '../../components/previews/eval-target-preview';
import { EvaluationSetPreview } from '../../components/previews/eval-set-preview';
import { ExperimentRunStatus } from '../../components/experiments/previews/experiment-run-status';
import LoopTableSortIcon from '../../components/dataset-list/sort-icon';
import ExperimentEvaluatorAggregatorScore from './experiment-evaluator-aggregator-score';

/** 实验列表列配置 */

export function getExperimentColumns({
  spaceID,
  enableSort = false,
}: {
  spaceID: Int64;
  enableSort?: boolean;
  onRefresh?: () => void;
}) {
  const columns: ColumnProps<Experiment>[] = [
    {
      title: I18n.t('experiment_name'),
      disableColumnManage: true,
      dataIndex: 'name',
      key: 'name',
      width: 200,
      render: text => <TypographyText>{text}</TypographyText>,
    },
    {
      title: I18n.t('evaluation_object_type'),
      dataIndex: 'type',
      key: 'type',
      width: 120,
      render(_, record) {
        return (
          <EvaluateTargetTypePreview
            type={record.eval_target?.eval_target_type}
          />
        );
      },
    },
    {
      title: I18n.t('evaluation_object'),
      dataIndex: 'eval_target',
      key: 'eval_target',
      width: 215,
      render(val) {
        return (
          <EvalTargetPreview
            spaceID={spaceID}
            evalTarget={val}
            enableLinkJump={true}
            showIcon={true}
            jumpBtnClassName={'show-in-table-row-hover'}
          />
        );
      },
    },
    {
      title: I18n.t('associated_evaluation_set'),
      dataIndex: 'eval_set',
      key: 'eval_set',
      width: 215,
      render: val => (
        <EvaluationSetPreview
          evalSet={val}
          enableLinkJump={true}
          jumpBtnClassName={'show-in-table-row-hover'}
        />
      ),
    },
    {
      title: I18n.t('status'),
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (_, record: Experiment) => (
        <div onClick={e => e.stopPropagation()}>
          <ExperimentRunStatus
            status={record.status}
            experiment={record}
            enableOnClick={false}
            showProcess={true}
          />
        </div>
      ),
    },
    {
      title: I18n.t('score'),
      dataIndex: 'score',
      key: 'score',
      width: 330,
      render: (_, record: Experiment) => (
        <ExperimentEvaluatorAggregatorScore
          evaluators={record.evaluators ?? []}
          spaceID={spaceID}
          evaluatorAggregateResult={
            record.expt_stats?.evaluator_aggregate_results
          }
        />
      ),
    },
    {
      title: I18n.t('description'),
      dataIndex: 'desc',
      key: 'desc',
      width: 160,
      render: val => <TypographyText>{val || '-'}</TypographyText>,
    },
    {
      title: I18n.t('creator'),
      dataIndex: 'base_info.created_by',
      key: 'create_by',
      width: 160,
      render: (val: UserInfo) =>
        val?.name ? (
          <UserProfile avatarUrl={val?.avatar_url} name={val?.name} />
        ) : (
          '-'
        ),
    },
    {
      title: I18n.t('create_time'),
      dataIndex: 'start_time',
      key: 'start_time',
      width: 180,
      sorter: enableSort,
      sortIcon: LoopTableSortIcon,
      render: val => formateTime(val),
    },
    {
      title: I18n.t('end_time'),
      dataIndex: 'end_time',
      key: 'end_time',
      width: 180,
      render: val => formateTime(val),
    },
  ];
  return columns;
}

export function handleDelete({
  record,
  spaceID,
  onRefresh,
}: {
  record: Experiment;
  spaceID: Int64;
  onRefresh?: () => void;
}) {
  Modal.confirm({
    title: I18n.t('delete_experiment'),
    content: I18n.t('confirm_to_delete_x', {
      name: <span className="font-medium px-[2px]">{record.name}</span>,
    }),
    okText: I18n.t('delete'),
    cancelText: I18n.t('Cancel'),
    okButtonColor: 'red',
    width: 420,
    autoLoading: true,
    async onOk() {
      if (record.id) {
        await StoneEvaluationApi.DeleteExperiment({
          workspace_id: spaceID,
          expt_id: record.id,
        });
        onRefresh?.();
      }
    },
  });
}
export function handleRetry({
  record,
  spaceID,
  onRefresh,
}: {
  record: Experiment;
  spaceID: Int64;
  onRefresh?: () => void;
}) {
  Modal.confirm({
    title: I18n.t('retry_experiment'),
    content: I18n.t('only_re_evaluate_failed_part'),
    okText: I18n.t('confirm'),
    cancelText: I18n.t('Cancel'),
    width: 420,
    autoLoading: true,
    async onOk() {
      await StoneEvaluationApi.RetryExperiment({
        workspace_id: spaceID,
        expt_id: record.id ?? '',
        retry_mode: ExptRetryMode.RetryAll,
      });
      onRefresh?.();
    },
  });
}

export function handleCopy({
  record,
  onOk,
}: {
  record: Experiment;
  onOk: () => void;
}) {
  Modal.confirm({
    title: I18n.t('copy_experiment_config'),
    content: I18n.t('copy_and_run_experiment', {
      name: <span className="font-medium px-[2px]">{record.name}</span>,
    }),
    okText: I18n.t('confirm'),
    cancelText: I18n.t('Cancel'),
    width: 420,
    onOk,
  });
}
