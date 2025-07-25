// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import classNames from 'classnames';

interface PrimaryPageHeaderProps {
  pageTitle?: string;
  filterSlot?: React.ReactNode;
  children?: React.ReactNode;
  contentClassName?: string;
  className?: string;
  titleSlot?: React.ReactNode;
}

export const PrimaryPage = ({
  pageTitle,
  filterSlot,
  children,
  contentClassName,
  className,
  titleSlot,
}: PrimaryPageHeaderProps) => (
  <div
    className={classNames(
      'pt-4 pb-6 h-full max-h-full flex flex-col',
      className,
    )}
  >
    <div className="flex items-center justify-between">
      <div className="text-[20px] font-medium leading-6 coz-fg-plus pt-2 pb-4 px-6">
        {pageTitle}
      </div>
      <div>{titleSlot}</div>
    </div>
    <div className="box-border coz-fg-secondary pt-1 pb-3 px-6">
      {filterSlot}
    </div>
    <div
      className={classNames(
        'flex-1 h-full max-h-full overflow-hidden px-6',
        contentClassName,
      )}
    >
      {children}
    </div>
  </div>
);
