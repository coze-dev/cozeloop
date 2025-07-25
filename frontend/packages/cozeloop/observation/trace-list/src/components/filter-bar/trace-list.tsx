// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useState } from 'react';

import { EVENT_NAMES, sendEvent } from '@cozeloop/tea-adapter';
import { useSpace } from '@cozeloop/biz-hooks-adapter';
import { Select, Tooltip } from '@coze-arch/coze-design';

import { type LogicValue } from '../logic-expr';

interface SpanTypeSelectProps {
  onChange: (value: string | number) => void;
  value: string | number;
  applyFilters?: LogicValue;
  selectedPlatform: string | number;
  spanListTypeEnumOptionList: {
    label: string;
    value: string | number;
  }[];
  tooltipContent?: Record<string, string>;
}

export const SpanTypeSelect = ({
  onChange,
  value,
  selectedPlatform,
  applyFilters,
  spanListTypeEnumOptionList,
  tooltipContent,
}: SpanTypeSelectProps) => {
  const [showTooltip, setShowTooltip] = useState(false);
  const [showDropdown, setShowDropdown] = useState(false);
  const { spaceID, space: { name: spaceName } = {} } = useSpace();
  return (
    <Tooltip
      content={tooltipContent?.[value]}
      theme="dark"
      visible={showDropdown ? false : showTooltip}
      trigger="custom"
    >
      <Select
        onMouseEnter={() => setShowTooltip(true)}
        onMouseLeave={() => setShowTooltip(false)}
        className="w-[144px] box-border"
        value={value}
        defaultValue={spanListTypeEnumOptionList[0]}
        onDropdownVisibleChange={setShowDropdown}
        optionList={spanListTypeEnumOptionList}
        onSelect={event => {
          onChange(event as string);
          sendEvent(
            EVENT_NAMES.cozeloop_observation_trace_list_span_type_switch,
            {
              platform: selectedPlatform,
              span_type: event,
              filters: JSON.stringify(applyFilters),
              space_id: spaceID,
              space_name: spaceName ?? '',
            },
          );
        }}
      />
    </Tooltip>
  );
};
