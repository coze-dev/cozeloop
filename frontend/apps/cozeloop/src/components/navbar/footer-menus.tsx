// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type ReactNode, useState } from 'react';

import cls from 'classnames';
import { I18n } from '@cozeloop/i18n-adapter';
import {
  IconCozDocument,
  IconCozArrowDown,
  IconCozArrowUp,
  IconCozLarkFill,
} from '@coze-arch/coze-design/icons';
import { Nav, Collapsible } from '@coze-arch/coze-design';

import {
  COZELOOP_DOC_URL,
  COZELOOP_LARK_GROUP_URL,
  COZELOOP_GITHUB_URL,
  COZELOOP_DISCORD_URL,
} from '@/constants';
import { ReactComponent as IconGithub } from '@/assets/images/github.svg';
import { ReactComponent as IconDiscord } from '@/assets/images/discord.svg';

import { ItemWithLink } from './item-with-link';

interface Props {
  isCollapsed?: boolean;
  isHovered?: boolean;
}

interface MenuItem {
  key: string;
  text: string;
  icon: ReactNode;
  className?: string;
  onClick: () => void;
}

export function FooterMenus({ isCollapsed, isHovered }: Props) {
  const [isShow, setIsShow] = useState(true);
  const menuItems: MenuItem[] = [
    {
      text: I18n.t('document'),
      key: 'actions/doc',
      icon: <IconCozDocument className="coz-fg-secondary" />,
      onClick: () => window.open(COZELOOP_DOC_URL),
    },
    {
      text: I18n.t('lark_group'),
      key: 'actions/lark',
      icon: <IconCozLarkFill className="coz-fg-secondary" />,
      onClick: () => window.open(COZELOOP_LARK_GROUP_URL),
    },
    {
      text: 'GitHub',
      key: 'actions/github',
      icon: <IconGithub className="w-[14px] h-[14px]" />,
      onClick: () => window.open(COZELOOP_GITHUB_URL),
    },
    {
      text: 'Discord',
      key: 'actions/discord',
      icon: <IconDiscord className="w-[14px] h-[14px]" />,
      onClick: () => window.open(COZELOOP_DISCORD_URL),
    },
  ];

  return (
    <>
      {isHovered ? (
        <div
          className="w-12 h-5 rounded-[6px] border border-solid border-[var(--coz-stroke-primary)] absolute left-1/2 -translate-x-1/2 -top-3 bg-white text-center cursor-pointer"
          onClick={() => {
            setIsShow(!isShow);
          }}
        >
          {isShow ? (
            <IconCozArrowDown className="coz-fg-primary" />
          ) : (
            <IconCozArrowUp className="coz-fg-primary" />
          )}
        </div>
      ) : null}

      <Collapsible isOpen={isShow}>
        <div className="flex flex-col w-full gap-1 mb-3">
          {menuItems.map(({ key, text, icon, className, onClick }) => (
            <Nav.Item
              key={key}
              itemKey={key}
              text={<ItemWithLink showLink={!isCollapsed}>{text}</ItemWithLink>}
              icon={icon}
              onClick={onClick}
              className={cls(
                'flex items-center h-[32px] rounded-[6px] hover:coz-mg-primary',
                className,
              )}
            />
          ))}
        </div>
      </Collapsible>
    </>
  );
}
