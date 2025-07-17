// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useEffect } from 'react';

import { I18n } from '@cozeloop/i18n-adapter';

import { useI18nStore } from '@/stores';

export function useSetupI18n() {
  const setLng = useI18nStore(s => s.setLng);

  useEffect(() => {
    setLng(I18n.lang);
  }, []);
}
