// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @coze-arch/max-line-per-function */
/* eslint-disable complexity */
import { useMemo, useRef } from 'react';

import { useRequest } from 'ahooks';
import { EVENT_NAMES, sendEvent } from '@cozeloop/tea-adapter';
import { useSpace } from '@cozeloop/biz-hooks-adapter';
import { type Prompt } from '@cozeloop/api-schema/prompt';
import { StonePromptApi } from '@cozeloop/api-schema';
import {
  Form,
  FormInput,
  FormTextArea,
  Modal,
  type FormApi,
} from '@coze-arch/coze-design';
import { I18n } from '@cozeloop/i18n-adapter';

interface PromptCreateProps {
  visible: boolean;
  data?: Prompt;
  isEdit?: boolean;
  isCopy?: boolean;
  onOk: (v: Prompt & { cloned_prompt_id?: Int64 }) => void;
  onCancel: () => void;
}
interface FormValueProps {
  prompt_key?: string;
  prompt_name?: string;
  prompt_description?: string;
  version?: string;
}
export function PromptCreate({
  visible,
  data,
  isCopy,
  isEdit,
  onOk,
  onCancel,
}: PromptCreateProps) {
  const formApi = useRef<FormApi<FormValueProps>>();
  const { spaceID } = useSpace();

  const createService = useRequest(
    (prompt: FormValueProps) =>
      StonePromptApi.CreatePrompt({
        prompt_key: prompt.prompt_key || '',
        prompt_name: prompt.prompt_name || '',
        prompt_description: prompt.prompt_description,
        workspace_id: spaceID,
        draft_detail: data?.prompt_commit?.detail,
      }),
    {
      manual: true,
    },
  );
  const updateService = useRequest(
    (prompt: FormValueProps) =>
      StonePromptApi.UpdatePrompt({
        prompt_id: data?.id || '',
        prompt_name: prompt.prompt_name || '',
        prompt_description: prompt.prompt_description,
      }),
    {
      manual: true,
    },
  );
  const copyService = useRequest(
    (prompt: FormValueProps) =>
      StonePromptApi.ClonePrompt({
        prompt_id: data?.id || '',
        cloned_prompt_key: prompt.prompt_key || '',
        cloned_prompt_name: prompt.prompt_name || '',
        cloned_prompt_description: prompt.prompt_description,
        commit_version: data?.prompt_commit?.commit_info?.version,
      }),
    {
      manual: true,
    },
  );
  const handleOk = async () => {
    const formData = await formApi.current?.validate();
    if (!formData) {
      return;
    }

    if (isCopy) {
      const res = await copyService.runAsync(formData);
      onOk({ ...data, cloned_prompt_id: res.cloned_prompt_id });
      sendEvent(EVENT_NAMES.prompt_create, {
        prompt_id: `${data?.id || ''}`,
        prompt_key: data?.prompt_key || '',
        original_version: formData?.version,
      });
    } else if (isEdit) {
      await updateService.runAsync(formData);
      sendEvent(EVENT_NAMES.prompt_create, {
        prompt_id: `${data?.id || ''}`,
        prompt_key: data?.prompt_key || '',
        is_update: true,
      });
      onOk({
        ...data,
        prompt_basic: {
          ...data?.prompt_basic,
          display_name: formData.prompt_name,
          description: formData.prompt_description,
        },
      });
    } else {
      const res = await createService.runAsync(formData);
      sendEvent(EVENT_NAMES.prompt_create, {
        prompt_id: `${data?.id || ''}`,
      });
      onOk({ ...data, id: res.prompt_id });
    }
  };

  const modalTitle = useMemo(() => {
    if (isEdit) {
      return I18n.t('edit_prompt');
    }
    if (isCopy) {
      return I18n.t('create_copy');
    }
    return I18n.t('create_prompt');
  }, [isCopy, isEdit]);

  return (
    <Modal
      title={modalTitle}
      visible={visible}
      onCancel={onCancel}
      onOk={handleOk}
      cancelText={I18n.t('Cancel')}
      okText={I18n.t('confirm')}
      okButtonProps={{
        loading:
          createService.loading || updateService.loading || copyService.loading,
      }}
      width={600}
    >
      <Form<FormValueProps>
        getFormApi={api => (formApi.current = api)}
        initValues={{
          prompt_key: isCopy
            ? `${
                (data?.prompt_key?.length || 0) < 95
                  ? `${data?.prompt_key}_copy`
                  : data?.prompt_key
              }`
            : data?.prompt_key,
          prompt_name: isCopy
            ? `${
                (data?.prompt_basic?.display_name?.length || 0) < 95
                  ? `${data?.prompt_basic?.display_name}_copy`
                  : data?.prompt_basic?.display_name
              }`
            : data?.prompt_basic?.display_name,
          prompt_description: data?.prompt_basic?.description,
        }}
      >
        <FormInput
          label="Prompt Key"
          field="prompt_key"
          placeholder={I18n.t('please_input', { field: I18n.t('prompt_key') })}
          rules={[
            {
              required: true,
              message: I18n.t('field_is_required', {
                field: I18n.t('prompt_key'),
              }),
            },
            {
              validator: (_rule, value) => {
                if (value && !/^[a-zA-Z][a-zA-Z0-9_.]*$/.test(value)) {
                  return new Error(I18n.t('prompt_key_format'));
                }
                return true;
              },
            },
          ]}
          maxLength={100}
          max={100}
          disabled={isEdit}
        />
        <FormInput
          label={I18n.t('prompt_name')}
          field="prompt_name"
          placeholder={I18n.t('please_input', { field: I18n.t('prompt_name') })}
          rules={[
            {
              required: true,
              message: I18n.t('field_is_required', {
                field: I18n.t('prompt_name'),
              }),
            },
            {
              validator: (_rule, value) => {
                if (value && !/^[\u4e00-\u9fa5a-zA-Z0-9_.-]+$/.test(value)) {
                  return new Error(I18n.t('prompt_name_format'));
                }
                if (value && /^[_.-]/.test(value)) {
                  return new Error(I18n.t('prompt_name_format'));
                }
                return true;
              },
            },
          ]}
          maxLength={100}
          max={100}
        />
        <FormTextArea
          label={I18n.t('prompt_description')}
          field="prompt_description"
          placeholder={I18n.t('please_input', {
            field: I18n.t('prompt_description'),
          })}
          maxCount={500}
          maxLength={500}
        />
      </Form>
    </Modal>
  );
}
