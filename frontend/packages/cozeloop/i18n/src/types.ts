import {
  type I18nKeysNoInterpolation,
  type I18nWithInterpolation,
} from './locale-types';

interface I18nFunction {
  <K extends keyof I18nWithInterpolation>(
    keys: K,
    options: I18nWithInterpolation[K],
    fallbackText?: string,
  ): string;
  <K extends I18nKeysNoInterpolation>(keys: K, fallbackText?: string): string;
  (
    keys: string,
    options?: Record<string, unknown>,
    fallbackText?: string,
  ): string;
}

export interface CozeloopI18n {
  t: I18nFunction;
}
