// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type PropsWithChildren } from 'react';

import { type GuardPoint } from './types';
import { useGuard } from './hooks/use-guard';

interface Props {
  point: GuardPoint;
}
export function GuardRoute({ point, children }: PropsWithChildren<Props>) {
  const { data } = useGuard({ point });

  if (data.readonly || data.hidden) {
    // 在独立包中，使用一个简单的提示组件代替 PageNoAuth
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-center p-4">
          <h2 className="text-xl font-bold mb-2">权限不足</h2>
          <p>您没有访问该页面的权限</p>
        </div>
      </div>
    );
  }
  return children;
}
