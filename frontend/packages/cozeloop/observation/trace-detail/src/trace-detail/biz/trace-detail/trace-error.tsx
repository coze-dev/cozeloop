// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import {
  IconCozIllusNone,
  IconCozIllusLock,
} from '@coze-arch/coze-design/illustrations';
import { Empty } from '@coze-arch/coze-design';

import { TRACE_EXPIRED_CODE } from '@/trace-detail/consts/code';
import { HorizontalTraceHeader } from '@/trace-detail/components/header';

export const getEmptyConfig = (statusCode: number) => {
  switch (statusCode) {
    case TRACE_EXPIRED_CODE:
      return {
        image: (
          <IconCozIllusNone className="text-[120px] w-[120px] h-[120px]" />
        ),
        description: I18n.t('current_trace_expired_to_view'),
      };
    default:
      return {
        image: (
          <IconCozIllusLock className="text-[120px] w-[120px] h-[120px]" />
        ),
        description: I18n.t('no_permission_to_view_trace'),
        title: I18n.t('no_permission_to_view'),
      };
  }
};

interface TraceDetailErrorProps {
  statusCode: number;
  spaceID: string;
  id?: string;
  searchType: string;
  headerConfig?: {
    onClose?: () => void;
    minColWidth?: number;
  };
}
export const TraceDetailError = (props: TraceDetailErrorProps) => {
  const { statusCode, spaceID, id, searchType } = props;
  const emptyConfig = getEmptyConfig(statusCode);
  return (
    <div className="w-full h-full">
      <div className="border-solid border-0 border-b border-[var(--coz-stroke-primary)]">
        <HorizontalTraceHeader
          showClose
          onClose={props.headerConfig?.onClose}
          minColWidth={props.headerConfig?.minColWidth}
          showTraceId={false}
          urlParams={{
            spaceID,
            id: id ?? '',
            searchType,
          }}
        />
      </div>
      <div className="flex-1 flex items-center justify-center h-full">
        <Empty {...emptyConfig} />
      </div>
    </div>
  );
};
