// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useState, type ReactNode } from 'react';

import { IconCozMore } from '@coze-arch/coze-design/icons';
import {
  Dropdown,
  Space,
  type TypographyProps,
  Typography,
} from '@coze-arch/coze-design';

import { TooltipWhenDisabled } from '../tooltip-with-disabled';

export interface TableColAction {
  label: ReactNode;
  icon?: ReactNode;
  disabled?: boolean;
  hide?: boolean;
  type?: TypographyProps['type'];
  onClick?: () => void;
}

interface Props {
  actions: TableColAction[];
  maxCount?: number;
  disabled?: boolean;
}

export function TableColActions({ actions, maxCount = 2, disabled }: Props) {
  const [visible, setVisible] = useState(false);
  const filteredActions = actions.filter(action => !action.hide);
  const firstActions = filteredActions.slice(0, maxCount);
  const moreActions = filteredActions.slice(maxCount);

  return (
    <div
      onClick={e => {
        e.stopPropagation();
      }}
    >
      <Space spacing={12}>
        {firstActions.map((action, index) => (
          <TooltipWhenDisabled
            key={index}
            content={action.label}
            disabled={Boolean(action.icon)}
          >
            <Typography.Text
              size="small"
              className={'!text-[13px]'}
              type={action.type}
              disabled={action.disabled ?? disabled}
              onClick={() => {
                if (!(action.disabled ?? disabled)) {
                  action.onClick?.();
                }
              }}
              link={!action.type}
            >
              {action.icon ? null : action.label}
            </Typography.Text>
          </TooltipWhenDisabled>
        ))}
        {moreActions.length > 0 && (
          <Dropdown
            position="bottomLeft"
            visible={visible}
            trigger="custom"
            onClickOutSide={() => setVisible(false)}
            render={
              <Dropdown.Menu mode="menu">
                {moreActions.map((action, index) => (
                  <Dropdown.Item
                    disabled={action.disabled ?? disabled}
                    key={index}
                    onClick={() => {
                      if (!(action.disabled ?? disabled)) {
                        setVisible(false);
                        action.onClick?.();
                      }
                    }}
                    className="min-w-[90px] !p-0 !pl-2"
                    icon={action.icon}
                    style={{ minWidth: '90px' }}
                  >
                    <Typography.Text
                      type={action.type}
                      size="small"
                      className="!text-[13px]"
                      link={!action.type}
                    >
                      {action.label}
                    </Typography.Text>
                  </Dropdown.Item>
                ))}
              </Dropdown.Menu>
            }
          >
            <div
              className="flex items-center justify-center"
              onClick={() => setVisible(true)}
            >
              <IconCozMore className="text-[#5A4DED]" />
            </div>
          </Dropdown>
        )}
      </Space>
    </div>
  );
}
