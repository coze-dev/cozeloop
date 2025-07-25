// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useNavigate, useParams } from 'react-router-dom';

import { useRequest } from 'ahooks';
import { SpaceType } from '@cozeloop/api-schema/foundation';
import { useSpaceStore } from '@cozeloop/account';

export enum SetupSpaceStatus {
  NOT_FOUND = 'not-found',
  FETCH_ERROR = 'fetch-error',
  OK = 'ok',
}

export function useSetupSpace() {
  const fetchSpaces = useSpaceStore(s => s.fetchSpaces);
  const patch = useSpaceStore(s => s.patch);
  const navigate = useNavigate();
  const { enterpriseID, spaceID } = useParams<{
    enterpriseID: string;
    spaceID?: string;
  }>();

  const { data: status, loading } = useRequest(
    async () => {
      try {
        const { spaces } = await fetchSpaces();

        const space = spaceID
          ? spaces?.find(it => it.id === spaceID)
          : spaces?.find(it => it.space_type === SpaceType.Personal) ||
            spaces?.[0];

        if (!space?.id) {
          return SetupSpaceStatus.NOT_FOUND;
        }

        patch({ space });

        if (spaceID !== space.id) {
          const url = `/console${enterpriseID ? `/enterprise/${enterpriseID}` : ''}/space/${space.id}`;

          navigate(url);
        }
        return SetupSpaceStatus.OK;
      } catch (e) {
        console.error(e);
        return SetupSpaceStatus.FETCH_ERROR;
      }
    },
    { refreshDeps: [enterpriseID, spaceID] },
  );

  return { status, loading };
}
