// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable max-params */
/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/naming-convention */
import type {
  LocaleData,
  I18nOptionsMap,
  I18nKeysHasOptionsType,
  I18nKeysNoOptionsType,
} from '@cozeloop/loop-lng';
import {
  type Intl,
  type I18nCore,
  type IIntlInitOptions,
  I18n as _I18n,
} from '@coze-arch/i18n/intl';

type Callback = Parameters<(typeof _I18n)['init']>[1];
type FallbackLng = ReturnType<(typeof _I18n)['getLanguages']>;
type IntlModule = Parameters<(typeof _I18n)['use']>[0];
type InitReturnType = ReturnType<(typeof _I18n)['init']>;

type I18nOptions<K extends LocaleData> = K extends keyof I18nOptionsMap
  ? I18nOptionsMap[K]
  : never;

class IntlX {
  plugins: any[] = [];
  public i18nInstance: I18nCore;
  constructor() {
    this.i18nInstance = _I18n.i18nInstance;
  }

  init(config: IIntlInitOptions, callback?: Callback): InitReturnType {
    return _I18n.init(config, callback);
  }

  use(plugin: IntlModule): Intl {
    return _I18n.use(plugin);
  }

  get language(): string {
    return _I18n.language;
  }

  setLangWithPromise(lng: string) {
    return this.i18nInstance.changeLanguageWithPromise(lng);
  }

  setLang(lng: string, callback?: Callback): void {
    return _I18n.setLang(lng, callback);
  }

  getLanguages(): FallbackLng {
    return _I18n.getLanguages();
  }

  dir(): 'ltr' | 'rtl' {
    return _I18n.dir();
  }

  addResourceBundle(
    lng: string,
    ns: string,
    resources: any,
    deep?: boolean,
    overwrite?: boolean,
  ) {
    return _I18n.addResourceBundle(lng, ns, resources, deep, overwrite);
  }

  t<K extends I18nKeysNoOptionsType>(
    keys: K,
    options?: Record<string, unknown>,
    fallbackText?: string,
  ): string;
  t<K extends I18nKeysHasOptionsType>(
    keys: K,
    options: I18nOptions<K>,
    fallbackText?: string,
  ): string;
  t<K extends LocaleData>(
    keys: K,
    options?: I18nOptions<K> | Record<string, unknown>,
    fallbackText?: string,
  ): string {
    return _I18n.t(keys, options, fallbackText);
  }
}

export const I18n = new IntlX();
