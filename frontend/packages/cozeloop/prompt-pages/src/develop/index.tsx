// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useParams, useSearchParams } from 'react-router-dom';
import { useEffect, useState } from 'react';

import { useShallow } from 'zustand/react/shallow';
import { useBreadcrumb } from '@cozeloop/base-hooks';

import { usePromptStore } from '@/store/use-prompt-store';
import { usePromptMockDataStore } from '@/store/use-mockdata-store';
import { useBasicStore } from '@/store/use-basic-store';
import { usePrompt } from '@/hooks/use-prompt';
import { PromptDev } from '@/components/prompt-dev';

export function PromptDevelop() {
  const { promptID } = useParams<{
    promptID: string;
  }>();
  const [searchParams] = useSearchParams();
  const queryVersion = searchParams.get('version') || undefined;
  const [getPromptLoading, setGetPromptLoading] = useState(true);

  const { getPromptByVersion } = usePrompt({ promptID, registerSub: true });
  const { clearStore: clearPromptStore, promptInfo } = usePromptStore(
    useShallow(state => ({
      clearStore: state.clearStore,
      promptInfo: state.promptInfo,
    })),
  );

  const { clearStore: clearBasicStore } = useBasicStore(
    useShallow(state => ({ clearStore: state.clearStore })),
  );

  const { clearMockDataStore } = usePromptMockDataStore(
    useShallow(state => ({
      clearMockDataStore: state.clearMockDataStore,
    })),
  );

  useBreadcrumb({
    text: promptInfo?.prompt_basic?.display_name || '',
  });

  useEffect(() => {
    if (promptID) {
      getPromptByVersion(queryVersion, true).then(() => {
        setGetPromptLoading(false);
      });
    }
    return () => {
      clearPromptStore();
      clearBasicStore();
      clearMockDataStore();
    };
  }, [promptID, queryVersion]);

  return <PromptDev getPromptLoading={getPromptLoading} />;
}
