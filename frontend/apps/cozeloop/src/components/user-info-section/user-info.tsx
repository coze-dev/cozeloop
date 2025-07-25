// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import classNames from 'classnames';
import { useUserStore } from '@cozeloop/account';
import { Typography, CozAvatar } from '@coze-arch/coze-design';

interface Props {
  isCollapsed?: boolean;
  active?: boolean;
  onClick?: () => void;
}

export function UserInfo({ isCollapsed, onClick, active }: Props) {
  const userInfo = useUserStore(s => s.userInfo);

  return (
    <div
      className={classNames(
        'flex items-center rounded-[6px] w-full h-[50px] box-border px-1 hover:coz-mg-primary cursor-pointer',
        {
          'coz-mg-primary font-medium': active,
        },
      )}
      onClick={onClick}
    >
      <CozAvatar className="flex-shrink-0" src={userInfo?.avatar_url}>
        {userInfo?.nick_name}
      </CozAvatar>

      {isCollapsed ? null : (
        <div className="flex flex-col w-full overflow-hidden ml-2">
          <Typography.Text
            strong
            ellipsis={{ showTooltip: true }}
            className="!coz-fg-plus text-[13px]"
          >
            {userInfo?.nick_name}
          </Typography.Text>
          <Typography.Text
            size="small"
            className="!coz-fg-secondary"
            ellipsis={{ showTooltip: true }}
          >
            @{userInfo?.name}
          </Typography.Text>
        </div>
      )}
    </div>
  );
}
