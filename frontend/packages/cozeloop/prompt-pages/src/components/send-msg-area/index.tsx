// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @coze-arch/max-line-per-function */
/* eslint-disable max-lines-per-function */
/* eslint-disable @typescript-eslint/no-magic-numbers */
/* eslint-disable complexity */
import { useEffect, useMemo, useRef, useState } from 'react';

import { useShallow } from 'zustand/react/shallow';
import { nanoid } from 'nanoid';
import cn from 'classnames';
import { useLatest } from 'ahooks';
import {
  EditorView,
  type Extension,
  keymap,
  Prec,
  PromptBasicEditor,
} from '@cozeloop/prompt-components';
import { I18n } from '@cozeloop/i18n-adapter';
import { useSpace } from '@cozeloop/biz-hooks-adapter';
import { uploadFile } from '@cozeloop/biz-components-adapter';
import {
  ContentType,
  type Message,
  Role,
  VariableType,
} from '@cozeloop/api-schema/prompt';
import {
  IconCozBroom,
  IconCozCrossCircleFillPalette,
  IconCozImage,
  IconCozInfoCircle,
  IconCozPlayCircle,
  IconCozStopCircle,
} from '@coze-arch/coze-design/icons';
import {
  Image,
  Badge,
  Button,
  IconButton,
  ImagePreview,
  Space,
  Spin,
  Toast,
  Tooltip,
  Typography,
  Upload,
  type UploadProps,
} from '@coze-arch/coze-design';

import { usePromptStore } from '@/store/use-prompt-store';
import {
  type ContentPartLoop,
  usePromptMockDataStore,
} from '@/store/use-mockdata-store';
import { MAX_FILE_SIZE, MAX_FILE_SIZE_MB, MAX_IMAGE_FILE } from '@/consts';

import styles from './index.module.less';

interface SendMsgAreaProps {
  streaming?: boolean;
  onMessageSend?: (queryMsg?: Message) => void;
  stopStreaming?: () => void;
}

export function SendMsgArea({
  streaming,
  onMessageSend,
  stopStreaming,
}: SendMsgAreaProps) {
  const { spaceID } = useSpace();
  const [editorActive, setEditorActive] = useState(false);

  const [queryMsg, setQueryMsg] = useState<Message>({
    role: Role.User,
  });
  const [queryMsgKey, setQueryMsgKey] = useState<string>(nanoid());

  const { variables, currentModel } = usePromptStore(
    useShallow(state => ({
      variables: state.variables,
      messageList: state.messageList,
      currentModel: state.currentModel,
    })),
  );

  const {
    setHistoricMessage,
    historicMessage,
    compareConfig,
    setHistoricMessageById,
  } = usePromptMockDataStore(
    useShallow(state => ({
      setHistoricMessage: state.setHistoricMessage,
      userDebugConfig: state.userDebugConfig,
      historicMessage: state.historicMessage,
      setHistoricMessageById: state.setHistoricMessageById,
      compareConfig: state.compareConfig,
    })),
  );

  const isCompare = Boolean(compareConfig?.groups?.length);

  const uploadRef = useRef<Upload>(null);

  const historicImgParts =
    historicMessage
      ?.map(
        it =>
          it?.parts?.filter(item => item.type === ContentType.ImageURL) || [],
      )
      ?.flat() || [];

  const imgParts: ContentPartLoop[] =
    queryMsg?.parts?.filter(it => it.type === ContentType.ImageURL) || [];

  const imgCount = historicImgParts.length + imgParts.length;
  const canUploadFileSize = MAX_IMAGE_FILE - historicImgParts.length;

  const isMaxImgSize = Boolean(imgCount >= MAX_IMAGE_FILE);
  const isMaxImgSizeRef = useLatest(isMaxImgSize);

  const fileUploading = imgParts.some(it => it.status === 'uploading');

  const inputReadonly = streaming;

  const executeDisabled = streaming || fileUploading || !currentModel?.model_id;

  const isMultiModal = currentModel?.ability?.multi_modal;
  const isMultiModalRef = useLatest(isMultiModal);

  const removePart = (part: ContentPartLoop) => {
    setQueryMsg(v => ({
      ...v,
      parts: (v?.parts || []).filter(
        (it: ContentPartLoop) => it.uid !== part.uid,
      ),
    }));
  };

  const handleUploadFile: UploadProps['customRequest'] = async ({
    fileInstance,
    file,
    onProgress,
    onSuccess,
    onError,
  }) => {
    const { uid } = file;
    const blobUrl = URL.createObjectURL(file.fileInstance as Blob);
    // 读取文件并转换为Base64
    setQueryMsg(v => ({
      ...v,
      parts: [
        ...(v?.parts || []),
        {
          type: ContentType.ImageURL,
          image_url: { url: blobUrl },
          status: file.status,
          uid,
        },
      ],
    }));

    try {
      const res = await uploadFile({
        file: fileInstance,
        fileType: fileInstance.type?.includes('image') ? 'image' : 'object',
        onProgress,
        onSuccess,
        onError,
        spaceID,
      });

      setQueryMsg(v => ({
        ...v,
        parts: (v?.parts || []).map((it: ContentPartLoop) => {
          if (it.uid === uid) {
            return {
              type: ContentType.ImageURL,
              image_url: { ...it.image_url, uri: res },
              status: 'success',
              uid,
            };
          }
          return it;
        }),
      }));
    } catch (error) {
      console.info('error', error);
      Toast.error(I18n.t('image_upload_error'));
      setQueryMsg(v => ({
        ...v,
        parts: (v?.parts || []).filter((it: ContentPartLoop) => it.uid !== uid),
      }));
    }
  };

  const handleSendMessage = () => {
    if (executeDisabled) {
      return;
    }

    onMessageSend?.(queryMsg);
    setQueryMsg({ role: Role.User });
    setQueryMsgKey(nanoid());
  };
  const handleSendMessageRef = useLatest(handleSendMessage);

  const handleUploadImgByEditor = (items?: DataTransferItemList) => {
    if (items?.length && isMultiModalRef.current) {
      for (const item of Array.from(items)) {
        if (item.type.includes('image')) {
          if (isMaxImgSizeRef.current) {
            Toast.warning(
              I18n.t('max_upload_picture_num', { num: MAX_IMAGE_FILE }),
            );
            return;
          }
          const file = item.getAsFile();
          if (file) {
            if (file.size / 1024 > MAX_FILE_SIZE) {
              Toast.error(
                I18n.t('image_size_not_exceed_num_mb', {
                  num: MAX_FILE_SIZE_MB,
                }),
              );
              return;
            }
            uploadRef.current?.insert([file], 0);
            uploadRef.current?.upload();
          }
        }
      }
    }
  };

  const clearHistoricChat = () => {
    setHistoricMessage([]);
    compareConfig?.groups?.forEach((_, idx) => setHistoricMessageById(idx, []));
  };

  const extensions: Extension[] = useMemo(
    () => [
      EditorView.theme({
        '.cm-gutters': {
          backgroundColor: 'transparent',
          borderRight: 'none',
        },
        '.cm-scroller': {
          paddingLeft: '10px',
          paddingRight: '6px !important',
        },
      }),
      Prec.high(
        keymap.of([
          {
            key: 'Enter',
            run: () => {
              handleSendMessageRef?.current();
              return true;
            },
          },
        ]),
      ),
      EditorView.domEventObservers({
        drop(event) {
          const items = event?.dataTransfer?.items;
          const hasImg = Array.from(items || []).some(it =>
            it.type.includes('image'),
          );
          if (hasImg) {
            event.preventDefault();
          }
          handleUploadImgByEditor(items);
          return true;
        },
        paste(event) {
          const items = event.clipboardData?.items;
          handleUploadImgByEditor(items);
          return true;
        },
      }),
    ],
    [],
  );

  useEffect(() => {
    if (!isMultiModal) {
      setQueryMsg(prev => ({
        ...prev,
        parts: [],
      }));
    }
  }, [isMultiModal]);

  return (
    <div className={styles['send-msg-area']}>
      <div className="flex items-center justify-end">
        <div className="flex-1 flex items-center justify-center">
          {streaming && stopStreaming ? (
            <Space align="center">
              <Button
                color="primary"
                icon={<IconCozStopCircle />}
                size="mini"
                onClick={stopStreaming}
              >
                {I18n.t('stop_respond')}
              </Button>
            </Space>
          ) : null}
        </div>
        {isCompare ? null : (
          <Tooltip content={I18n.t('clear_history_messages')} theme="dark">
            <IconButton
              icon={<IconCozBroom />}
              onClick={clearHistoricChat}
              color="secondary"
              disabled={streaming}
            />
          </Tooltip>
        )}
      </div>
      <div
        className={cn(styles['send-msg-area-content'], {
          [styles['editor-active']]: editorActive,
        })}
      >
        {imgParts?.length ? (
          <ImagePreview closable className={styles['msg-files']}>
            {imgParts?.map((it: ContentPartLoop) => (
              <Badge
                className={styles['msg-files-badge']}
                count={
                  <IconCozCrossCircleFillPalette
                    className={styles['msg-files-badge-icon']}
                    onClick={() => removePart(it)}
                  />
                }
                key={it.image_url?.url || it.image_url?.uri || it.uid}
              >
                <Spin
                  style={{ width: 45, height: 45 }}
                  spinning={it.status === 'uploading'}
                  size="small"
                >
                  <Image
                    width={45}
                    height={45}
                    src={it.image_url?.url}
                    imgStyle={{ objectFit: 'contain' }}
                  />
                </Spin>
              </Badge>
            ))}
          </ImagePreview>
        ) : null}
        <div className={cn('w-full flex-1 gap-0.5')}>
          <PromptBasicEditor
            key={queryMsgKey}
            defaultValue={queryMsg?.content}
            onChange={value =>
              setQueryMsg(v => ({
                ...v,
                content: value,
              }))
            }
            height={44}
            variables={variables?.filter(it => it.type === VariableType.String)}
            readOnly={streaming || inputReadonly}
            linePlaceholder={I18n.t('input_question_tip')}
            customExtensions={extensions}
            onFocus={() => setEditorActive(true)}
            onBlur={() => setEditorActive(false)}
          />
        </div>
        <div className="flex items-center justify-between w-full gap-0.5 px-3">
          <div className="flex items-center gap-2">
            {isCompare ? (
              <Tooltip content={I18n.t('clear_history_messages')} theme="dark">
                <IconButton
                  icon={<IconCozBroom />}
                  onClick={clearHistoricChat}
                  color="secondary"
                  disabled={streaming}
                />
              </Tooltip>
            ) : null}
            <Upload
              key={queryMsgKey}
              ref={uploadRef}
              action=""
              customRequest={handleUploadFile}
              accept="image/*"
              showUploadList={false}
              maxSize={MAX_FILE_SIZE}
              limit={canUploadFileSize}
              onSizeError={() =>
                Toast.error(
                  Toast.error(
                    I18n.t('image_size_not_exceed_num_mb', {
                      num: MAX_FILE_SIZE_MB,
                    }),
                  ),
                )
              }
              onExceed={() =>
                Toast.warning(
                  I18n.t('max_upload_picture_num', { num: MAX_IMAGE_FILE }),
                )
              }
              multiple
              fileList={imgParts.map(it => ({
                uid: it.uid || '',
                url: it.image_url?.url,
                status: it.status || 'success',
                name: it.uid || '',
                size: '0',
              }))}
            >
              <IconButton
                icon={<IconCozImage />}
                color="primary"
                disabled={streaming || isMaxImgSize || !isMultiModal}
              />
            </Upload>
            {isMultiModal ? (
              <Typography.Text size="small" type="tertiary">
                {imgCount} / 20
              </Typography.Text>
            ) : (
              <Typography.Text
                size="small"
                type="tertiary"
                icon={<IconCozInfoCircle />}
              >
                {I18n.t('model_not_support_picture')}
              </Typography.Text>
            )}
          </div>
          <Button
            icon={<IconCozPlayCircle />}
            onClick={handleSendMessage}
            disabled={executeDisabled}
          >
            {I18n.t('run')}
          </Button>
        </div>
      </div>
      <Typography.Text size="small" type="tertiary" className="text-center">
        {I18n.t('generated_by_ai_tip')}
      </Typography.Text>
    </div>
  );
}
