// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { Typography } from '@coze-arch/coze-design';

import { type Span, SpanType } from '@/trace-detail/typings/params';

import { getNodeConfig } from '../../utils/span';
import { CustomIconWrapper } from '../../consts/span';

import styles from './index.module.less';
export const SpanDetailHeader: React.FC<{
  span: Span;
  moduleName?: string;
}> = ({ span }) => {
  const { type, span_type } = span;
  const nodeConfig = getNodeConfig({
    spanTypeEnum: type ?? SpanType.Unknown,
    spanType: span_type,
  });

  return (
    <div className={styles['detail-header']}>
      <div className={styles['detail-title']}>
        <span className={styles['icon-wrapper']}>
          {nodeConfig.icon ? (
            nodeConfig.icon({ className: '!w-[16px] !h-[16px]', size: 'large' })
          ) : (
            <CustomIconWrapper color={nodeConfig.color} size="large">
              {nodeConfig.character}
            </CustomIconWrapper>
          )}
        </span>
        <Typography.Text
          ellipsis={{ rows: 1 }}
          className="text-[16px] !font-semibold"
        >
          {span.span_name}
        </Typography.Text>
      </div>
    </div>
  );
};
