// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
export function formatTimeDuration(time: number) {
  if (time < 1000) {
    return `${time}ms`;
  } else if (time < 60000) {
    return `${(time / 1000).toFixed(2)}s`;
  } else if (time < 3600000) {
    return `${(time / 60000).toFixed(2)}min`;
  } else if (time < 86400000) {
    return `${(time / 3600000).toFixed(2)}h`;
  } else {
    return `${(time / 86400000).toFixed(2)}d`;
  }
}
