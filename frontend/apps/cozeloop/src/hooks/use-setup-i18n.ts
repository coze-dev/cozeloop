// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useEffect } from 'react';

import { useI18nStore } from '@cozeloop/stores';
import { I18n } from '@cozeloop/i18n-adapter';

export function useSetupI18n() {
  const setLng = useI18nStore(s => s.setLng);

  useEffect(() => {
    setLng(I18n.lang);
  }, []);
}
