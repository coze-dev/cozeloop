// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @coze-arch/max-line-per-function */
/* eslint-disable complexity */
/* eslint-disable @typescript-eslint/no-explicit-any */
import React, { useEffect, useMemo, useState } from 'react';

import classNames from 'classnames';
import { safeParseJson } from '@cozeloop/toolkit';
import { BaseJsonEditor, BaseRawTextEditor } from '@cozeloop/prompt-components';
import { handleCopy, TooltipWhenDisabled } from '@cozeloop/components';
import { ToolType } from '@cozeloop/api-schema/prompt';
import { IconCozCopy } from '@coze-arch/coze-design/icons';
import {
  Button,
  Col,
  Modal,
  Row,
  Select,
  Space,
  Toast,
  Typography,
} from '@coze-arch/coze-design';

import { type ToolWithMock } from '.';
import { I18n } from '@cozeloop/i18n-adapter';

interface ToolModalProps {
  visible?: boolean;
  data?: ToolWithMock;
  disabled?: boolean;
  tools?: ToolWithMock[];
  onConfirm?: (
    data: ToolWithMock,
    isUpdate?: boolean,
    oldData?: ToolWithMock,
  ) => void;
  onClose?: () => void;
}

interface SchemaEditorProps {
  value?: string;
  readOnly?: boolean;
  onChange?: (value?: string) => void;
  language?: string;
  placeholder?: string;
  showLineNumbs?: boolean;
  className?: string;
}

const TEMPLATE_DATA = `{
  "name": "get_weather",
  "description": "Determine weather in my location",
  "parameters": {
    "type": "object",
    "properties": {
      "location": {
        "type": "string",
        "description": "The city and state e.g. San Francisco, CA"
      },
      "unit": {
        "type": "string",
        "enum": [
          "c",
          "f"
        ]
      }
    },
    "required": [
      "location"
    ]
  }
}`;

export const SchemaEditor = ({
  value,
  onChange,
  placeholder,
  readOnly,
  language,
  className,
}: SchemaEditorProps) => (
  <div
    className={classNames(
      'w-full h-[500px] border border-solid coz-stroke-primary rounded-[4px] overflow-hidden relative bg-white',
      className,
    )}
  >
    {language === 'json' ? (
      <BaseJsonEditor
        className="w-full h-full overflow-y-auto"
        onChange={onChange}
        value={value || ''}
        placeholder={placeholder}
        readonly={readOnly}
        borderRadius={4}
      />
    ) : (
      <BaseRawTextEditor
        className="w-full h-full overflow-y-auto"
        onChange={onChange}
        value={value || ''}
        placeholder={placeholder}
        readonly={readOnly}
      />
    )}
  </div>
);

interface ToolSchemaProps {
  name?: string;
  description?: string;
  parameters?: any;
}
export function ToolModal({
  visible,
  disabled,
  data,
  onClose,
  onConfirm,
  tools,
}: ToolModalProps) {
  const [mockType, setMockType] = useState('text');

  const toolSchema = useMemo(() => {
    if (data?.function) {
      const toolObj: ToolSchemaProps = {
        name: data.function.name,
        description: data.function.description,
      };
      if (data.function.parameters) {
        toolObj.parameters = safeParseJson(data.function.parameters);
      }
      return JSON.stringify(toolObj, null, 2);
    }

    return '';
  }, [JSON.stringify(data || {})]);

  const [schema, setSchema] = useState<string>();
  const [mockValue, setMockValue] = useState<string>();
  const isCreate = !data;

  const canSaveTool = useMemo(() => {
    if (schema) {
      const schemaObj = safeParseJson<ToolSchemaProps>(schema);
      if (
        !schemaObj?.name ||
        !/^[a-zA-Z][a-zA-Z0-9_-]{0,63}$/.test(schemaObj.name)
      ) {
        return false;
      }
      return true;
    }
  }, [schema, JSON.stringify(tools), isCreate]);

  const handleSaveTool = () => {
    if (disabled) {
      onClose?.();
    }
    if (!schema) {
      return;
    }

    const schemaObj = safeParseJson<ToolSchemaProps>(schema);
    const toolObj = {
      name: schemaObj?.name,
      description: schemaObj?.description,
      parameters: schemaObj?.parameters
        ? JSON.stringify(schemaObj?.parameters)
        : '',
    };
    const tool = {
      type: ToolType.Function,
      function: toolObj,
      mock_response: mockValue,
    };

    const hasItem =
      tools?.find(it => it?.function?.name === toolObj?.name) &&
      data?.function?.name !== toolObj?.name;

    if (hasItem) {
      Toast.warning({
        content: I18n.t('method_exists'),
        zIndex: 99999,
      });
      return;
    }
    onConfirm?.(tool, !isCreate, data);
  };

  useEffect(() => {
    setSchema(toolSchema);
  }, [toolSchema]);

  useEffect(() => {
    setMockValue(data?.mock_response);
  }, [data?.mock_response]);

  useEffect(() => {
    if (!visible) {
      setSchema(undefined);
      setMockValue(undefined);
    }
  }, [visible]);

  return (
    <Modal
      title={data?.function?.name || I18n.t('new_function')}
      width={960}
      visible={visible}
      onCancel={onClose}
      okButtonProps={{ disabled: !canSaveTool }}
      maskClosable={false}
      footer={
        disabled ? null : (
          <Space>
            <Button className="mr-2" onClick={onClose} color="primary">
              {I18n.t('cancel')}
            </Button>
            <TooltipWhenDisabled
              content={I18n.t('method_name_rule')}
              disabled={Boolean(schema && !canSaveTool)}
            >
              <Button onClick={handleSaveTool} disabled={!canSaveTool}>
                {I18n.t('confirm')}
              </Button>
            </TooltipWhenDisabled>
          </Space>
        )
      }
      hasScroll={false}
    >
      <Row gutter={16}>
        <Col span={14}>
          <div className="flex justify-between items-center w-full h-8 mb-2">
            <Typography.Text className="font-semibold items" type="tertiary">
              SCHEMA
              <IconCozCopy
                className="ml-2 hover:text-semi-primary cursor-pointer"
                onClick={() => handleCopy(schema || '')}
              />
            </Typography.Text>
            {disabled ? null : (
              <Button
                size="small"
                onClick={() => {
                  setSchema(TEMPLATE_DATA);
                  setMockType('text');
                  setMockValue('Sunny');
                }}
              >
                {I18n.t('insert_template')}
              </Button>
            )}
          </div>
          <SchemaEditor
            language="json"
            value={schema}
            onChange={v => setSchema(v)}
            showLineNumbs
            readOnly={disabled}
          />
        </Col>
        <Col span={10}>
          <div className="flex justify-between items-center w-full h-8 mb-2">
            <Typography.Text className="font-semibold" type="tertiary">
              {I18n.t('default_mock_value')}
            </Typography.Text>
            <Select
              value={mockType}
              onChange={v => setMockType(v as string)}
              size="small"
              zIndex={2001}
            >
              <Select.Option value="text">TEXT</Select.Option>
              <Select.Option value="json">JSON</Select.Option>
            </Select>
          </div>
          <SchemaEditor
            language={mockType}
            value={mockValue}
            onChange={v => setMockValue(v)}
            placeholder={I18n.t('input_mock_value_here')}
            readOnly={disabled}
          />
        </Col>
      </Row>
    </Modal>
  );
}
