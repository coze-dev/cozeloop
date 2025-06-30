// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { RouterProvider, createBrowserRouter } from 'react-router-dom';
import { Suspense } from 'react';

import { I18n } from '@cozeloop/i18n-adapter';
import { PageLoading } from '@cozeloop/components';
import { zh_CN } from '@coze-arch/coze-design/locales';
import { CDLocaleProvider } from '@coze-arch/coze-design';

import { routeConfig } from './routes';

import './index.css';

const router = createBrowserRouter(routeConfig);

export function App() {
  return (
    <Suspense fallback={<PageLoading />}>
      <CDLocaleProvider locale={zh_CN} i18n={I18n}>
        <RouterProvider router={router} />
      </CDLocaleProvider>
    </Suspense>
  );
}
