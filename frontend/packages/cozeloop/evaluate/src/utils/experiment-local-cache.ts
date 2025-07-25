// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { safeParseJson } from '@cozeloop/toolkit';
import { type AggregatorType } from '@cozeloop/api-schema/evaluation';

/** 实验详情本地缓存数据 */
export interface ExprimentDetailLocalCache {
  overviewAggregatorType?: AggregatorType;
  evaluatorVersionIds?: Int64[];
}

const EXPERIMENT_KEY_PREFIX = 'experiment_local_cache_';

/** 读取实验详情本地缓存数据 */
export function getExperimentDetailLocalCache(experimentID: Int64) {
  const key = `${EXPERIMENT_KEY_PREFIX}${experimentID}`;
  const cacheStr = localStorage.getItem(key);
  const cache: ExprimentDetailLocalCache | undefined =
    safeParseJson(cacheStr ?? '') || undefined;
  return cache;
}

/** 设置实验详情本地缓存数据 */
export function setExperimentDetailLocalCache(
  experimentID: Int64,
  data: ExprimentDetailLocalCache,
  mode: 'merge' | 'replace' = 'merge',
) {
  const newData =
    mode === 'merge'
      ? {
          ...getExperimentDetailLocalCache(experimentID),
          ...data,
        }
      : data;
  localStorage.setItem(
    `${EXPERIMENT_KEY_PREFIX}${experimentID}`,
    JSON.stringify(newData),
  );
}
