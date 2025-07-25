// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useEffect } from 'react';

import { $notification } from '@cozeloop/api-schema';
import { Toast } from '@coze-arch/coze-design';

export function useApiErrorToast() {
  useEffect(() => {
    const onApiError = (msg: string) => {
      Toast.error({
        className: 'api-error-toast',
        content: (
          <span className="inline-block max-w-[100%] break-all whitespace-normal">
            {msg}
          </span>
        ),
      });
    };

    $notification.addListener('apiError', onApiError);

    return () => {
      $notification.removeListener('apiError', onApiError);
    };
  }, []);
}
