// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable complexity */
import cs from 'classnames';
import { formatTimestampToString } from '@cozeloop/toolkit';
import { I18n } from '@cozeloop/i18n-adapter';
import { UserProfile } from '@cozeloop/components';
import { type CommitInfo } from '@cozeloop/api-schema/prompt';
import { type UserInfoDetail } from '@cozeloop/api-schema/foundation';
import { Descriptions, Tag, Typography } from '@coze-arch/coze-design';

import styles from './index.module.less';

export default function VersionItem({
  version,
  active,
  className,
  onClick,
}: {
  version?: CommitInfo & { user?: UserInfoDetail };
  active?: boolean;
  className?: string;
  onClick?: () => void;
}) {
  const isDraft = !version?.version;
  return (
    <div className={`group flex cursor-pointer ${className}`} onClick={onClick}>
      <div className="w-6 h-10 flex items-center shrink-0">
        <div
          className={`w-2 h-2 rounded-full ${active ? 'bg-green-700' : 'bg-gray-300'} `}
        />
      </div>
      <div
        className={`grow px-2 pt-2 rounded-m ${active ? 'bg-gray-100' : ''} group-hover:bg-gray-100`}
      >
        <Descriptions
          align="left"
          className={cs(styles.description, className)}
        >
          <Tag color={isDraft ? 'primary' : 'green'} className="mb-2">
            {isDraft ? I18n.t('current_draft') : I18n.t('submit')}
          </Tag>
          {isDraft ? null : (
            <Descriptions.Item itemKey={I18n.t('version')}>
              <span className="font-medium">{version.version ?? '-'}</span>
            </Descriptions.Item>
          )}
          {!version?.committed_at ? null : (
            <Descriptions.Item
              itemKey={isDraft ? I18n.t('save_time') : I18n.t('submit_time')}
              className="!text-[13px]"
            >
              <span className="font-medium !text-[13px]">
                {version?.committed_at
                  ? formatTimestampToString(
                      version?.committed_at,
                      'YYYY-MM-DD HH:mm:ss',
                    )
                  : '-'}
              </span>
            </Descriptions.Item>
          )}
          {isDraft && !version?.committed_by ? null : (
            <Descriptions.Item
              itemKey={I18n.t('submitter')}
              className="!text-[13px]"
            >
              <UserProfile
                avatarUrl={version?.user?.avatar_url}
                name={version?.user?.nick_name}
              />
            </Descriptions.Item>
          )}
          {isDraft ? null : (
            <Descriptions.Item
              itemKey={I18n.t('version_description')}
              className="!text-[13px]"
            >
              <Typography.Text
                ellipsis={{ rows: 2, showTooltip: true }}
                className="!text-[13px]"
              >
                {version.description || '-'}
              </Typography.Text>
            </Descriptions.Item>
          )}
        </Descriptions>
      </div>
    </div>
  );
}
