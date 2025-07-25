// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useRef, useState } from 'react';

import { useRequest } from 'ahooks';
import { useSpace, useCurrentEnterpriseId } from '@cozeloop/biz-hooks-adapter';
import type { View } from '@cozeloop/api-schema/observation';
import { observabilityTrace } from '@cozeloop/api-schema';

const MAX_VIEW_COUNT = 5;

interface UseCustomViewProps {
  onSuccess?: (viewList: View[]) => void;
}
export const useCustomView = ({ onSuccess }: UseCustomViewProps = {}) => {
  const [viewNames, setViewNames] = useState<string[]>([]);
  const [viewList, setViewList] = useState<View[]>([]);
  const [autoSelectedViewId, setAutoSelectedViewId] = useState<
    number | string | null
  >(null);

  const [activeViewKey, setActiveViewKey] = useState<number | string | null>(
    null,
  );
  const [updateViewFlag, setUpdateViewFlag] = useState(0);
  const { spaceID } = useSpace();
  const enterpriseId = useCurrentEnterpriseId();
  const lastVisibleIds = useRef<(string | number)[] | null>(null);
  const [visibleViewIds, setVisibleViewIds] = useState<(string | number)[]>([]);

  useRequest(
    async () => {
      const result = await observabilityTrace.ListViews({
        enterprise_id: enterpriseId,
        workspace_id: spaceID,
      });

      return (result.views ?? []).map(view => {
        if (!view.is_system) {
          return view;
        }
        return {
          ...view,
          span_list_type: 'root_span',
          platform_type: 'cozeloop',
        };
      });
    },
    {
      refreshDeps: [updateViewFlag],
      onSuccess: viewListData => {
        setViewList(viewListData as View[]);
        setViewNames(viewListData.map(item => item.view_name));
        if (autoSelectedViewId) {
          setActiveViewKey(autoSelectedViewId);
          setAutoSelectedViewId(null);
        }
        onSuccess?.(viewListData as View[]);
        const showViewIds = viewListData.map(item => item.id);
        const needShowIds =
          lastVisibleIds.current !== null
            ? lastVisibleIds.current
            : showViewIds.slice(0, MAX_VIEW_COUNT);
        setVisibleViewIds(needShowIds);
        lastVisibleIds.current = needShowIds;
      },
    },
  );

  const updateViewList = () => {
    setUpdateViewFlag(updateViewFlag + 1);
  };

  return {
    viewNames,
    viewList,
    setAutoSelectedViewId,
    setViewNames,
    setViewList,
    autoSelectedViewId,
    activeViewKey,
    setActiveViewKey,
    updateViewList,
    visibleViewIds,
    setVisibleViewIds,
    lastVisibleIds,
  };
};
