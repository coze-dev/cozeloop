// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type PropsWithChildren } from 'react';

import { Navbar } from '../navbar';
import { MainBreadcrumb } from '../breadcrumb';

export function TemplateLayout({ children }: PropsWithChildren) {
  return (
    <div className="relative h-full min-h-0 flex-shrink flex overflow-y-hidden">
      <Navbar />
      <div className="flex flex-col flex-1 overflow-hidden coz-bg-plus">
        <MainBreadcrumb />
        <div className="flex-1 overflow-x-auto overflow-y-hidden min-h-0">
          <div className="min-w-[960px] h-full">{children}</div>
        </div>
      </div>
    </div>
  );
}
