// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useEffect, useMemo, useRef, useState } from 'react';

import { useRequest } from 'ahooks';
import { EvaluatorSelectLocalData } from '@cozeloop/evaluate-components';
import { type Experiment } from '@cozeloop/api-schema/evaluation';
import { Spin } from '@coze-arch/coze-design';

import {
  getExperimentDetailLocalCache,
  setExperimentDetailLocalCache,
} from '@/utils/experiment-local-cache';
import { batchGetExperimentAggrResult } from '@/request/experiment';

import EvaluatorsScoreChart from './evaluators-score-chart';
import { EvaluatorsDraggableCard } from './draggable-cards';
import { I18n } from '@cozeloop/i18n-adapter';

export default function ExperimentChart({
  spaceID,
  experiment,
  experimentID,
  loading = false,
}: {
  spaceID: string;
  experiment: Experiment | undefined;
  experimentID: string;
  loading?: boolean;
}) {
  const [selectedEvaluatorIds, setSelectedEvaluatorIds] = useState<Int64[]>([]);
  const currentExperimentRef = useRef<Experiment | undefined>();

  const service = useRequest(
    async () => {
      const res = await batchGetExperimentAggrResult({
        workspace_id: spaceID,
        experiment_ids: [experimentID ?? 0],
      });
      return res.expt_aggregate_result?.[0]?.evaluator_results ?? [];
    },
    { refreshDeps: [experimentID] },
  );

  const evaluatorAggregateResult = useMemo(
    () => Object.values(service.data ?? {}),
    [service.data],
  );

  const allEvaluatorIds = useMemo(
    () => experiment?.evaluators?.map(e => e.current_version?.id ?? '') ?? [],
    [experiment],
  );

  useEffect(() => {
    // 相同实验id不刷新评估器选中状态
    if (
      currentExperimentRef.current &&
      currentExperimentRef.current.id === experiment?.id
    ) {
      return;
    }
    currentExperimentRef.current = experiment;
    const evaluatorVersionIds =
      getExperimentDetailLocalCache(experimentID)?.evaluatorVersionIds ??
      experiment?.evaluators?.map(e => e.current_version?.id ?? '') ??
      [];
    setSelectedEvaluatorIds(evaluatorVersionIds);
  }, [experiment]);

  return (
    <div className=" flex flex-col gap-4">
      <Spin spinning={loading || service.loading}>
        <div className="flex items-center text-sm font-semibold mb-3 h-[32px]">
          总览
        </div>
        <EvaluatorsScoreChart
          selectedEvaluatorIds={allEvaluatorIds}
          evaluatorAggregateResult={evaluatorAggregateResult}
          ready={!service.loading && !loading}
          experiment={experiment}
          spaceID={spaceID}
        />
        <div className="flex items-center gap-2 h-[32px] mt-5 mb-3">
          <div className="text-sm font-semibold">
            {I18n.t('score_details_data_item_distribution')}
          </div>
          <EvaluatorSelectLocalData
            multiple
            maxTagCount={1}
            prefix={I18n.t('indicator')}
            placeholder={I18n.t('please_select', {
              field: I18n.t('indicator'),
            })}
            className="ml-auto"
            style={{ minWidth: 220 }}
            evaluators={experiment?.evaluators}
            value={selectedEvaluatorIds}
            onChange={val => {
              setSelectedEvaluatorIds(val as Int64[]);
              setExperimentDetailLocalCache(experimentID, {
                evaluatorVersionIds: val as Int64[],
              });
            }}
          />
        </div>
        <EvaluatorsDraggableCard
          spaceID={spaceID}
          ready={!service.loading && !loading}
          evaluators={
            experiment?.evaluators?.filter(item =>
              selectedEvaluatorIds.find(id => id === item.current_version?.id),
            ) ?? []
          }
          evaluatorAggregateResult={evaluatorAggregateResult}
        />
      </Spin>
    </div>
  );
}
