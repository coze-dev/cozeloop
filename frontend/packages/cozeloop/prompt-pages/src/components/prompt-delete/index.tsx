// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useEffect, useState } from 'react';

import { useRequest } from 'ahooks';
import { I18n } from '@cozeloop/i18n-adapter';
import { TextWithCopy } from '@cozeloop/components';
import { type Prompt } from '@cozeloop/api-schema/prompt';
import { StonePromptApi } from '@cozeloop/api-schema';
import { Input, Modal, Space, Toast, Typography } from '@coze-arch/coze-design';

interface PromptDeleteProps {
  data?: Prompt;
  visible: boolean;
  onCacnel?: () => void;
  onOk?: () => void;
}
export function PromptDelete({
  data,
  visible,
  onCacnel,
  onOk,
}: PromptDeleteProps) {
  const [deleteKey, setDeleteKey] = useState('');
  const { runAsync: runDeleteAsync, loading: deleteLoading } = useRequest(
    promptId => StonePromptApi.DeletePrompt({ prompt_id: promptId }),
    {
      manual: true,
      onSuccess: () => {
        Toast.success({
          content: I18n.t('delete_success'),
          showClose: false,
        });
        onOk?.();
      },
    },
  );

  const handleOk = () => {
    if (data?.id && deleteKey === data?.prompt_key) {
      runDeleteAsync(data?.id);
    }
  };

  useEffect(() => {
    if (visible) {
      setDeleteKey('');
    }
  }, [visible]);
  return (
    <Modal
      title={I18n.t('delete_prompt')}
      visible={visible}
      onCancel={onCacnel}
      onOk={handleOk}
      okButtonProps={{
        disabled: Boolean(!deleteKey || deleteKey !== data?.prompt_key),
        loading: deleteLoading,
      }}
      okText={I18n.t('confirm')}
      cancelText={I18n.t('Cancel')}
    >
      <Space vertical style={{ width: '100%' }} align="start">
        <Typography.Text>
          {I18n.t('input_prompt_key_to_delete')}
          <TextWithCopy
            content={data?.prompt_key}
            maxWidth={400}
            className="gap-2"
            copyTooltipText={I18n.t('copy_prompt_key')}
          />
        </Typography.Text>
        <Input
          style={{ width: '100%' }}
          placeholder={I18n.t('prompt_key_again_confirm')}
          value={deleteKey}
          onChange={setDeleteKey}
        />
      </Space>
    </Modal>
  );
}
