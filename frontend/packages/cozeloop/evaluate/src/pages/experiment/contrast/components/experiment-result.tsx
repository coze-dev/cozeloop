// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useMemo } from 'react';

import { EVENT_NAMES, sendEvent } from '@cozeloop/tea-adapter';
import {
  AutoOverflowList,
  EvaluatorNameScore,
} from '@cozeloop/evaluate-components';
import {
  type Evaluator,
  type Experiment,
  type ExperimentTurnPayload,
  type EvaluatorRecord,
} from '@cozeloop/api-schema/evaluation';

import {
  ActualOutputWithTrace,
  ExperimentRunDataSummary,
} from '@/components/experiment';

export interface ExperimentContrastResultProps {
  result: ExperimentTurnPayload | undefined;
  experiment: Experiment | undefined;
  expand?: boolean;
  hiddenFieldMap?: Record<string, boolean>;
  spaceID?: Int64;
  onRefresh?: () => void;
}

interface Item extends Evaluator {
  evaluatorRecord: EvaluatorRecord | undefined;
}

export default function ExperimentResult({
  result,
  experiment,
  expand,
  spaceID,
  hiddenFieldMap = {},
  onRefresh,
}: ExperimentContrastResultProps) {
  const items = useMemo(
    () =>
      experiment?.evaluators?.map(evaluator => {
        const evaluatorRecord =
          result?.evaluator_output?.evaluator_records?.[
            evaluator.current_version?.id ?? ''
          ];
        return { ...evaluator, evaluatorRecord };
      }) ?? [],
    [experiment?.evaluators, result],
  );
  const actualOutput =
    result?.target_output?.eval_target_record?.eval_target_output_data
      ?.output_fields?.actual_output;
  const targetTraceID = result?.target_output?.eval_target_record?.trace_id;
  const onReportCalibration = () => {
    sendEvent(EVENT_NAMES.cozeloop_experiment_detailsdrawer_editscore, {
      from: 'experiment_contrast_result',
    });
  };
  const onReportEvaluatorTrace = () => {
    sendEvent(EVENT_NAMES.cozeloop_experiment_detailsdrawer_trace, {
      from: 'experiment_contrast_result',
    });
  };
  return (
    <div className="flex flex-col gap-2" onClick={e => e.stopPropagation()}>
      <ActualOutputWithTrace
        expand={expand}
        content={actualOutput}
        traceID={targetTraceID}
        startTime={experiment?.start_time}
        endTime={experiment?.end_time}
      />
      <AutoOverflowList<Item>
        itemKey={'current_version.id'}
        items={items}
        itemRender={({ item, inOverflowPopover }) => {
          const { evaluatorRecord } = item;
          const evaluatorResult =
            evaluatorRecord?.evaluator_output_data?.evaluator_result;
          return (
            <EvaluatorNameScore
              key={item.current_version?.id}
              evaluator={item}
              evaluatorResult={evaluatorResult}
              experiment={experiment}
              updateUser={evaluatorRecord?.base_info?.updated_by}
              spaceID={spaceID}
              traceID={evaluatorRecord?.trace_id}
              evaluatorRecordID={evaluatorRecord?.id}
              enablePopover={!inOverflowPopover}
              enableEditScore={false}
              border={!inOverflowPopover}
              showVersion={true}
              defaultShowAction={inOverflowPopover}
              onEditScoreSuccess={onRefresh}
              onReportCalibration={onReportCalibration}
              onReportEvaluatorTrace={onReportEvaluatorTrace}
            />
          );
        }}
      />
      <ExperimentRunDataSummary
        result={result}
        latencyHidden={true}
        tokenHidden={true}
        statusHidden={hiddenFieldMap.status}
      />
    </div>
  );
}
