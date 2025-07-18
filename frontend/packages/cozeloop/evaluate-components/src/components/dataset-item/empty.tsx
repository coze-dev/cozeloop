// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import { IconCozEmpty } from '@coze-arch/coze-design/icons';

export const EmptyDatasetItem = () => (
  <div className="flex flex-col items-center justify-center h-[80px] p-3 coz-fg-dim ">
    <IconCozEmpty className="h-[32px] w-[32px]" />
    <span className="text-[12px]">{I18n.t('no_data')}</span>
  </div>
);
