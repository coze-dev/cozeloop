// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @coze-arch/max-line-per-function */
/* eslint-disable max-lines-per-function */
/* eslint-disable complexity */
import { useEffect, useState } from 'react';

import { useShallow } from 'zustand/react/shallow';
import { nanoid } from 'nanoid';
import { debounce, isEqual } from 'lodash-es';
import { useSpace } from '@cozeloop/biz-hooks-adapter';
import { Role } from '@cozeloop/api-schema/prompt';

import { getPromptStorageInfo, setPromptStorageInfo } from '@/utils/prompt';
import { type PromptState, usePromptStore } from '@/store/use-prompt-store';
import {
  type PromptMockDataState,
  usePromptMockDataStore,
} from '@/store/use-mockdata-store';
import { useBasicStore } from '@/store/use-basic-store';
import { CALL_SLEEP_TIME, PromptStorageKey } from '@/consts';

type PlaygroundInfoStorage = Record<string, PromptState>;

type PlaygroundMockSetStorage = Record<string, PromptMockDataState>;

export const usePlayground = () => {
  const { spaceID } = useSpace();
  const {
    setPromptInfo,
    setMessageList,
    setModelConfig,
    setToolCallConfig,
    setTools,
    setVariables,
    setCurrentModel,
    clearStore: clearPromptStore,
  } = usePromptStore(
    useShallow(state => ({
      setPromptInfo: state.setPromptInfo,
      setMessageList: state.setMessageList,
      setModelConfig: state.setModelConfig,
      setToolCallConfig: state.setToolCallConfig,
      setTools: state.setTools,
      setVariables: state.setVariables,
      setCurrentModel: state.setCurrentModel,
      clearStore: state.clearStore,
    })),
  );
  const { setAutoSaving, clearStore: clearBasicStore } = useBasicStore(
    useShallow(state => ({
      setAutoSaving: state.setAutoSaving,
      clearStore: state.clearStore,
    })),
  );
  const {
    setHistoricMessage,
    setMockVariables,
    setUserDebugConfig,
    clearMockDataStore,
    setCompareConfig,
  } = usePromptMockDataStore(
    useShallow(state => ({
      setHistoricMessage: state.setHistoricMessage,
      setMockVariables: state.setMockVariables,
      setUserDebugConfig: state.setUserDebugConfig,
      compareConfig: state.compareConfig,
      setCompareConfig: state.setCompareConfig,
      clearMockDataStore: state.clearMockDataStore,
    })),
  );

  const [initPlaygroundLoading, setInitPlaygroundLoading] = useState(true);

  useEffect(() => {
    setInitPlaygroundLoading(true);
    const storagePlaygroundInfo = getPromptStorageInfo<PlaygroundInfoStorage>(
      PromptStorageKey.PLAYGROUND_INFO,
    );
    const info = storagePlaygroundInfo?.[spaceID];

    const storagePlaygroundMockSet =
      getPromptStorageInfo<PlaygroundMockSetStorage>(
        PromptStorageKey.PLAYGROUND_MOCKSET,
      );
    const mock = storagePlaygroundMockSet?.[spaceID];

    if (mock) {
      setHistoricMessage(mock?.historicMessage || []);
      setMockVariables(mock?.mockVariables || []);
      setUserDebugConfig(mock?.userDebugConfig || {});
      setCompareConfig(mock?.compareConfig || {});
    }

    setTools(info?.tools || []);
    setModelConfig(info?.modelConfig || {});
    setToolCallConfig(info?.toolCallConfig || {});
    setVariables(info?.variables || []);
    setMessageList(
      info?.messageList || [{ role: Role.System, content: '', key: nanoid() }],
    );
    setCurrentModel(info?.currentModel || {});
    setPromptInfo(info?.promptInfo || { workspace_id: spaceID });

    setInitPlaygroundLoading(false);

    return () => {
      setInitPlaygroundLoading(true);
      setTimeout(() => {
        clearPromptStore();
        clearBasicStore();
        clearMockDataStore();
      }, 0);
    };
  }, [spaceID]);

  const saveMockSet = debounce((mockSet: PromptMockDataState, sID: string) => {
    const storagePlaygroundMockSet =
      getPromptStorageInfo<PlaygroundMockSetStorage>(
        PromptStorageKey.PLAYGROUND_MOCKSET,
      );
    setPromptStorageInfo<PlaygroundMockSetStorage>(
      PromptStorageKey.PLAYGROUND_MOCKSET,
      { ...storagePlaygroundMockSet, [sID]: mockSet },
    );
    setAutoSaving(false);
  }, CALL_SLEEP_TIME);

  const saveInfo = debounce((info: PromptState, sID: string) => {
    const storagePlaygroundInfo = getPromptStorageInfo<PlaygroundInfoStorage>(
      PromptStorageKey.PLAYGROUND_INFO,
    );
    setPromptStorageInfo<PlaygroundInfoStorage>(
      PromptStorageKey.PLAYGROUND_INFO,
      {
        ...storagePlaygroundInfo,
        [sID]: info,
      },
    );
    setAutoSaving(false);
  }, CALL_SLEEP_TIME);

  useEffect(() => {
    const dataSub = usePromptStore.subscribe(
      state => ({
        toolCallConfig: state.toolCallConfig,
        variables: state.variables,
        modelConfig: state.modelConfig,
        tools: state.tools,
        messageList: state.messageList,
        promptInfo: state.promptInfo,
        currentModel: state.currentModel,
      }),
      val => {
        if (!initPlaygroundLoading) {
          setAutoSaving(true);
          saveInfo(val, spaceID);
        }
      },
      {
        equalityFn: isEqual,
        fireImmediately: true, // 是否在第一次调用（初始化时）立刻执行
      },
    );
    const mockSub = usePromptMockDataStore.subscribe(
      state => ({
        historicMessage: state.historicMessage,
        userDebugConfig: state.userDebugConfig,
        mockVariables: state.mockVariables,
        compareConfig: state.compareConfig,
      }),
      val => {
        if (!initPlaygroundLoading) {
          setAutoSaving(true);
          saveMockSet(val, spaceID);
        }
      },
      {
        equalityFn: isEqual,
        fireImmediately: true, // 是否在第一次调用（初始化时）立刻执行
      },
    );

    return () => {
      dataSub?.();
      mockSub?.();
    };
  }, [initPlaygroundLoading, spaceID]);

  return {
    initPlaygroundLoading,
  };
};
