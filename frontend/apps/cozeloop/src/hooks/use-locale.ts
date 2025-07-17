// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { en_US, zh_CN } from '@coze-arch/coze-design/locales';

import { useI18nStore } from '@/stores';

function langToLocale(lang: string) {
  if (!lang) {
    return zh_CN;
  }
  switch (lang) {
    case 'zh':
    case 'zh-CN':
      return zh_CN;
    default:
      return en_US;
  }
}

export function useLocale() {
  const lng = useI18nStore(s => s.lng);

  return langToLocale(lng);
}
