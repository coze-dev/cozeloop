// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useEffect } from 'react';

import { useShallow } from 'zustand/react/shallow';
import { useUIStore, type BreadcrumbItemConfig } from '@cozeloop/stores';

export function useBreadcrumb(config: BreadcrumbItemConfig) {
  const { pushBreadcrumb, popBreadcrumb } = useUIStore(
    useShallow(store => ({
      pushBreadcrumb: store.pushBreadcrumb,
      popBreadcrumb: store.popBreadcrumb,
    })),
  );

  useEffect(() => {
    pushBreadcrumb(config);
    return () => {
      popBreadcrumb();
    };
  }, [config]);
}
