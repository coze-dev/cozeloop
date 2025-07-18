// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @coze-arch/max-line-per-function */
import { I18n } from '@cozeloop/i18n-adapter';
import { TooltipWhenDisabled } from '@cozeloop/components';
import {
  IconCozArrowDown,
  IconCozArrowRight,
  IconCozCopy,
  IconCozTrashCan,
} from '@coze-arch/coze-design/icons';
import {
  Button,
  Collapse,
  FormInput,
  FormSelect,
  Popconfirm,
  Tooltip,
  Typography,
  useFieldApi,
} from '@coze-arch/coze-design';

import {
  DISPLAY_TYPE_MAP,
  DATA_TYPE_LIST,
  type DataType,
  DISPLAY_FORMAT_MAP,
} from '../dataset-item/type';
import { columnNameRuleValidator } from '../../utils/source-name-rule';

interface ColumnRenderProps {
  fieldKey: string;
  index: number;
  onDelete: () => void;
  onCopy: () => void;
  size?: 'large' | 'small';
  activeKey: string[];
  setActiveKey: (key: string[]) => void;
  disabledDataTypeSelect?: boolean;
}

export const ColumnRender = ({
  fieldKey,
  index,
  onDelete,
  onCopy,
  size = 'large',
  activeKey,
  setActiveKey,
  disabledDataTypeSelect,
}: ColumnRenderProps) => {
  const typeField = useFieldApi(`${fieldKey}.${index}.type`);
  const keyField = useFieldApi(`${fieldKey}.${index}.key`);
  const nameField = useFieldApi(`${fieldKey}.${index}.name`);
  const displayFormatField = useFieldApi(
    `${fieldKey}.${index}.default_display_format`,
  );
  const allColumnField = useFieldApi(fieldKey);
  const type = typeField.getValue() as DataType;
  const isExist = keyField.getValue() !== undefined;
  const getHeader = () => (
    <div className="flex w-full justify-between">
      <div className="flex items-center gap-[4px]">
        <Typography.Text className="text-[14px] !font-semibold">
          {nameField.getValue() || I18n.t('column_index', { index: index + 1 })}
        </Typography.Text>
        {activeKey.includes(`${index}`) ? (
          <IconCozArrowDown
            onClick={() =>
              setActiveKey(activeKey.filter(key => key !== `${index}`))
            }
            className="cursor-pointer w-[16px] h-[16px]"
          />
        ) : (
          <IconCozArrowRight
            onClick={() => setActiveKey([...activeKey, `${index}`])}
            className="cursor-pointer w-[16px] h-[16px]"
          />
        )}
      </div>
      <div
        onClick={e => e.stopPropagation()}
        className="group-hover:block hidden"
      >
        <Tooltip content={I18n.t('copy')} theme="dark" className="mr-[2px]">
          <Button
            color="secondary"
            size="mini"
            icon={<IconCozCopy className="w-[14px] h-[14px]" />}
            onClick={() => onCopy()}
          ></Button>
        </Tooltip>
        {isExist ? (
          <Popconfirm
            content={
              <Typography.Text className="break-all text-[12px] !coz-fg-secondary">
                {I18n.t('confirm_delete_x_columns', {
                  num: (
                    <Typography.Text className="!font-medium">
                      {nameField.getValue()}
                    </Typography.Text>
                  ),
                })}
              </Typography.Text>
            }
            title={I18n.t('delete_column')}
            okText={I18n.t('delete')}
            zIndex={1062}
            okButtonProps={{
              color: 'red',
            }}
            cancelText={I18n.t('Cancel')}
            style={{ width: 280 }}
            onConfirm={() => {
              onDelete();
            }}
          >
            <Button
              icon={<IconCozTrashCan className="w-[14px] h-[14px]" />}
              color="secondary"
              size="mini"
            ></Button>
          </Popconfirm>
        ) : (
          <Tooltip content={I18n.t('delete')} theme="dark">
            <Button
              icon={<IconCozTrashCan className="w-[14px] h-[14px]" />}
              color="secondary"
              size="mini"
              onClick={() => onDelete()}
            ></Button>
          </Tooltip>
        )}
      </div>
    </div>
  );

  return (
    <Collapse.Panel
      className="group"
      itemKey={`${index}`}
      header={getHeader()}
      showArrow={false}
    >
      <div className="flex flex-col justify-stretch">
        <div className="flex gap-[20px]">
          <FormInput
            fieldClassName="flex-1"
            label={I18n.t('column_name')}
            placeholder={I18n.t('please_input', {
              field: I18n.t('column_name'),
            })}
            maxLength={50}
            autoComplete=""
            field={`${fieldKey}.${index}.name`}
            rules={[
              {
                required: true,
                message: I18n.t('please_input', {
                  field: I18n.t('column_name'),
                }),
              },
              {
                validator: columnNameRuleValidator,
              },
              {
                validator: (_, value) => {
                  if (!value) {
                    return true;
                  }
                  const allColumnData = allColumnField.getValue();
                  // 判断之前的列名称中是否有与自己相同的name
                  const hasSameName = allColumnData
                    ?.slice(0, index)
                    ?.some(
                      (data, dataIndex) =>
                        dataIndex !== index && data.name === value,
                    );

                  return !hasSameName;
                },
                message: I18n.t('field_exists', {
                  field: I18n.t('column_name'),
                }),
              },
            ]}
          ></FormInput>
          <TooltipWhenDisabled
            disabled={disabledDataTypeSelect && isExist}
            content={I18n.t('cannot_modify_data_type_tip')}
            theme="dark"
            className="top-9"
          >
            <FormSelect
              label={I18n.t('data_type')}
              labelWidth={90}
              zIndex={1070}
              fieldClassName="w-[190px]"
              disabled={disabledDataTypeSelect && isExist}
              optionList={DATA_TYPE_LIST}
              onChange={newType => {
                displayFormatField.setValue(
                  DISPLAY_TYPE_MAP?.[newType as DataType]?.[0],
                );
              }}
              field={`${fieldKey}.${index}.type`}
              className="w-full"
              rules={[
                {
                  required: true,
                  message: I18n.t('please_select', {
                    field: I18n.t('data_type'),
                  }),
                },
              ]}
            ></FormSelect>
          </TooltipWhenDisabled>
          <div>
            <FormSelect
              label={I18n.t('view_format')}
              zIndex={1070}
              labelWidth={90}
              disabled={DISPLAY_TYPE_MAP[type]?.length <= 1}
              fieldClassName="w-[190px]"
              field={`${fieldKey}.${index}.default_display_format`}
              className={'w-full '}
              optionList={DISPLAY_TYPE_MAP[type]?.map(item => ({
                label: DISPLAY_FORMAT_MAP[item],
                value: item,
              }))}
              rules={[
                {
                  required: true,
                  message: I18n.t('please_select', {
                    field: I18n.t('view_format'),
                  }),
                },
              ]}
            ></FormSelect>
          </div>
        </div>
        <div className="flex-grow-1">
          <FormInput
            label={I18n.t('column_description')}
            placeholder={I18n.t('please_input', {
              field: I18n.t('column_description'),
            })}
            maxLength={200}
            field={`${fieldKey}.${index}.description`}
            autoComplete="off"
          ></FormInput>
        </div>
      </div>
    </Collapse.Panel>
  );
};
