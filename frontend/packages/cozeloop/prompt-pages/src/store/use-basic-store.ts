// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type Dispatch, type SetStateAction } from 'react';

import { create } from 'zustand';
import { nanoid } from 'nanoid';
import { produce } from 'immer';
interface BasicState {
  readonly?: boolean;
  saveLock?: boolean;
  autoSaving?: boolean;
  streaming?: boolean;
  versionChangeVisible?: boolean;
  versionChangeLoading?: boolean;
  debugId?: Int64;
  executeHistoryVisible?: boolean;
  optimizeEditorKey: string;
}

type PromptActionType<S> = Dispatch<SetStateAction<S>>;
interface BasicAction {
  setReadonly: PromptActionType<boolean | undefined>;
  setAutoSaving: PromptActionType<boolean | undefined>;
  setStreaming: PromptActionType<boolean | undefined>;
  setVersionChangeVisible: PromptActionType<boolean | undefined>;
  setVersionChangeLoading: PromptActionType<boolean | undefined>;
  setDebugId: PromptActionType<Int64 | undefined>;
  setExecuteHistoryVisible: PromptActionType<boolean | undefined>;
  setSaveLock: PromptActionType<boolean | undefined>;
  clearStore: () => void;
}

export const useBasicStore = create<BasicState & BasicAction>()((set, get) => ({
  autoSaving: false,
  optimizeEditorKey: nanoid(),
  saveLock: true,
  setReadonly: (val: SetStateAction<boolean | undefined>) =>
    set(
      produce((state: BasicState) => {
        state.readonly = val instanceof Function ? val(get().readonly) : val;
      }),
    ),
  setAutoSaving: (val: SetStateAction<boolean | undefined>) =>
    set(
      produce((state: BasicState) => {
        state.autoSaving =
          val instanceof Function ? val(get().autoSaving) : val;
      }),
    ),
  streaming: false,
  setStreaming: (val: SetStateAction<boolean | undefined>) =>
    set(
      produce((state: BasicState) => {
        state.streaming = val instanceof Function ? val(get().streaming) : val;
      }),
    ),
  versionChangeVisible: false,
  setVersionChangeVisible: (val: SetStateAction<boolean | undefined>) =>
    set(
      produce((state: BasicState) => {
        state.versionChangeVisible =
          val instanceof Function ? val(get().versionChangeVisible) : val;
      }),
    ),
  versionChangeLoading: false,
  setVersionChangeLoading: (val: SetStateAction<boolean | undefined>) =>
    set(
      produce((state: BasicState) => {
        state.versionChangeLoading =
          val instanceof Function ? val(get().versionChangeLoading) : val;
      }),
    ),
  debugId: undefined,
  setDebugId: (val: SetStateAction<Int64 | undefined>) =>
    set(
      produce((state: BasicState) => {
        state.debugId = val instanceof Function ? val(get().debugId) : val;
      }),
    ),
  executeHistoryVisible: false,
  setExecuteHistoryVisible: (val: SetStateAction<boolean | undefined>) =>
    set(
      produce((state: BasicState) => {
        state.executeHistoryVisible =
          val instanceof Function ? val(get().executeHistoryVisible) : val;
      }),
    ),
  setSaveLock: (val: SetStateAction<boolean | undefined>) =>
    set(
      produce((state: BasicState) => {
        state.saveLock = val instanceof Function ? val(get().saveLock) : val;
      }),
    ),
  clearStore: () =>
    set({
      debugId: undefined,
      autoSaving: false,
      streaming: false,
      versionChangeVisible: false,
      versionChangeLoading: false,
      executeHistoryVisible: false,
      saveLock: true,
      optimizeEditorKey: nanoid(),
    }),
}));
