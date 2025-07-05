// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import localeZhCN from './locales/zh-CN.json';
import localeEn from './locales/en-US.json';

const defaultConfig = {
  en: { i18n: localeEn },
  'zh-CN': { i18n: localeZhCN },
};

export { localeEn, localeZhCN, defaultConfig };
export type {
  I18nOptionsMap,
  I18nKeysHasOptionsType,
  I18nKeysNoOptionsType,
  LocaleData,
} from './locale-data';
