// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import classNames from 'classnames';
import { SpanType } from '@cozeloop/api-schema/observation';
import {
  IconCozCopy,
  IconCozCrossFill,
  IconCozArrowLeft,
  IconCozArrowRight,
} from '@coze-arch/coze-design/icons';
import { Button, Typography, Divider, Tooltip } from '@coze-arch/coze-design';

import { getNodeConfig } from '@/trace-detail/utils/span';
import { useDetailCopy } from '@/trace-detail/hooks/use-detail-copy';
import {
  BROKEN_ROOT_SPAN_ID,
  NORMAL_BROKEN_SPAN_ID,
} from '@/trace-detail/consts/span';

import { type TraceHeaderProps } from '../typing';

import styles from './index.module.less';
import { I18n } from '@cozeloop/i18n-adapter';

export const HorizontalTraceHeader = ({
  rootSpan,
  className: propsClassName,
  showClose,
  onClose,
  switchConfig,
  showTraceId = true,
}: TraceHeaderProps & { showTraceId?: boolean }) => {
  const { type, span_name = '', trace_id, span_id, span_type } = rootSpan || {};
  const handleCopy = useDetailCopy();
  const isBroken = [NORMAL_BROKEN_SPAN_ID, BROKEN_ROOT_SPAN_ID].includes(
    span_id || '',
  );

  const traceName = isBroken ? 'Unknown Trace' : span_name;

  const nodeConfig = getNodeConfig({
    spanTypeEnum: type ?? SpanType.Unknown,
    spanType: span_type ?? SpanType.Unknown,
  });

  return (
    <div className={classNames(styles['horizontal-header'], propsClassName)}>
      <div className={styles['trace-profile']}>
        {
          <span className={styles['icon-wrapper']}>
            {nodeConfig.icon?.({
              className: '!w-[16px] !h-[16px]',
              size: 'large',
            })}
          </span>
        }
        <div className={styles.desc}>
          <div className={styles.name}>
            <Typography.Text
              ellipsis={{
                showTooltip: {
                  type: 'tooltip',
                  opts: {
                    position: 'bottom',
                    theme: 'dark',
                  },
                },
              }}
              className="coz-fg-plus text-[16px]"
            >
              {traceName}
            </Typography.Text>
          </div>
        </div>
        {showTraceId ? (
          <div
            onClick={() => {
              handleCopy(trace_id || '', 'trace_id');
            }}
            className="flex items-center gap-x-1 ml-1"
          >
            <Typography.Text
              type="secondary"
              size="small"
              className="cursor-pointer coz-fg-primary text-[14px]"
            >
              {I18n.t('trace_id')}
            </Typography.Text>

            <Tooltip theme="dark" content={I18n.t('Copy')}>
              <Button
                size="small"
                className="w-[24px] h-[24px]"
                color="secondary"
                icon={
                  <IconCozCopy className="!w-[14px] !h-[14px] coz-fg-secondary" />
                }
              />
            </Tooltip>
          </div>
        ) : null}
      </div>
      <div className="flex-1 flex justify-end items-center">
        {switchConfig ? (
          <>
            <Button
              icon={<IconCozArrowLeft />}
              color="secondary"
              disabled={!switchConfig?.canSwitchPre}
              className="text-[13px] coz-fg-secondary"
              size="default"
              onClick={() => {
                switchConfig?.onSwitch('pre');
              }}
            >
              {I18n.t('prev_item')}
            </Button>
            <Button
              icon={<IconCozArrowRight />}
              iconPosition="right"
              className="text-[13px] coz-fg-secondary ml-2"
              color="secondary"
              size="default"
              disabled={!switchConfig?.canSwitchNext}
              onClick={() => {
                switchConfig?.onSwitch('next');
              }}
            >
              {I18n.t('next_item')}
            </Button>
          </>
        ) : null}

        {switchConfig ? (
          <Divider className="mx-2 h-[12px]" layout="vertical" />
        ) : null}

        {showClose ? (
          <Button
            color="secondary"
            onClick={onClose}
            className="!h-[32px] !w-[32px]"
          >
            <IconCozCrossFill className="!w-[16px] !h-[16px]" />
          </Button>
        ) : null}
      </div>
    </div>
  );
};
