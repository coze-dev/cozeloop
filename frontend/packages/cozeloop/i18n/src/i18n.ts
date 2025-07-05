// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import LanguageDetector from 'i18next-browser-languagedetector';
import {
  localeEn as loopLocaleEn,
  localeZhCN as loopLocaleZhCN,
} from '@cozeloop/loop-lng';
import { type IIntlInitOptions } from '@coze-arch/i18n/intl';

import { detectLng } from './utils';
import { I18n } from './intl';

/**
 * initialize I18n
 */
async function initIntl(options?: IIntlInitOptions) {
  const { lng = detectLng('zh-CN'), ...restOptions } = options || {};

  return new Promise(resolve => {
    I18n.use(LanguageDetector).init(
      {
        lng,
        detection: {
          order: [
            'querystring',
            'cookie',
            'localStorage',
            'navigator',
            'htmlTag',
          ],
          lookupQuerystring: 'lng',
          lookupCookie: 'i18next',
          lookupLocalStorage: 'i18next',
          caches: ['cookie'],
        },
        react: {
          useSuspense: false,
        },
        keySeparator: false,
        resources: {
          'zh-CN': {
            translation: Object.assign({}, loopLocaleZhCN),
          },
          'en-US': {
            translation: Object.assign({}, loopLocaleEn),
          },
        },
        ...(restOptions ?? {}),
      },
      resolve,
    );
  });
}

export { I18n, initIntl };
