// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useNavigate, useLocation } from 'react-router-dom';
import type React from 'react';
import { useMemo, useRef } from 'react';

import qs from 'query-string';
import type { ParseOptions, StringifyOptions } from 'query-string';
import { useMemoizedFn, useUpdate } from 'ahooks';

export interface Options {
  navigateMode?: 'push' | 'replace';
  parseOptions?: ParseOptions;
  stringifyOptions?: StringifyOptions;
}

const baseParseConfig: ParseOptions = {
  parseNumbers: false,
  parseBooleans: false,
  arrayFormat: 'bracket',
};

const baseStringifyConfig: StringifyOptions = {
  skipNull: false,
  skipEmptyString: false,
  arrayFormat: 'bracket',
};

type UrlState = Record<string, unknown>;

export const useUrlState = <S extends UrlState = UrlState>(
  initialState?: S | (() => S),
  options?: Options,
) => {
  type State = S;
  const {
    navigateMode = 'replace',
    parseOptions,
    stringifyOptions,
  } = options || {};

  const mergedParseOptions = { ...baseParseConfig, ...parseOptions };
  const mergedStringifyOptions = {
    ...baseStringifyConfig,
    ...stringifyOptions,
  };

  const location = useLocation();
  const navigate = useNavigate();
  const update = useUpdate();

  const initialStateRef = useRef(
    typeof initialState === 'function'
      ? (initialState as () => S)()
      : initialState || {},
  );

  const queryFromUrl = useMemo(
    () => qs.parse(location.search, mergedParseOptions),
    [location.search],
  );

  const targetQuery = useMemo(
    () => ({
      ...initialStateRef.current,
      ...queryFromUrl,
    }),
    [queryFromUrl],
  ) as State;

  const setState = (s: React.SetStateAction<State>) => {
    const newQuery = typeof s === 'function' ? s(targetQuery) : s;
    // 1. 如果 setState 后，search 没变化，就需要 update 来触发一次更新。
    // 2. update 和 history 的更新会合并，不会造成多次更新
    update();
    navigate(
      {
        hash: location.hash,
        search:
          qs.stringify(
            { ...queryFromUrl, ...newQuery },
            mergedStringifyOptions,
          ) || '?',
      },
      {
        replace: navigateMode === 'replace',
        state: location.state,
      },
    );
  };

  return [targetQuery, useMemoizedFn(setState)] as const;
};
