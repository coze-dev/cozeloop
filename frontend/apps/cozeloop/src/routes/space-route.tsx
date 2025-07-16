// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { Navigate } from 'react-router-dom';

import { I18n } from '@cozeloop/i18n-adapter';
import { PageNoContent } from '@cozeloop/components';
import { useSpaceStore } from '@cozeloop/account';

interface Props {
  index?: boolean;
}

export function SpaceRoute({ index }: Props) {
  const space = useSpaceStore(s => s.space);

  if (!space?.id) {
    return (
      <PageNoContent
        title={I18n.t('no_space')}
        description={I18n.t('not_join_space')}
      />
    );
  }

  const path = index ? `space/${space.id}` : space.id;

  return <Navigate to={path} replace />;
}
