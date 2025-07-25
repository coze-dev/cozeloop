// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import classNames from 'classnames';
import { OverflowList, Space, Tag, Tooltip } from '@coze-arch/coze-design';

export function PromptVariablesList({
  className,
  variables,
}: {
  className?: string;
  variables: string[];
}) {
  const renderItem = (item: string, index: number) => (
    <>
      <div
        className="flex flex-row items-center h-9 text-xs coz-fg-plus font-bold"
        key={item}
      >
        {index === 0 ? (
          <div className="text-[13px] coz-fg-dim ml-3 font-['JetBrainsMonoBold'] font-normal">
            {'Variables'}
          </div>
        ) : null}
        <div className="mx-3">{item}</div>
        {index !== variables.length - 1 ? (
          <div className="text-[var(--coz-stroke-primary)]">|</div>
        ) : null}
      </div>
    </>
  );

  const renderOverflow = (items: string[]) =>
    items.length ? (
      <Tooltip
        content={
          <Space wrap spacing={3}>
            {items.map(item => (
              <Tag color="primary">{item}</Tag>
            ))}
          </Space>
        }
      >
        <div className="flex flex-row items-center h-9 text-xs coz-fg-plus font-bold mx-3">
          <div className="mx-3">+{items.length}</div>
        </div>
      </Tooltip>
    ) : null;

  return (
    <div
      className={classNames(
        'coz-bg-primary border border-solid coz-stroke-primary rounded-[6px]',
        className,
      )}
    >
      <OverflowList
        // @ts-expect-error OverflowList类型过于垃圾
        items={variables}
        // @ts-expect-error OverflowList类型过于垃圾
        visibleItemRenderer={renderItem}
        // @ts-expect-error OverflowList类型过于垃圾
        overflowRenderer={renderOverflow}
      />
    </div>
  );
}
