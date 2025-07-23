// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { type PropsWithChildren } from 'react';

import { useI18nStore } from '@cozeloop/stores';
import { I18n } from '@cozeloop/i18n-adapter';
import { en_US, zh_CN } from '@coze-arch/coze-design/locales';
import {
  CDLocaleProvider,
  ConfigProvider,
  enUS,
  zhCN,
} from '@coze-arch/coze-design';

function langToLocale(lang: string) {
  if (!lang) {
    return { locale: zhCN, cdLocale: zh_CN };
  }
  switch (lang) {
    case 'zh':
    case 'zh-CN':
      return { locale: zhCN, cdLocale: zh_CN };
    default:
      return { locale: enUS, cdLocale: en_US };
  }
}

export function LocaleProvider({ children }: PropsWithChildren) {
  const lang = useI18nStore(s => s.lng);
  const { locale, cdLocale } = langToLocale(lang);

  return (
    <ConfigProvider locale={locale}>
      <CDLocaleProvider
        locale={cdLocale}
        i18n={{ t: (key: string) => I18n.unsafeT(key) }}
      >
        {children}
      </CDLocaleProvider>
    </ConfigProvider>
  );
}
