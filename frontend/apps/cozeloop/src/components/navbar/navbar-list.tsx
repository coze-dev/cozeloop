// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import cls from 'classnames';
import { useSpaceStore } from '@cozeloop/account';
import { Divider, Nav, Typography } from '@coze-arch/coze-design';

import { type MenuConfig } from './menu-config';
import { ItemWithLink } from './item-with-link';

import styles from './index.module.less';

interface Props {
  isCollapsed?: boolean;
  selectedKeys: string[];
  menus: MenuConfig[];
  className?: string;
}

export function NavbarList({
  isCollapsed,
  selectedKeys,
  menus,
  className,
}: Props) {
  const space = useSpaceStore(s => s.space);

  return (
    <div className={cls(styles['navbar-list'], 'styled-scrollbar', className)}>
      {menus
        .filter(item => item.visible?.({ space }) ?? true)
        .map(menu => {
          if (menu.hideInNavbar) {
            return null;
          }
          if (menu.items?.length) {
            return (
              <div key={menu.itemKey} className="text-xs mb-[10px]">
                {isCollapsed ? (
                  <Divider margin={18} />
                ) : (
                  <div className="flex items-center h-[32px] min-w-fit">
                    <Typography.Text
                      size="small"
                      type="tertiary"
                      className="pl-3 rounded whitespace-nowrap"
                    >
                      {menu.text}
                    </Typography.Text>
                  </div>
                )}
                <div className="flex flex-col gap-1">
                  {menu.items.map(item => (
                    <Nav.Item
                      key={item.itemKey}
                      itemKey={item.itemKey}
                      text={
                        <ItemWithLink showLink={!!item.link && !isCollapsed}>
                          {item.text}
                        </ItemWithLink>
                      }
                      icon={
                        selectedKeys.includes(item.itemKey)
                          ? item.selectedIcon
                          : item.icon
                      }
                      disabled={item.disabled}
                      className={cls(
                        'flex items-center h-[32px] rounded-[6px]',
                        selectedKeys.includes(item.itemKey)
                          ? 'font-medium coz-fg-plus'
                          : 'coz-fg-primary',
                        item.disabled ? 'coz-fg-dim' : 'hover:coz-mg-primary',
                      )}
                      onClick={() => item.link && window.open(item.link)}
                    />
                  ))}
                </div>
              </div>
            );
          }
          return (
            <div key={menu.itemKey}>
              <Nav.Item
                itemKey={menu.itemKey}
                text={menu.text}
                icon={menu.icon}
                className={
                  'items-center height-[32px] text-coz-mg-primary rounded'
                }
              />
            </div>
          );
        })}
    </div>
  );
}
