// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useMemo } from 'react';

import { IconCozClockFill } from '@coze-arch/coze-design/icons';
import { Tag } from '@coze-arch/coze-design';

import { formatTimeDuration } from '@/utils/time';

import { CustomTableTooltip } from '../table-cell-text';

interface TimeTagProps {
  latency?: number; // in seconds
}

const FAST_LATENCY = 50000;
const SLOW_LATENCY = 600000;

export const LatencyTag: React.FC<TimeTagProps> = ({ latency }) => {
  const [bgClassName, textClassName] = useMemo(() => {
    if (Number(latency) > SLOW_LATENCY) {
      return ['!bg-[rgba(255,235,233,1)]', '!text-[rgba(208,41,47,1)]'];
    } else if (Number(latency) > FAST_LATENCY) {
      return ['!bg-[rgba(251,238,225,1)]', '!text-[rgba(160,95,1,1)]'];
    } else {
      return ['!bg-[rgba(230,247,237,1)]', '!text-[rgba(0,129,92,1)]'];
    }
  }, [latency]);

  if (!latency) {
    return <div className="flex items-center justify-start">-</div>;
  }

  return (
    <Tag
      size="small"
      className={`m-w-full border-box ${bgClassName}`}
      prefixIcon={<IconCozClockFill className={`${textClassName} !w-3 !h-3`} />}
    >
      <CustomTableTooltip textClassName={textClassName}>
        {formatTimeDuration(Number(latency))}
      </CustomTableTooltip>
    </Tag>
  );
};
