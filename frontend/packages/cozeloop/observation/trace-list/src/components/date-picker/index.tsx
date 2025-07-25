// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { forwardRef, type ReactNode, useImperativeHandle, useRef } from 'react';

import dayjs from 'dayjs';
import {
  IconCozArrowDown,
  IconCozCalendar,
} from '@coze-arch/coze-design/icons';
import { type DatePickerProps } from '@coze-arch/coze-design';
import { DatePicker, Select, InputGroup } from '@coze-arch/coze-design';

import { CURRENT_TIMEZONE } from '../../utils/dayjs';
import { calcPresetTime, PresetRange, YEAR_DAY_COUNT } from '../../consts/time';

import styles from './index.module.less';

export interface TimeStamp {
  startTime: number;
  endTime: number;
}

interface PreselectedDatePickerOption {
  label: ReactNode;
  value: number | string;
  disabled?: boolean;
}

interface PreselectedDatePickerProps {
  preset: PresetRange;
  timeStamp: TimeStamp;
  datePickerOptions: PreselectedDatePickerOption[];
  maxPastDateRange?: number;
  datePickerProps?: DatePickerProps;
  onPresetChange: (preset: PresetRange, presetTimeStamp?: TimeStamp) => void;
  oneTimeStampChange: (timeStamp: TimeStamp) => void;
}

export interface PreselectedDatePickerRef {
  getCurrentTime: () =>
    | {
        startTime: number;
        endTime: number;
      }
    | undefined;
  closeSelect: () => void;
  closeDatePicker: () => void;
}

export const PreselectedDatePicker = forwardRef<
  PreselectedDatePickerRef,
  PreselectedDatePickerProps
>((props, ref) => {
  const {
    preset,
    timeStamp: { startTime, endTime },
    maxPastDateRange,
    datePickerOptions,
    datePickerProps,
    onPresetChange,
    oneTimeStampChange,
  } = props;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const selectRef = useRef<any>(null);
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const datePickerRef = useRef<any>(null);
  useImperativeHandle(
    ref,
    () => ({
      getCurrentTime: () => {
        const time = calcPresetTime(preset as PresetRange);
        return time;
      },
      closeSelect: () => {
        selectRef.current?.close?.();
      },
      closeDatePicker: () => {
        datePickerRef.current?.close?.();
      },
    }),
    [preset],
  );

  // 格式化时间戳显示函数
  const formatTimeDisplay = () => {
    if (!startTime || !endTime) {
      return '';
    }

    const startFormatted = dayjs(startTime)
      .tz(CURRENT_TIMEZONE)
      .format('YYYY-MM-DD HH:mm:ss');
    const endFormatted = dayjs(endTime)
      .tz(CURRENT_TIMEZONE)
      .format('YYYY-MM-DD HH:mm:ss');

    return `${startFormatted} ~ ${endFormatted}`;
  };

  // 自定义日期选择器渲染函数
  const customTriggerRender = () => (
    <div className={styles.datePickerTrigger}>
      <Select
        value={formatTimeDisplay()}
        showArrow={false}
        showClear={false}
        emptyContent={null}
        disabled={preset !== PresetRange.Unset}
        suffix={
          <IconCozArrowDown className="coz-fg-secondary !mx-2 w-[14px] h-[14px] !text-[14px]" />
        }
        prefix={
          <IconCozCalendar className="coz-fg-secondary mx-2 w-3.5 h-3.5" />
        }
      ></Select>
    </div>
  );

  return (
    <InputGroup className="box-border !flex-nowrap">
      <Select
        ref={selectRef}
        className="min-w-[136px] max-w-[136px] box-border"
        maxHeight={320}
        dropdownClassName={styles.presetTimeDropdown}
        optionList={datePickerOptions}
        value={preset}
        onChange={v => {
          const time = calcPresetTime(v as PresetRange);
          onPresetChange(v as PresetRange, time);
        }}
      />
      <DatePicker
        ref={datePickerRef}
        className={styles.datePicker}
        disabledDate={date => {
          if (date && date.getTime() > dayjs().endOf('day').valueOf()) {
            return true;
          }
          const dayCount = dayjs().diff(dayjs(date), 'days');
          return dayCount > (maxPastDateRange ?? YEAR_DAY_COUNT);
        }}
        disabled={preset !== PresetRange.Unset}
        type="dateTimeRange"
        value={[startTime, endTime]}
        timeZone={CURRENT_TIMEZONE}
        onChange={range => {
          const [start, end] = (range || []) as string[] | Date[];
          oneTimeStampChange({
            startTime: new Date(start).valueOf(),
            endTime: new Date(end).valueOf(),
          });
        }}
        showClear={false}
        triggerRender={customTriggerRender}
        {...datePickerProps}
      />
    </InputGroup>
  );
});
