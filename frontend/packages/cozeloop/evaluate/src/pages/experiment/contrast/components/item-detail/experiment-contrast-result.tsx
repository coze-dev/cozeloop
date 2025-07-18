// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { EVENT_NAMES, sendEvent } from '@cozeloop/tea-adapter';
import { I18n } from '@cozeloop/i18n-adapter';
import {
  TraceTrigger,
  EvaluatorNameScore,
  useGlobalEvalConfig,
} from '@cozeloop/evaluate-components';
import {
  type Experiment,
  type ExperimentTurnPayload,
} from '@cozeloop/api-schema/evaluation';
import { IconCozInfoCircle } from '@coze-arch/coze-design/icons';
import { Divider, Tooltip } from '@coze-arch/coze-design';

import { CellContentRender } from '@/utils/experiment';
import { ExperimentRunDataSummary } from '@/components/experiment';

export default function ExperimentContrastResult({
  result,
  experiment,
  expand,
  spaceID,
  onRefresh,
}: {
  experiment: Experiment | undefined;
  result: ExperimentTurnPayload | undefined;
  expand?: boolean;
  spaceID?: Int64;
  onRefresh?: () => void;
}) {
  const { traceEvalTargetPlatformType } = useGlobalEvalConfig();
  const actualOutput =
    result?.target_output?.eval_target_record?.eval_target_output_data
      ?.output_fields?.actual_output;
  const targetTraceID = result?.target_output?.eval_target_record?.trace_id;
  const onReportCalibration = () => {
    sendEvent(EVENT_NAMES.cozeloop_experiment_detailsdrawer_editscore, {
      from: 'experiment_contrast_item_detail',
    });
  };
  const onReportEvaluatorTrace = () => {
    sendEvent(EVENT_NAMES.cozeloop_experiment_detailsdrawer_trace, {
      from: 'experiment_contrast_item_detail',
    });
  };
  return (
    <div className="group flex flex-col gap-2 h-full">
      <div className="flex gap-2 flex-wrap">
        {experiment?.evaluators?.map(item => {
          const evaluatorRecord =
            result?.evaluator_output?.evaluator_records?.[
              item.current_version?.id ?? ''
            ];
          // 评估器聚合结果
          const evaluatorResult =
            evaluatorRecord?.evaluator_output_data?.evaluator_result;
          return (
            <EvaluatorNameScore
              evaluator={item}
              evaluatorResult={evaluatorResult}
              experiment={experiment}
              updateUser={evaluatorRecord?.base_info?.updated_by}
              spaceID={spaceID}
              traceID={evaluatorRecord?.trace_id}
              evaluatorRecordID={evaluatorRecord?.id}
              enablePopover={true}
              showVersion={true}
              onEditScoreSuccess={onRefresh}
              onReportCalibration={onReportCalibration}
              onReportEvaluatorTrace={onReportEvaluatorTrace}
            />
          );
        })}
      </div>
      <ExperimentRunDataSummary
        result={result}
        latencyHidden={true}
        tokenHidden={true}
      />
      <Divider />
      <div className="flex items-center gap-1">
        <div className="text-[var(--coz-fg-secondary)]">actual_output</div>
        <Tooltip
          theme="dark"
          content={I18n.t('evaluation_object_actual_output')}
        >
          <IconCozInfoCircle className="text-[var(--coz-fg-secondary)] hover:text-[var(--coz-fg-primary)]" />
        </Tooltip>
      </div>

      <div className="group flex leading-5 w-full grow min-h-[20px] overflow-hidden">
        <CellContentRender
          expand={expand}
          content={actualOutput}
          displayFormat={true}
          className="!max-h-[none]"
        />
        {targetTraceID ? (
          <div className="flex ml-auto" onClick={e => e.stopPropagation()}>
            <TraceTrigger
              className="ml-1 invisible group-hover:visible"
              traceID={targetTraceID ?? ''}
              platformType={traceEvalTargetPlatformType}
              startTime={experiment?.start_time}
              endTime={experiment?.end_time}
            />
          </div>
        ) : null}
      </div>
    </div>
  );
}
