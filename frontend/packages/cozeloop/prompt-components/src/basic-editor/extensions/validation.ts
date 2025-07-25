// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable security/detect-non-literal-regexp */
/* eslint-disable @typescript-eslint/naming-convention */
import { type ReactNode, useLayoutEffect } from 'react';

import { useInjector } from '@coze-editor/editor/react';
import { astDecorator } from '@coze-editor/editor';

import { VARIABLE_MAX_LEN } from '@/consts';

import styles from './validation.module.css';

function validate(text: string) {
  const regex = new RegExp(`^[a-zA-Z][\\w]{0,${VARIABLE_MAX_LEN - 1}}$`, 'gm');
  if (regex.test(text)) {
    return true;
  }

  return false;
}

export function Validation(): ReactNode {
  const injector = useInjector();

  // 用于校验 {{ }} 中的变量，如果变量无效，使用灰色标识
  useLayoutEffect(
    () =>
      injector.inject([
        astDecorator.whole.of((cursor, state) => {
          if (
            cursor.name === 'JinjaExpression' &&
            cursor.node.firstChild?.name === 'JinjaExpressionStart' &&
            cursor.node.lastChild?.name === 'JinjaExpressionEnd'
          ) {
            const from = cursor.node.firstChild.to;
            const to = cursor.node.lastChild.from;
            const text = state.sliceDoc(from, to);
            if (validate(text)) {
              return {
                type: 'className',
                className: styles.valid,
                from,
                to,
              };
            }

            return {
              type: 'className',
              className: styles.invalid,
              from,
              to,
            };
          }
        }),
      ]),
    [injector],
  );

  return null;
}

export default Validation;
