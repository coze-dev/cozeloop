// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import classNames from 'classnames';
import { I18n } from '@cozeloop/i18n-adapter';
import { IconCozRefresh } from '@coze-arch/coze-design/icons';
import {
  Button,
  Space,
  Tooltip,
  type ButtonProps,
} from '@coze-arch/coze-design';

import { ColumnSelector, type ColumnSelectorProps } from '@/columns-select';

import styles from './index.module.less';

export interface TableHeaderProps {
  columnSelectorConfigProps?: ColumnSelectorProps;
  filterForm?: React.ReactNode;
  actions?: React.ReactNode;
  className?: string;
  style?: React.CSSProperties;
  spaceProps?: Record<string, unknown>;
  refreshButtonPros?: ButtonProps;
}

export function TableHeader({
  columnSelectorConfigProps,
  filterForm,
  actions,
  className,
  style,
  spaceProps,
  refreshButtonPros,
}: TableHeaderProps) {
  return (
    <div
      className={classNames('flex flex-row justify-between w-full ', className)}
      style={style}
    >
      <div className={classNames('flex flex-row', styles['table-header-form'])}>
        {filterForm}
      </div>
      <Space {...(spaceProps || {})}>
        {refreshButtonPros ? (
          <Tooltip content={I18n.t('refresh')} theme="dark">
            <Button
              color="primary"
              icon={<IconCozRefresh />}
              {...refreshButtonPros}
            ></Button>
          </Tooltip>
        ) : null}
        {columnSelectorConfigProps ? (
          <ColumnSelector
            className={classNames(columnSelectorConfigProps.className)}
            {...columnSelectorConfigProps}
          />
        ) : null}
        {actions}
      </Space>
    </div>
  );
}
