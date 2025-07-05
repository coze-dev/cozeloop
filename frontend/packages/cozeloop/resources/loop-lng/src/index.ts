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
