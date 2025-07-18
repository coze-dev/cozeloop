// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @coze-arch/use-error-in-catch */
import JSONBig from 'json-bigint';
import Decimal from 'decimal.js';
import { I18n } from '@cozeloop/i18n-adapter';
import { type FieldSchema } from '@cozeloop/api-schema/evaluation';

import { ContentType, COLUMN_TYPE_MAP, DataType } from './type';

const jsonBig = JSONBig({ storeAsString: true });

Decimal.set({ precision: 300 });
Decimal.set({ toExpNeg: -7, toExpPos: 21 });
export const getDataType = (fieldSchema?: FieldSchema) => {
  try {
    const json = JSON.parse(fieldSchema?.text_schema || '{}');
    if (Object.values(DataType).includes(json.type)) {
      return json.type;
    }
    return DataType.String;
  } catch (error) {
    console.error(error);
    return DataType.String;
  }
};

export const getColumnType = (fieldSchema?: FieldSchema) => {
  if (fieldSchema?.content_type === ContentType.Text) {
    return getDataType(fieldSchema);
  }
  return COLUMN_TYPE_MAP[fieldSchema?.content_type || ContentType.Text];
};

export const saftJsonParse = (value?: string) => {
  try {
    return JSON.parse(value || '');
  } catch (error) {
    return '';
  }
};

export const saftJsonBigParse = (jsonStr?: string) => {
  try {
    const parsed = jsonBig.parse(jsonStr || '');
    return parsed;
  } catch (error) {
    return '';
  }
};

export const getSchemaConfig = (schema?: string) => {
  const config = saftJsonBigParse(schema);
  console.log(config);
  return {
    multipleOf: config?.multipleOf,
    maximum: config?.maximum,
    minimum: config?.minimum,
  };
};

export const validateAndFormat = ({
  val,
  minimum,
  maximum,
  multipleOf,
}: {
  val: string;
  minimum?: Decimal;
  maximum?: Decimal;
  multipleOf?: Decimal;
}): string => {
  try {
    //去除str中不符合数字规范的内容，科学技术法要单独保留
    const newStr = val.replace(/[^\d.eE+-]/g, '');
    let decimalValue = new Decimal(newStr);
    // 检查范围
    if (minimum) {
      const minValue = new Decimal(minimum);
      if (decimalValue.lt(minValue)) {
        decimalValue = minValue;
      }
    }
    if (maximum) {
      const maxValue = new Decimal(maximum);
      if (decimalValue.gt(maxValue)) {
        decimalValue = maxValue;
      }
    }
    if (!multipleOf) {
      return decimalValue.toString();
    }
    // 调整到最近的 multipleOf 的倍数
    const multipleOfDecimal = new Decimal(multipleOf);
    decimalValue = decimalValue
      .div(multipleOfDecimal)
      .round()
      .mul(multipleOfDecimal);

    // 确定小数位数
    const decimalPlaces = Math.max(
      0,
      multipleOfDecimal.decimalPlaces(),
      decimalValue.decimalPlaces(),
    );
    let formattedStr = decimalValue.toFixed(decimalPlaces);
    if (formattedStr.includes('.')) {
      formattedStr = formattedStr.replace(/\.?0+$/, '');
    }

    return formattedStr;
  } catch {
    return '';
  }
};

// eslint-disable-next-line complexity -- skip
export const validarDatasetItem = (
  value: string,
  callback: (error?: string) => void,
  fieldSchema?: FieldSchema,
) => {
  const type = getColumnType(fieldSchema);
  if (type !== DataType.Float && type !== DataType.Integer) {
    return true;
  }
  if (!/^-?(?:0|[1-9]\d*)(?:\.\d+)?$/.test(value)) {
    callback(I18n.t('please_input', { field: I18n.t('number') }));
    return false;
  }
  // 校验value 是否为数字；
  let decimalValue;
  try {
    decimalValue = new Decimal(value);
    const { minimum, maximum, multipleOf } = getSchemaConfig(
      fieldSchema?.text_schema,
    );
    const minValue = minimum ? new Decimal(minimum) : undefined;
    const maxValue = maximum ? new Decimal(maximum) : undefined;
    if (minValue && decimalValue.lt(minValue)) {
      callback(I18n.t('input_num_gte', { num: I18n.t('number') }));
      return false;
    }
    if (maxValue && decimalValue.gt(maxValue)) {
      callback(I18n.t('input_num_lte', { num: I18n.t('number') }));
      return false;
    }

    if (type === DataType.Integer && decimalValue.isInteger() === false) {
      callback(I18n.t('please_input', { field: I18n.t('integer') }));
      return false;
    }
    if (type === DataType.Float && multipleOf) {
      const multipleOfDecimal = new Decimal(multipleOf);
      const division = decimalValue.dividedBy(multipleOfDecimal);
      if (!division.isInteger()) {
        callback(
          I18n.t('support_precision', {
            precision: multipleOfDecimal.decimalPlaces(),
          }),
        );
        return false;
      }
    }
    return true;
  } catch (error) {
    callback(I18n.t('please_input', { field: I18n.t('number') }));
    return false;
  }
};
