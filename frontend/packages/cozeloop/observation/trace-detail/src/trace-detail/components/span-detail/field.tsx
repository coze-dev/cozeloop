// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { isEmpty } from 'lodash-es';
import dayjs from 'dayjs';
import { I18n } from '@cozeloop/i18n-adapter';
import { SpanStatus, SpanType } from '@cozeloop/api-schema/observation';
import {
  IconCozCheckMarkCircleFillPalette,
  IconCozCrossCircleFill,
} from '@coze-arch/coze-design/icons';
import { Tooltip, Tag } from '@coze-arch/coze-design';

import { type Span } from '@/trace-detail/typings/params';

import { type Field } from '../common/field-list';
import { getNodeConfig } from '../../utils/span';

export function geSpanOverviewField(span: Span): Field[] {
  const {
    span_id,
    duration,
    started_at,
    status,
    type,
    custom_tags,
    span_type,
    status_code,
  } = span;
  const { tokens, input_tokens, output_tokens } = custom_tags ?? {};

  const nodeConfig = getNodeConfig({
    spanTypeEnum: type ?? SpanType.Unknown,
    spanType: span_type,
  });
  const isSuccess = status === SpanStatus.Success;
  return [
    {
      key: 'status',
      title: 'Status',
      item: (
        <Tag
          color={isSuccess ? 'green' : 'red'}
          prefixIcon={
            isSuccess ? (
              <IconCozCheckMarkCircleFillPalette />
            ) : (
              <IconCozCrossCircleFill />
            )
          }
        >
          {isSuccess ? 'Success' : 'Error'}
        </Tag>
      ),
      width: 224,
    },
    {
      key: 'status_code',
      title: I18n.t('observation_status_code'),
      item:
        status !== SpanStatus.Success ? (
          <span
            style={{
              color: '#FF441E',
            }}
          >
            {status_code}
          </span>
        ) : (
          '-'
        ),
      width: 224,
    },
    {
      key: 'span_id',
      title: I18n.t('observation_query_detail_span_id'),
      item: span_id || '-',
      width: 224,
      enableCopy: true,
    },
    {
      key: 'type',
      title: 'Type',
      item: span_type || nodeConfig.typeName || '-',
      width: 224,
    },
    {
      key: 'latency',
      title: 'Latency',
      item:
        duration === undefined ? '-' : `${Number(duration).toLocaleString()}ms`,
      width: 224,
    },
    {
      key: 'start_time',
      title: 'StartTime',
      item: dayjs(Number(started_at)).format('YYYY-MM-DD HH:mm:ss'),
      width: 240,
    },
    {
      key: 'tokens',
      title: 'Tokens',
      item: tokens && (
        <Tooltip
          theme="dark"
          content={
            <>
              {input_tokens !== undefined && (
                <div>
                  {I18n.t('observation_input_tokens_count', {
                    count: Number(input_tokens),
                  })}
                </div>
              )}
              {output_tokens !== undefined && (
                <div>
                  {I18n.t('observation_output_tokens_count', {
                    count: Number(output_tokens),
                  })}
                </div>
              )}
            </>
          }
        >
          {Number(tokens)}
        </Tooltip>
      ),
      width: 240,
    },
  ].filter(field => !isEmpty(field.item)) satisfies Field[];
}
