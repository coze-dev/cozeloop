// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useNavigate } from 'react-router-dom';
import { useState } from 'react';

import { I18n } from '@cozeloop/i18n-adapter';
import { AccountSetting } from '@cozeloop/auth-pages';
import { useLogout, useUserStore, useSpaceStore } from '@cozeloop/account';
import { Popover, Modal } from '@coze-arch/coze-design';

import { LOGIN_PATH } from '@/constants';

import { UserInfo } from './user-info';
import { SettingsMenu } from './settings-menu';

interface Props {
  isCollapsed?: boolean;
}

export function UserInfoSection({ isCollapsed }: Props) {
  const logout = useLogout();
  const navigate = useNavigate();
  const [visible, setVisible] = useState(false);
  const resetUser = useUserStore(s => s.reset);
  const resetSpace = useSpaceStore(s => s.reset);

  const handleAction = (action: 'logout' | 'setting') => {
    switch (action) {
      case 'logout':
        Modal.confirm({
          title: I18n.t('confirm_logout'),
          okText: I18n.t('logout'),
          cancelText: I18n.t('cancel'),
          type: 'modal',
          autoLoading: true,
          okButtonProps: {
            style: {
              backgroundColor: '#FF2710',
              color: '#ffffff',
            },
          },
          onOk: async () => {
            await logout.runAsync();
            resetUser();
            resetSpace();
            navigate(LOGIN_PATH);
          },
        });
        break;
      case 'setting':
        setVisible(true);
        break;
      default:
        break;
    }
  };
  return (
    <>
      <Popover
        position="rightBottom"
        content={<SettingsMenu onAction={handleAction} />}
        trigger="click"
        className="!p-0 rounded-[6px]"
        spacing={26}
        clickToHide={true}
        stopPropagation={true}
      >
        <div>
          <UserInfo isCollapsed={isCollapsed} />
        </div>
      </Popover>
      <Modal
        visible={visible}
        height={600}
        width={1120}
        title={null}
        footer={null}
        keepDOM={false}
        onCancel={() => setVisible(false)}
      >
        <AccountSetting />
      </Modal>
    </>
  );
}
