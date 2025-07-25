// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @coze-arch/max-line-per-function */
/* eslint-disable complexity */

import { useNavigate } from 'react-router-dom';
import { useEffect, useState } from 'react';

import { useShallow } from 'zustand/react/shallow';
import { useRequest } from 'ahooks';
import { sendEvent, EVENT_NAMES } from '@cozeloop/tea-adapter';
import { PromptCreate } from '@cozeloop/prompt-components';
import { I18n } from '@cozeloop/i18n-adapter';
import { getBaseUrl } from '@cozeloop/components';
import { useSpace } from '@cozeloop/biz-hooks-adapter';
import { useModalData } from '@cozeloop/base-hooks';
import { type CommitInfo, type Prompt } from '@cozeloop/api-schema/prompt';
import { promptManage } from '@cozeloop/api-schema';
import { IconCozDuplicate, IconCozUpdate } from '@coze-arch/coze-design/icons';
import {
  Button,
  List,
  Modal,
  Space,
  Spin,
  Toast,
} from '@coze-arch/coze-design';

import { sleep } from '@/utils/prompt';
import { usePromptStore } from '@/store/use-prompt-store';
import { useBasicStore } from '@/store/use-basic-store';
import { useVersionList } from '@/hooks/use-version-list';
import { usePrompt } from '@/hooks/use-prompt';
import { CALL_SLEEP_TIME } from '@/consts';

import VersionItem from './version-item';

export function VersionList() {
  const { spaceID } = useSpace();
  const navigate = useNavigate();
  const baseURL = getBaseUrl(spaceID);
  const { promptInfo } = usePromptStore(
    useShallow(state => ({ promptInfo: state.promptInfo })),
  );
  const {
    setVersionChangeLoading,
    setVersionChangeVisible,
    versionChangeLoading,
  } = useBasicStore(
    useShallow(state => ({
      setVersionChangeLoading: state.setVersionChangeLoading,
      setVersionChangeVisible: state.setVersionChangeVisible,
      versionChangeLoading: state.versionChangeLoading,
    })),
  );
  const [draftVersion, setDraftVersion] = useState<CommitInfo>();

  const { getPromptByVersion } = usePrompt({ promptID: promptInfo?.id });

  const [activeVersion, setActiveVersion] = useState<string | undefined>();

  const [getDraftLoading, setGetDraftLoading] = useState(true);

  const promptInfoModal = useModalData<Prompt>();

  const {
    versionListData,
    versionListLoadMore,
    versionListLoading,
    versionListReload,
    versionListLoadingMore,
  } = useVersionList({
    promptID: promptInfo?.id,
    draftVersion,
  });

  const isActionButtonShow = Boolean(activeVersion);

  const { runAsync: rollbackRunAsync } = useRequest(
    () =>
      promptManage.RevertDraftFromCommit({
        prompt_id: promptInfo?.id,
        commit_version_reverting_from: activeVersion,
      }),
    {
      manual: true,
      ready: Boolean(spaceID && promptInfo?.id && activeVersion),
      refreshDeps: [spaceID, promptInfo?.id, activeVersion],
      onSuccess: async () => {
        Toast.success(I18n.t('rollback_success'));
        setVersionChangeLoading(true);
        await sleep(CALL_SLEEP_TIME);
        getPromptByVersion()
          .then(() => {
            setVersionChangeLoading(false);
            setVersionChangeVisible(false);
          })
          .catch(() => {
            setVersionChangeLoading(false);
            setVersionChangeVisible(false);
          });
      },
    },
  );

  const handleVersionChange = (version?: string) => {
    if (version === activeVersion) {
      return;
    }
    setVersionChangeLoading(true);
    getPromptByVersion(version || '', true)
      .then(() => {
        setVersionChangeLoading(false);
        sendEvent(EVENT_NAMES.cozeloop_pe_version, {
          prompt_id: `${promptInfo?.id || 'playground'}`,
        });
      })
      .catch(() => {
        setVersionChangeLoading(false);
      });

    setActiveVersion(version);
  };

  useEffect(() => {
    if (spaceID && promptInfo?.id) {
      promptInfo?.prompt_draft?.draft_info &&
        setDraftVersion({
          version: '',
          base_version:
            promptInfo?.prompt_draft?.draft_info?.base_version || '',
          description: '',
          committed_by: '',
          committed_at: promptInfo?.prompt_draft?.draft_info?.updated_at,
        });
      setActiveVersion('');
      setGetDraftLoading(false);
      setTimeout(() => {
        versionListReload();
      }, CALL_SLEEP_TIME);
    }
    return () => {
      setActiveVersion(undefined);
      setGetDraftLoading(true);
    };
  }, [spaceID, promptInfo?.id]);

  return (
    <div className="flex-1 w-full h-full py-6 flex flex-col gap-2 overflow-hidden ">
      <div
        className="w-full h-full overflow-y-auto px-6"
        onScroll={e => {
          const target = e.currentTarget;

          const isAtBottom =
            target.scrollHeight - target.scrollTop <= target.clientHeight + 1;

          if (
            !versionListData?.hasMore ||
            !isAtBottom ||
            versionListLoadingMore
          ) {
            return;
          }
          versionListLoadMore();
        }}
      >
        <List
          dataSource={versionListData?.list || []}
          renderItem={item => (
            <VersionItem
              className="cursor-pointer mb-3"
              key={item.version}
              active={activeVersion === item.version}
              version={item}
              onClick={() => handleVersionChange(item.version)}
            />
          )}
          size="small"
          emptyContent={
            versionListLoading || getDraftLoading ? <div></div> : null
          }
          loadMore={
            versionListLoadingMore || getDraftLoading ? (
              <div className="w-full text-center">
                <Spin />
              </div>
            ) : null
          }
        />
      </div>

      {isActionButtonShow ? (
        <Space className="w-full flex-shrink-0 px-6">
          <Button
            className="flex-1"
            color="primary"
            disabled={versionChangeLoading}
            icon={<IconCozDuplicate />}
            onClick={() => promptInfoModal.open(promptInfo)}
          >
            {I18n.t('create_copy')}
          </Button>
          <Button
            className="flex-1"
            color="red"
            disabled={versionChangeLoading}
            icon={<IconCozUpdate />}
            onClick={() =>
              Modal.confirm({
                title: I18n.t('restore_to_this_version'),
                content: I18n.t('restore_version_tip'),
                onOk: rollbackRunAsync,
                cancelText: I18n.t('Cancel'),
                okText: I18n.t('restore'),
                okButtonProps: {
                  color: 'red',
                },
                autoLoading: true,
              })
            }
          >
            {I18n.t('restore_to_this_version')}
          </Button>
        </Space>
      ) : null}
      <PromptCreate
        visible={promptInfoModal.visible}
        onCancel={promptInfoModal.close}
        data={promptInfoModal?.data}
        isCopy
        onOk={res => {
          navigate(`${baseURL}/pe/prompts/${res.cloned_prompt_id}`);
          promptInfoModal.close();
        }}
      />
    </div>
  );
}
