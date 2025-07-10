// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import {
  localeEnUS as loopLocaleEnUS,
  localeZhCN as loopLocaleZhCN,
} from '@cozeloop/loop-lng';
import { intlClient, type IntlClientOptions } from '@cozeloop/intl';

import { type CozeloopI18n } from './types';

// eslint-disable-next-line @typescript-eslint/naming-convention -- skip
const I18n: CozeloopI18n = intlClient;

/**
 * initialize I18n
 */
async function initIntl(options: IntlClientOptions = {}) {
  await intlClient.init({
    ...options,
    detection: {
      order: ['querystring', 'cookie', 'localStorage', 'navigator', 'htmlTag'],
      lookupQuerystring: 'locale',
      lookupCookie: 'locale',
      lookupLocalStorage: 'locale',
      caches: ['cookie'],
    },
    resources: {
      'zh-CN': {
        translation: Object.assign({}, loopLocaleZhCN),
      },
      'en-US': {
        translation: Object.assign({}, loopLocaleEnUS),
      },
    },
  });
}

export { I18n, initIntl };
