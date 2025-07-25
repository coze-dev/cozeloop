// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useState, useEffect, useRef } from 'react';

import { useDocumentVisibility } from 'ahooks';
import { sendEvent, EVENT_NAMES } from '@cozeloop/tea-adapter';
import { useSpace } from '@cozeloop/biz-hooks-adapter';

export function usePageStay() {
  const visibleState = useDocumentVisibility();
  const stayTimeRef = useRef<number>(0);
  const startTimeRef = useRef<number | null>(null);
  const [stayTime, setStayTime] = useState<number>(0);
  const { spaceID, space: { name: spaceName } = {} } = useSpace();

  useEffect(() => {
    if (visibleState !== 'hidden') {
      startTimeRef.current = Date.now();
    } else if (startTimeRef.current) {
      stayTimeRef.current = Date.now() - startTimeRef.current;
      startTimeRef.current = null;
      setStayTime(stayTimeRef.current);

      sendEvent(EVENT_NAMES.cozeloop_observation_trace_page_stay, {
        duration: stayTimeRef.current,
        space_id: spaceID,
        space_name: spaceName || '',
      });
    }

    return () => {
      if (startTimeRef.current) {
        stayTimeRef.current = Date.now() - startTimeRef.current;
        startTimeRef.current = null;
        setStayTime(stayTimeRef.current);

        sendEvent(EVENT_NAMES.cozeloop_observation_trace_page_stay, {
          duration: stayTimeRef.current,
          space_id: spaceID,
          space_name: spaceName || '',
        });
      }
    };
  }, [visibleState]);

  return stayTime;
}
