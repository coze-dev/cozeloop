// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useState } from 'react';

import { Loading } from '@coze-arch/coze-design';

export const useEditorLoading = () => {
  const [loading, setLoading] = useState(true);
  const onMount = () => {
    setLoading(false);
  };
  const LoadingNode = loading && (
    <div className="absolute bg-[white] z-20  top-0 bottom-0 left-0 right-0 flex items-center justify-center">
      <Loading loading={true} />
    </div>
  );
  return { LoadingNode, loading, onMount };
};
