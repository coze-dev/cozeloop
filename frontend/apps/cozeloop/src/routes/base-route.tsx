// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { Outlet, Navigate } from 'react-router-dom';

import { GuardProvider } from '@cozeloop/guard';
import { PageLoading } from '@cozeloop/components';
import { useCheckLogin, useLoginStatus } from '@cozeloop/account';

import { useApiErrorToast } from '@/hooks';
import { LOGIN_PATH } from '@/constants';

export function BaseRoute() {
  const loginStatus = useLoginStatus();
  useCheckLogin();
  useApiErrorToast();

  switch (loginStatus) {
    case 'settling':
      return <PageLoading />;
    case 'not_login':
      return <Navigate to={LOGIN_PATH} />;
    case 'logined':
      return (
        <GuardProvider>
          <Outlet />
        </GuardProvider>
      );
    default:
      return null;
  }
}
