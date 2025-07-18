// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

import { I18n } from '@cozeloop/i18n-adapter';

export const evaluateTargetValidators = {
  evalTargetType: [
    {
      required: true,
      message: I18n.t('please_select', { field: I18n.t('type') }),
    },
  ],
  evalTarget: [
    {
      required: true,
      message: I18n.t('please_select', { field: I18n.t('evaluation_object') }),
    },
  ],
  evalTargetVersion: [
    {
      required: true,
      message: I18n.t('please_select', {
        field: I18n.t('evaluation_object_version'),
      }),
    },
  ],
  // todo: 这里注册进来
  evalTargetMapping: [
    { required: true, message: I18n.t('config_evaluation_object_mapping') },
  ],
};
