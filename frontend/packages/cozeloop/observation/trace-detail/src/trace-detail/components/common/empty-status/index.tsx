// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import {
  IconCozIllusEmpty,
  IconCozIllusError,
} from '@coze-arch/coze-design/illustrations';
import { Empty } from '@coze-arch/coze-design';

export const NodeDetailEmpty = () => (
  <Empty
    className="w-full h-full flex items-center justify-center"
    image={<IconCozIllusEmpty style={{ width: 150, height: 150 }} />}
    title={I18n.t('observation_empty_node_unselected')}
    description={I18n.t('observation_empty_to_select_node')}
  />
);

export const RunTreeEmpty = () => (
  <Empty
    className="w-full h-full flex items-center justify-center"
    image={<IconCozIllusError style={{ width: 150, height: 150 }} />}
    title={I18n.t('observation_empty_run_tree_failure')}
    description={I18n.t('observation_empty_data_abnormal')}
  />
);
