// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useEffect, useRef } from 'react';

import { useDocumentVisibility, useMemoizedFn, useUnmount } from 'ahooks';
import { type ParamsTypeDefine, sendEvent } from '@cozeloop/tea-adapter';

import { ViewDurationManager } from '@/utils/duration';

export const useTeaDuration = (
  event: string,
  option: ParamsTypeDefine[typeof event],
) => {
  const viewDurationManager = useRef(new ViewDurationManager()).current;

  const documentVisibility = useDocumentVisibility();

  // 上报耗时埋点
  const handleViewDuration = useMemoizedFn(() => {
    if (viewDurationManager.status !== 'finished') {
      const duration = viewDurationManager.finish();
      sendEvent(event, {
        ...option,
        duration,
      });
    }
  });

  // 组件卸载（离开）时上报埋点
  useUnmount(handleViewDuration);

  // 页面隐藏上报埋点
  useEffect(() => {
    if (
      documentVisibility === 'visible' &&
      viewDurationManager.status === 'finished'
    ) {
      viewDurationManager.reset();
    }

    if (
      documentVisibility === 'hidden' &&
      viewDurationManager.status === 'running'
    ) {
      handleViewDuration();
    }
  }, [documentVisibility]);

  // 刷新关闭页面上报埋点
  useEffect(() => {
    window.addEventListener('beforeunload', handleViewDuration);

    return () => {
      window.removeEventListener('beforeunload', handleViewDuration);
    };
  }, []);
};
