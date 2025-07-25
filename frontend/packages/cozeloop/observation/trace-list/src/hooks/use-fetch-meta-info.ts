// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useRequest } from 'ahooks';
import { I18n } from '@cozeloop/i18n-adapter';
import { useSpace } from '@cozeloop/biz-hooks-adapter';
import {
  type PlatformType,
  type SpanListType,
} from '@cozeloop/api-schema/observation';
import { observabilityTrace } from '@cozeloop/api-schema';
import { Toast } from '@coze-arch/coze-design';

import { useTraceStore } from '../stores/trace';

export const useFetchMetaInfo = () => {
  const { spaceID } = useSpace();
  const { selectedPlatform, selectedSpanType, setFieldMetas } = useTraceStore();
  useRequest(
    async () => {
      const result = await observabilityTrace.GetTracesMetaInfo(
        {
          platform_type: selectedPlatform as PlatformType,
          span_list_type: selectedSpanType as SpanListType,
          workspace_id: spaceID,
        },
        {
          __disableErrorToast: true,
        },
      );
      setFieldMetas(result?.field_metas ?? {});
    },
    {
      refreshDeps: [selectedPlatform, selectedSpanType],
      onError(e) {
        Toast.error(
          I18n.t('observation_fetch_meta_error', {
            msg: e.message || '',
          }),
        );
      },
    },
  );
};
