// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
const MAX_NAME_LENGTH = 20;
export const validateViewName = (name: string, viewNames: string[]) => {
  if (name.trim() === '') {
    return {
      isValid: false,
      message: '不允许为空',
    };
  }

  if (name.trim().length > MAX_NAME_LENGTH) {
    return {
      isValid: false,
      message: `名称长度不能超过${MAX_NAME_LENGTH}个字符`,
    };
  }
  if (viewNames.includes(name.trim())) {
    return {
      isValid: false,
      message: '视图名称已存在',
    };
  }
  return {
    isValid: true,
    message: '',
  };
};
