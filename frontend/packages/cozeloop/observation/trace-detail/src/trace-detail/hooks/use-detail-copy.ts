// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useCallback } from 'react';

import { handleCopy as copy } from '@cozeloop/components';
import { useSpace } from '@cozeloop/biz-hooks-adapter';

/** 上报埋点需要 */
export const useDetailCopy = (moduleName?: string) => {
  const { spaceID, space: { name } = {} } = useSpace();

  const handleCopy = useCallback(
    (text: string, point: string) => {
      copy(text);

      if (!moduleName) {
        return;
      }
    },
    [moduleName, spaceID, name],
  );
  return handleCopy;
};
