// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { CozeLoopStorage } from '@cozeloop/toolkit';

export {
  LOGIN_PATH,
  CONSOLE_PATH,
  COZELOOP_DOC_URL,
  COZELOOP_GITHUB_URL,
  COZELOOP_DISCORD_URL,
  COZELOOP_LARK_GROUP_URL,
} from './home';

export const storage = new CozeLoopStorage({
  field: 'base',
});
