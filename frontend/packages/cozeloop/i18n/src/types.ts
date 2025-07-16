// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import {
  type I18nKeysNoInterpolation,
  type I18nWithInterpolation,
} from './locale-types';

import { type intlClient } from '@cozeloop/intl';

interface I18nFunction {
  /** 🔵 I18n **with** interpolation */
  <K extends keyof I18nWithInterpolation>(
    keys: K,
    options: I18nWithInterpolation[K],
    fallbackText?: string,
  ): string;
  /** 🟣 I18n **without** interpolation {@link I18nKeysNoInterpolation keys} */
  <K extends I18nKeysNoInterpolation>(keys: K, fallbackText?: string): string;
}

interface UnsafeI18nFunction {
  /** trust the key */
  (
    key: string,
    options?: Record<string, unknown>,
    fallbackText?: string,
  ): string;
}

export type CozeloopI18n = Omit<typeof intlClient, 't' | 'unsafeT'> & {
  t: I18nFunction;
  unsafeT: UnsafeI18nFunction;
};
