// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import { ItemErrorType } from '@cozeloop/api-schema/data';
import {
  type DatasetIOJobProgress,
  type ItemErrorGroup,
} from '@cozeloop/api-schema/data';
import { Typography } from '@coze-arch/coze-design';

const ErrorTypeMap = {
  [ItemErrorType.MismatchSchema]: I18n.t('schema_mismatch'),
  [ItemErrorType.EmptyData]: I18n.t('empty_data'),
  [ItemErrorType.ExceedMaxItemSize]: I18n.t('single_data_size_exceeded'),
  [ItemErrorType.ExceedDatasetCapacity]: I18n.t('dataset_capacity_exceeded'),
  [ItemErrorType.MalformedFile]: I18n.t('file_format_error'),
  [ItemErrorType.InternalError]: I18n.t('system_error'),
  [ItemErrorType.IllegalContent]: I18n.t('contains_illegal_content'),
};

export const ImportResultInfo = ({
  progress,
  errors,
}: {
  progress?: DatasetIOJobProgress;
  errors?: ItemErrorGroup[];
}) => (
  <div>
    <div className="flex gap-2 items-center">
      <Typography.Text className="flex-1 leading-[16px]">
        {I18n.t('success')}
        <Typography.Text className="!font-medium mx-1">
          {progress?.added || 0}
        </Typography.Text>
        {I18n.t('item_unit')}， {I18n.t('failure')}
        <Typography.Text className="!font-medium mx-1">
          {Number(progress?.processed) - Number(progress?.added) || 0}
        </Typography.Text>
        {I18n.t('item_unit')}
      </Typography.Text>
    </div>
    {errors?.length ? (
      <div className="mt-2 rounded-[4px] p-2 coz-mg-secondary border border-solid border-[var(--coz-stroke-primary)]">
        <Typography.Text size="small" className="coz-fg-secondary">
          {I18n.t('failure_reasons_and_retry')}
        </Typography.Text>
        {errors.map(log => (
          <div className="flex items-center">
            <span className="rounded-[50%] w-[4px] h-[4px] mx-2 bg-[black]"></span>
            <Typography.Text size="small" className="coz-fg-secondary">
              {ErrorTypeMap[log?.type || ItemErrorType.InternalError]}
              <Typography.Text size="small" className="!font-medium">
                {log?.error_count && log?.error_count > 0
                  ? `（${log?.error_count}${I18n.t('item_unit')}）`
                  : ''}
              </Typography.Text>
            </Typography.Text>
          </div>
        ))}
      </div>
    ) : null}
  </div>
);
