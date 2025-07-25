// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useMemo } from 'react';

import classNames from 'classnames';
import { type Message, Role } from '@cozeloop/api-schema/evaluation';

import { splitStringByDoubleBrace } from '../../../utils/double-brace';

const ROLE_MAP: Record<Role | 'Placeholder', string> = {
  [Role.System]: 'System',
  [Role.User]: 'User',
  [Role.Assistant]: 'Assistant',
  [Role.Tool]: 'Tool',
  Placeholder: 'Placeholder',
};

export function PromptMessage({
  className,
  message,
}: {
  className?: string;
  message?: Message;
}) {
  const items = useMemo(
    () => splitStringByDoubleBrace(message?.content?.text || ''),
    [message],
  );

  return (
    <div
      className={classNames(
        'w-full rounded-[6px] border border-solid coz-stroke-primary overflow-hidden',
        className,
      )}
    >
      <div className="rounded-t-[6px] coz-bg-secondary px-3 text-[13px] leading-9 font-['JetBrainsMonoBold'] font-normal coz-fg-dim border-0 border-b border-solid coz-stroke-primary ">
        {message?.role ? ROLE_MAP[message?.role] : null}
      </div>
      <div className="px-3 py-2 coz-bg-primary coz-text-primary text-[13px] leading-5 font-normal break-words whitespace-break-spaces min-h-[36px] max-h-[500px] overflow-auto">
        {items.map((i, idx) => (
          <span
            key={idx}
            className={i.isDoubleBrace ? 'text-[#00A136]' : undefined}
          >
            {i.text}
          </span>
        ))}
      </div>
    </div>
  );
}
