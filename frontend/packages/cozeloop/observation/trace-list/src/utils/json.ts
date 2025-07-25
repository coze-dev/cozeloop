// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import JSONBig from 'json-bigint';

const jsonBig = JSONBig({ storeAsString: true });

export const jsonFormat = (json: string): object | string => {
  try {
    return JSON.parse(JSON.stringify(jsonBig.parse(json)));
  } catch (e) {
    return json;
  }
};

export const decodeJSON = <T>(jsonStr: string) =>
  JSON.parse(decodeURIComponent(jsonStr)) as T;

export const encodeJSON = <T>(json: T) =>
  encodeURIComponent(JSON.stringify(json));
