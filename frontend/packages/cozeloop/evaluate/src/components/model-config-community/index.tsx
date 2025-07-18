// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import classNames from 'classnames';
import { PopoverModelConfigEditorQuery } from '@cozeloop/prompt-components';
import { type ModelConfigEditorProps } from '@cozeloop/evaluate-components';
import { type Model } from '@cozeloop/api-schema/llm-manage';
import { IconCozSetting } from '@coze-arch/coze-design/icons';
import { I18n } from '@cozeloop/i18n-adapter';

export function ModelConfigCommunity(props: ModelConfigEditorProps) {
  const renderDisplayContent = (
    selectModel?: Model,
    isPopoverVisible?: boolean,
  ) => (
    <div
      className={classNames(
        'flex flex-row items-center h-8 border border-solid coz-stroke-plus rounded-[6px] px-2 hover:coz-mg-primary-hovered active:coz-mg-primary-pressed active:coz-stroke-hglt cursor-pointer',
        {
          '!coz-stroke-hglt': isPopoverVisible,
        },
      )}
    >
      <div className="flex-1 text-sm coz-fg-primary font-normal">
        {selectModel ? (
          selectModel?.name
        ) : (
          <span className="coz-fg-dim">{I18n.t('choose_model')}</span>
        )}
      </div>
      {props.disabled ? (
        <div className="flex-shrink-0 text-sm text-brand-9 font-normal cursor-pointer">
          {I18n.t('check_parameters')}
        </div>
      ) : (
        <IconCozSetting className="flex-shrink-0 w-4 h-4 ml-6px coz-fg-secondary" />
      )}
    </div>
  );

  return (
    <PopoverModelConfigEditorQuery
      key={props.refreshModelKey}
      value={props.value}
      onChange={v => {
        props.onChange?.(v);
      }}
      disabled={props.disabled}
      renderDisplayContent={renderDisplayContent}
      defaultActiveFirstModel={true}
      popoverProps={
        props?.popoverProps ?? {
          position: 'bottomLeft',
        }
      }
    />
  );
}
