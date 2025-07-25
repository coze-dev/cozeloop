// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type PropsWithChildren } from 'react';

import { IconCozLongArrowTopRight } from '@coze-arch/coze-design/icons';

export function ItemWithLink({
  showLink,
  children,
}: PropsWithChildren<{ showLink: boolean }>) {
  return (
    <div className="flex items-center group ">
      {children}
      {showLink ? (
        <IconCozLongArrowTopRight className="ml-auto text-[14px] group-hover:visible invisible" />
      ) : null}
    </div>
  );
}
