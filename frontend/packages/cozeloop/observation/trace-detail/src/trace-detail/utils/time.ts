// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import dayjs from '@/utils/dayjs';

const offsetTime = 30;

export const getEndTime = (
  startTime: number | string,
  latency: number | string,
) =>
  dayjs(Number(startTime))
    .add(Number(latency) + 1000, 'millisecond')
    .add(offsetTime, 'minute')
    .valueOf()
    .toString();

export const getStartTime = (startTime: number | string) =>
  dayjs(Number(startTime)).subtract(offsetTime, 'minute').valueOf().toString();
