// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import {
  IconCozInfoCircle,
  IconCozWarningCircleFillPalette,
} from '@coze-arch/coze-design/icons';
import { Banner, Tooltip } from '@coze-arch/coze-design';

import { type ExperimentItem } from '@/types/experiment';
import { useExperiment } from '@/hooks/use-experiment';
import { ActualOutputWithTrace } from '@/components/experiment';

export default function EvalActualOutputTable({
  item,
  expand,
}: {
  item: ExperimentItem;
  expand?: boolean;
}) {
  const experiment = useExperiment();
  return (
    <div className="text-sm py-3">
      <div className="flex items-center gap-1 mt-2 mb-3 px-5">
        <div className="font-medium text-xs">actual_output</div>
        <Tooltip
          theme="dark"
          content={I18n.t('evaluation_object_actual_output')}
        >
          <IconCozInfoCircle className="text-[var(--coz-fg-secondary)] hover:text-[var(--coz-fg-primary)]" />
        </Tooltip>
      </div>
      {item.targetErrorMsg ? (
        <Banner
          type="danger"
          className="rounded-small !px-3 !py-2"
          fullMode={false}
          icon={
            <div className="h-[22px] flex items-center">
              <IconCozWarningCircleFillPalette className="text-[16px] text-[rgb(var(--coze-red-5))]" />
            </div>
          }
          description={item.targetErrorMsg}
        />
      ) : (
        <div className="px-5">
          <ActualOutputWithTrace
            expand={expand}
            content={item.actualOutput}
            traceID={item?.evalTargetTraceID}
            displayFormat={true}
            startTime={experiment?.start_time}
            endTime={experiment?.end_time}
          />
        </div>
      )}
    </div>
  );
}
