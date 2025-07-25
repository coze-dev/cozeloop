// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import { Modal, Radio, RadioGroup } from '@coze-arch/coze-design';

interface OverWriteFieldProps {
  value: boolean;
  onChange?: (value: boolean) => void;
}
export const OverWriteField = ({ value, onChange }: OverWriteFieldProps) => (
  <RadioGroup
    value={value ? 'true' : 'false'}
    onChange={e => {
      const newValue = e.target.value === 'true';
      if (newValue) {
        Modal.confirm({
          title: I18n.t('confirm_select_full_coverage'),
          content: I18n.t('continue_will_override_existing_data'),
          okText: I18n.t('confirm'),
          cancelText: I18n.t('Cancel'),
          onOk: () => {
            onChange?.(true);
          },
          onCancel: () => {
            onChange?.(false);
          },
          okButtonProps: {
            color: 'yellow',
          },
        });
      } else {
        onChange?.(newValue);
      }
    }}
  >
    <Radio value={'false'}>{I18n.t('append_data')}</Radio>
    <Radio value={'true'}>{I18n.t('overwrite_data')}</Radio>
  </RadioGroup>
);
