// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import {
  type I18nKeysNoInterpolation,
  type I18nWithInterpolation,
} from './locale-types';

interface I18nFunction {
  /** 🔵 I18n **with** interpolation */
  <K extends keyof I18nWithInterpolation>(
    keys: K,
    options: I18nWithInterpolation[K],
    fallbackText?: string,
  ): string;
  /** 🟣 I18n **without** interpolation */
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

export interface CozeloopI18n {
  t: I18nFunction;
  unsafeT: UnsafeI18nFunction;
}
