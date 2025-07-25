// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import { type RuleItem } from '@coze-arch/coze-design';

export const sourceNameRuleValidator: RuleItem['validator'] = (
  rule,
  value,
  callback,
) => {
  // 复合正则表达式验证 [4,7](@ref)
  const pattern = /^[a-zA-Z0-9\u4e00-\u9fa5][\w\u4e00-\u9fa5\-\.]*$/;
  if (!pattern.test(value)) {
    // 错误类型细分 [2,5](@ref)
    const firstChar = value.charAt(0);
    console.log(firstChar);
    if (/^[-_.]/.test(firstChar)) {
      callback(
        I18n.t(
          'only_support_english_letters_numbers_and_chinese_at_the_beginning',
        ),
      );
    } else {
      callback(I18n.t('only_support_english_letters_numbers_and_chinese_and_'));
    }
  }
  return true;
};

export const columnNameRuleValidator: RuleItem['validator'] = (
  rule,
  value,
  callback,
) => {
  if (!/^[a-zA-Z][a-zA-Z0-9_]*$/.test(value)) {
    callback(
      I18n.t(
        'only_support_english_numbers_and_underscores_and_start_with_a_letter',
      ),
    );
  }
  return true;
};
