// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @coze-arch/use-error-in-catch */
export const safeJsonParse = (value?: string | null) => {
  try {
    return JSON.parse(value || '');
  } catch (error) {
    return '';
  }
};
