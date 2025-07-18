// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @coze-arch/max-line-per-function */
import { useEffect, useRef, useState } from 'react';

import { I18n } from '@cozeloop/i18n-adapter';
import { GuardPoint, useGuards, GuardActionType } from '@cozeloop/guard';
import {
  TableColActions,
  IDRender,
  type TableColAction,
} from '@cozeloop/components';
import { useNavigateModule } from '@cozeloop/biz-hooks-adapter';
import { ExptStatus } from '@cozeloop/api-schema/evaluation';
import { type Experiment } from '@cozeloop/api-schema/evaluation';
import { Tooltip, type ColumnProps } from '@coze-arch/coze-design';

import { dealColumnsFromStorage } from '../../components/common';
import {
  getExperimentColumns,
  handleDelete,
  handleRetry,
  handleCopy,
} from './utils';

function isExperimentFail(status: ExptStatus | undefined) {
  return [
    ExptStatus.Failed,
    ExptStatus.SystemTerminated,
    ExptStatus.Terminated,
  ].includes(status as ExptStatus);
}

export interface UseExperimentListColumnsProps {
  spaceID: Int64;
  /** 开启ID列，默认为true */
  enableIdColumn?: boolean;
  /** 开启操作列，默认为true  */
  enableActionColumn?: boolean;
  /** 开启列排序，默认为false */
  enableSort?: boolean;
  /** 表格行操作显示控制，默认为true显示 */
  actionVisibleControl?: {
    copy?: boolean;
    retry?: boolean;
    delete?: boolean;
  };
  /** 详情跳转的来源路径（在实验详情页面点击返回跳转的路径） */
  detailJumpSourcePath?: string;
  columnManageStorageKey?: string;
  onRefresh?: () => void;
  onDetailClick?: (e: Experiment) => void;
}

/** 实验列表列配置 */
export function useExperimentListColumns({
  spaceID,
  enableIdColumn = true,
  enableActionColumn = true,
  enableSort = false,
  columnManageStorageKey,
  detailJumpSourcePath,
  actionVisibleControl,
  onRefresh,
  onDetailClick,
}: UseExperimentListColumnsProps) {
  const guards = useGuards({
    points: [
      GuardPoint['eval.experiments.copy'],
      GuardPoint['eval.experiments.delete'],
      GuardPoint['eval.experiments.retry'],
    ],
  });

  const navigate = useNavigateModule();

  const [columns, setColumns] = useState<ColumnProps[]>([]);
  const [defaultColumns, setDefaultColumns] = useState<ColumnProps[]>([]);

  const copyGuardType = guards.data[GuardPoint['eval.experiments.copy']].type;
  const retryGuardType = guards.data[GuardPoint['eval.experiments.retry']].type;
  const deleteGuardType =
    guards.data[GuardPoint['eval.experiments.delete']].type;

  const guardsRef = useRef(guards);
  guardsRef.current = guards;

  const handleRetryOnCLick = (record: Experiment) => {
    const action = () => {
      handleRetry({ record, spaceID, onRefresh });
    };
    if (retryGuardType === GuardActionType.GUARD) {
      guardsRef.current.data[GuardPoint['eval.experiments.retry']].preprocess(
        action,
      );
    } else {
      action();
    }
  };

  const handleDetailOnClick = (record: Experiment) => {
    onDetailClick?.(record);
    navigate(
      `evaluation/experiments/${record.id}`,
      detailJumpSourcePath
        ? { state: { from: detailJumpSourcePath } }
        : undefined,
    );
  };

  const handleCopyOnClick = (record: Experiment) => {
    const action = () => {
      handleCopy({
        record,
        onOk: () => {
          navigate(
            `evaluation/experiments/create?copy_experiment_id=${record.id}`,
          );
        },
      });
    };

    if (copyGuardType === GuardActionType.GUARD) {
      guardsRef.current.data[GuardPoint['eval.experiments.copy']].preprocess(
        action,
      );
    } else {
      action();
    }
  };

  useEffect(() => {
    const actionsColumn: ColumnProps<Experiment> = {
      title: I18n.t('operation'),
      disableColumnManage: true,
      dataIndex: 'action',
      key: 'action',
      fixed: 'right',
      align: 'right',
      width: 176,
      render: (_: unknown, record: Experiment) => {
        const hideRun =
          !isExperimentFail(record.status) ||
          actionVisibleControl?.retry === false;
        const actions: TableColAction[] = [
          {
            label: (
              <Tooltip content={I18n.t('re_evaluate_failed_only')} theme="dark">
                {I18n.t('retry')}
              </Tooltip>
            ),
            hide: hideRun,
            disabled: retryGuardType === GuardActionType.READONLY,
            onClick: () => handleRetryOnCLick(record),
          },
          {
            label: (
              <Tooltip content={I18n.t('view_detail')} theme="dark">
                {I18n.t('detail')}
              </Tooltip>
            ),
            onClick: () => handleDetailOnClick(record),
          },
          {
            label: (
              <Tooltip
                content={I18n.t('copy_and_create_experiment')}
                theme="dark"
              >
                {I18n.t('copy')}
              </Tooltip>
            ),
            hide: actionVisibleControl?.copy === false,
            disabled: copyGuardType === GuardActionType.READONLY,
            onClick: () => handleCopyOnClick(record),
          },
        ];
        // 收起来的操作
        const shrinkActions: TableColAction[] = [
          {
            label: I18n.t('delete'),
            type: 'danger',
            hide: actionVisibleControl?.delete === false,
            disabled: deleteGuardType === GuardActionType.READONLY,
            onClick: () => handleDelete({ record, spaceID, onRefresh }),
          },
        ];
        const maxCount = actions.filter(item => !item.hide).length;
        return (
          <TableColActions
            actions={[...actions, ...shrinkActions]}
            maxCount={maxCount}
          />
        );
      },
    };
    const idColumn: ColumnProps<Experiment> = {
      title: 'ID',
      disableColumnManage: true,
      dataIndex: 'id',
      key: 'id',
      width: 110,
      render(val: Int64) {
        return <IDRender id={val} useTag={true} />;
      },
    };
    const newColumns: ColumnProps<Experiment>[] = [
      ...(enableIdColumn ? [idColumn] : []),
      ...getExperimentColumns({ spaceID, enableSort }),
    ];

    setColumns([
      ...dealColumnsFromStorage(newColumns, columnManageStorageKey),
      ...(enableActionColumn ? [actionsColumn] : []),
    ]);
    setDefaultColumns([
      ...newColumns,
      ...(enableActionColumn ? [actionsColumn] : []),
    ]);
  }, [
    spaceID,
    copyGuardType,
    retryGuardType,
    deleteGuardType,
    actionVisibleControl,
  ]);

  return {
    columns,
    defaultColumns,
    setColumns,
  };
}
