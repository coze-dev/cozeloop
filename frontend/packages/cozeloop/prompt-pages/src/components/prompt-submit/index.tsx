// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @coze-arch/max-line-per-function */
/* eslint-disable complexity */

import { useNavigate } from 'react-router-dom';
import { useEffect, useRef, useState } from 'react';

import { useShallow } from 'zustand/react/shallow';
import classNames from 'classnames';
import { useRequest } from 'ahooks';
import { EVENT_NAMES, sendEvent } from '@cozeloop/tea-adapter';
import { I18n } from '@cozeloop/i18n-adapter';
import { getBaseUrl } from '@cozeloop/components';
import { useSpace } from '@cozeloop/biz-hooks-adapter';
import { type Prompt } from '@cozeloop/api-schema/prompt';
import { promptManage } from '@cozeloop/api-schema';
import {
  IconCozIllusDone,
  IconCozIllusDoneDark,
} from '@coze-arch/coze-design/illustrations';
import { IconCozCheckMarkFill } from '@coze-arch/coze-design/icons';
import {
  EmptyState,
  Form,
  type FormApi,
  Loading,
  Modal,
  Skeleton,
  Typography,
} from '@coze-arch/coze-design';

import { sleep, versionValidate } from '@/utils/prompt';
import { usePromptStore } from '@/store/use-prompt-store';
import { CALL_SLEEP_TIME } from '@/consts';

import { DiffContent } from './diff-content';

import styles from './index.module.less';

interface PromptSubmitProps {
  visible: boolean;
  initVersion?: string;
  onCancel?: () => void;
  onOk?: (version: { version?: string }) => void;
}

export function PromptSubmit({
  visible,
  onOk,
  onCancel,
  initVersion,
}: PromptSubmitProps) {
  const formApi = useRef<FormApi<{ version?: string; description?: string }>>();
  const { spaceID } = useSpace();
  const { promptInfo } = usePromptStore(
    useShallow(state => ({ promptInfo: state.promptInfo })),
  );
  const navigate = useNavigate();
  const baseURL = getBaseUrl(spaceID);

  const [basePrompt, setBasePrompt] = useState<Prompt>();
  const [currentPrompt, setCurrentPrompt] = useState<Prompt>();

  const [okButtonText, setOkButtonText] = useState('继续');

  const { runAsync: getPromptByVersion } = useRequest(
    (version?: string) =>
      promptManage.GetPrompt({
        prompt_id: promptInfo?.id ?? '',
        with_draft: !version,
        with_commit: Boolean(version),
      }),
    {
      manual: true,
      ready: Boolean(spaceID && visible),
    },
  );

  const showSuccessModal = () => {
    const modal = Modal.info({
      title: I18n.t('submit_new_version'),
      width: 960,
      closable: true,
      content: (
        <div className="w-full h-[470px] flex items-center justify-center">
          <EmptyState
            icon={<IconCozIllusDone width="160" height="160" />}
            darkModeIcon={<IconCozIllusDoneDark width="160" height="160" />}
            title={
              <Typography.Title heading={5} className="!my-4">
                {I18n.t('version_submit_success')}
              </Typography.Title>
            }
            description={
              <div className="flex flex-col items-center gap-2 w-[400px]">
                <Typography.Text className="flex gap-2 items-center">
                  {I18n.t('cozeloop_sdk_data_report_observation')}
                  <Typography.Text
                    link
                    onClick={() => {
                      navigate(`${baseURL}/observation/traces`);
                      modal.destroy();
                    }}
                  >
                    {I18n.t('go_immediately')}
                  </Typography.Text>
                </Typography.Text>
                <Typography.Text className="flex gap-2 items-center">
                  {I18n.t('prompt_effect_evaluation')}
                  <Typography.Text
                    link
                    onClick={() => {
                      navigate(`${baseURL}/evaluation/datasets`);
                      modal.destroy();
                    }}
                  >
                    {I18n.t('go_immediately')}
                  </Typography.Text>
                </Typography.Text>
              </div>
            }
          />
        </div>
      ),
      okText: I18n.t('close'),
    });
  };

  const { loading: submitLoading, runAsync: submitRunAsync } = useRequest(
    async () => {
      const values = await formApi.current
        ?.validate()
        ?.catch(e => console.error(e));
      if (!values) {
        return;
      }

      try {
        await promptManage.CommitDraft({
          prompt_id: promptInfo?.id || '',
          commit_version: values?.version || '',
          commit_description: values?.description,
        });
        sendEvent(EVENT_NAMES.prompt_submit_info, {
          prompt_id: `${promptInfo?.id || 'playground'}`,
          prompt_key: promptInfo?.prompt_key || 'playground',
          version: values?.version || '',
        });
        await sleep(CALL_SLEEP_TIME);

        await onOk?.({ version: values?.version });

        showSuccessModal();
      } catch (e) {
        console.error(e);
      }
    },
    {
      manual: true,
      ready: Boolean(spaceID && promptInfo?.id),
      refreshDeps: [spaceID, promptInfo?.id],
    },
  );

  useEffect(() => {
    if (visible && promptInfo?.prompt_draft?.draft_info?.base_version) {
      getPromptByVersion().then(vres => {
        setCurrentPrompt(vres.prompt);
        getPromptByVersion(
          promptInfo?.prompt_draft?.draft_info?.base_version,
        ).then(res => {
          setBasePrompt(res?.prompt);
        });
      });
    } else {
      setOkButtonText(I18n.t('continue'));
      setBasePrompt(undefined);
      setCurrentPrompt(undefined);
      formApi.current?.reset();
    }
  }, [visible, promptInfo]);

  const submitForm = (
    <Form
      initValues={{ version: initVersion }}
      getFormApi={api => (formApi.current = api)}
    >
      <Form.Input
        label={{
          text: I18n.t('version'),
          required: true,
        }}
        field="version"
        required
        validate={val => versionValidate(val, initVersion)}
        placeholder={I18n.t('input_version_number')}
      />
      <Form.TextArea
        label={I18n.t('version_description')}
        field="description"
        placeholder={I18n.t('please_input', {
          field: I18n.t('version_description'),
        })}
        maxCount={200}
        maxLength={200}
        rows={5}
      />
    </Form>
  );

  const handleSubmit = () => {
    if (okButtonText === I18n.t('continue')) {
      setOkButtonText(I18n.t('submit'));
    } else {
      submitRunAsync();
    }
  };

  return (
    <Modal
      className="min-h-[calc(100vh - 140px)]"
      width={900}
      visible={visible}
      title={I18n.t('submit_new_version')}
      onCancel={onCancel}
      okText={basePrompt ? okButtonText : I18n.t('submit')}
      cancelText={I18n.t('cancel')}
      onOk={basePrompt ? handleSubmit : submitRunAsync}
      okButtonProps={{ loading: submitLoading }}
      height="fit-content"
    >
      <Skeleton
        loading={Boolean(
          !currentPrompt && promptInfo?.prompt_draft?.draft_info?.base_version,
        )}
        placeholder={
          <div className="w-full flex items-center justify-center  h-[470px]">
            <Loading loading />
          </div>
        }
      >
        <div className="w-full overflow-y-auto">
          {basePrompt ? (
            <div className="flex flex-col gap-2">
              <div className={styles['tab-header']}>
                <Typography.Text
                  className={classNames(
                    styles['tab-step'],
                    styles['tab-active'],
                    'cursor-pointer',
                  )}
                  icon={
                    okButtonText === I18n.t('submit') ? (
                      <span className={styles['tab-icon']}>
                        <IconCozCheckMarkFill />
                      </span>
                    ) : (
                      <span className={styles['tab-icon']}>1</span>
                    )
                  }
                  onClick={() => setOkButtonText(I18n.t('continue'))}
                >
                  {I18n.t('confirm_version_difference')}
                </Typography.Text>
                <Typography.Text
                  className={classNames(styles['tab-step'], {
                    [styles['tab-active']]: okButtonText === I18n.t('submit'),
                  })}
                  icon={<span className={styles['tab-icon']}>2</span>}
                >
                  {I18n.t('confirm_version_info')}
                </Typography.Text>
              </div>
              <div className="flex-1">
                {okButtonText === I18n.t('continue') ? (
                  <DiffContent base={basePrompt} current={currentPrompt} />
                ) : (
                  submitForm
                )}
              </div>
            </div>
          ) : (
            submitForm
          )}
        </div>
      </Skeleton>
    </Modal>
  );
}
