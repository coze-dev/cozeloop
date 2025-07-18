// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import { IconCozIllusAdd } from '@coze-arch/coze-design/illustrations';
import { EmptyState } from '@coze-arch/coze-design';

export function ExperimentListEmptyState({
  hasFilterCondition,
}: {
  hasFilterCondition: boolean;
}) {
  return (
    <EmptyState
      size="full_screen"
      icon={<IconCozIllusAdd />}
      title={
        hasFilterCondition
          ? I18n.t('failed_to_find_related_results')
          : I18n.t('no_experiment')
      }
      description={
        hasFilterCondition
          ? I18n.t('try_other_keywords_or_modify_filter_options')
          : I18n.t('click_to_create_experiment')
      }
    />
  );
}
