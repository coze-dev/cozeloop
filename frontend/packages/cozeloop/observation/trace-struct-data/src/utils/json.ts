// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import JSONBig from 'json-bigint';
import { logger } from '@coze-arch/logger';

const jsonBig = JSONBig({ storeAsString: true });

export const safeJsonParse = (json: string): object | string => {
  try {
    return JSON.parse(JSON.stringify(jsonBig.parse(json)));
  } catch (e) {
    logger.error({ error: e as unknown as Error });
    return json;
  }
};

export const beautifyJson = (data: object) => JSON.stringify(data, null, 2);
