// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import { type ModelConfig } from '@cozeloop/api-schema/evaluation';

import { useGlobalEvalConfig } from '@/stores/eval-global-config';

export function ModelConfigInfo({ data }: { data?: ModelConfig }) {
  const { modelConfigEditor: ModelConfigEditor } = useGlobalEvalConfig();

  return (
    <>
      <div className="text-sm font-medium coz-fg-primary mb-2">
        {I18n.t('model')}
      </div>
      {ModelConfigEditor && data ? (
        <ModelConfigEditor
          value={data}
          disabled={true}
          popoverProps={{ position: 'bottomRight' }}
        />
      ) : (
        '-'
      )}
    </>
  );
}
