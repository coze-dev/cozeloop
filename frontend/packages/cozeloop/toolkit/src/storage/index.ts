// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type LocalStorageKeys, cacheConfig } from './config';
interface CozeLoopStorageConfig {
  // 业务域，唯一标识
  field: string;
}

export class CozeLoopStorage {
  private field: string;
  private static userID: string;

  static setUserID(id: string) {
    CozeLoopStorage.userID = id;
  }
  constructor(config: CozeLoopStorageConfig) {
    this.field = config.field;
  }

  private makeKey(key: LocalStorageKeys) {
    if (cacheConfig[key]?.bindAccount) {
      return `[${this.field}]:[${CozeLoopStorage.userID}]:${key}`;
    }
    return `[${this.field}]:${key}`;
  }
  setItem(key: LocalStorageKeys, value: string) {
    localStorage.setItem(`${this.makeKey(key)}`, value);
  }

  getItem(key: LocalStorageKeys) {
    return localStorage.getItem(`${this.makeKey(key)}`);
  }

  removeItem(key: LocalStorageKeys) {
    localStorage.removeItem(`${this.makeKey(key)}`);
  }

  getKey(key: LocalStorageKeys) {
    return this.makeKey(key);
  }
}
