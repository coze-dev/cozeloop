// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import dayjs from 'dayjs';

const UNIX_LEN = 10;
const UNIX_LEN2 = 13;

export const formatTimestampToString = (
  timestamp: string | number,
  format = 'YYYY-MM-DD HH:mm:ss',
) => {
  const strLen = `${timestamp}`.length;
  if (strLen === UNIX_LEN) {
    return dayjs.unix(Number(timestamp)).format(format);
  } else if (strLen === UNIX_LEN2) {
    return dayjs(Number(timestamp)).format(format);
  }
  return '-';
};

export const safeParseJson = <T>(
  jsonString?: string,
  fallback?: T,
): T | undefined => {
  try {
    if (jsonString) {
      return JSON.parse(jsonString) as T;
    }
  } catch (e) {
    return fallback;
  }
};

export const formateMsToSeconds = (ms?: number | string) => {
  if (ms === undefined || ms === null) {
    return '-';
  }
  if (Number(ms) < 100) {
    return `${ms}ms`;
  }
  return `${(Number(ms) / 1000).toFixed(2)}s`;
};
