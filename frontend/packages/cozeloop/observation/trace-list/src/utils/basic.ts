// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { isNil, isNumber } from 'lodash-es';

export const textWithFallback = (text?: string | number) =>
  isNil(text) || text === '' || (isNumber(text) && !isFinite(text))
    ? '-'
    : text.toString();

export const formatNumberWithCommas = (number?: string | number) =>
  number ? number.toString().replace(/(\d)(?=(?:\d{3})+$)/g, '$1,') : number;
