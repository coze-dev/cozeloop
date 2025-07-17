// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { RouterProvider, createBrowserRouter } from 'react-router-dom';
import { Suspense } from 'react';

import { I18n } from '@cozeloop/i18n-adapter';
import { PageLoading } from '@cozeloop/components';
import { CDLocaleProvider } from '@coze-arch/coze-design';

import { useI18nStore } from './stores';
import { routeConfig } from './routes';
import { useLocale, useSetupI18n } from './hooks';

import './index.css';

const router = createBrowserRouter(routeConfig);

export function App() {
  useSetupI18n();
  const lng = useI18nStore(s => s.lng);
  const locale = useLocale();

  return (
    <Suspense fallback={<PageLoading />}>
      <CDLocaleProvider key={lng} locale={locale} i18n={I18n}>
        <RouterProvider router={router} />
      </CDLocaleProvider>
    </Suspense>
  );
}
