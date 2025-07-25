// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useLayoutEffect } from 'react';

import { useInjector } from '@coze-editor/editor/react';
import { astDecorator } from '@coze-editor/editor';
import { EditorView } from '@codemirror/view';

function JinjaHighlight() {
  const injector = useInjector();

  useLayoutEffect(
    () =>
      injector.inject([
        astDecorator.whole.of(cursor => {
          if (
            cursor.name === 'JinjaStatementStart' ||
            cursor.name === 'JinjaStatementEnd'
          ) {
            return {
              type: 'className',
              className: 'jinja-statement-bracket',
            };
          }

          if (cursor.name === 'JinjaComment') {
            return {
              type: 'className',
              className: 'jinja-comment',
            };
          }

          if (cursor.name === 'JinjaExpression') {
            return {
              type: 'className',
              className: 'jinja-expression',
            };
          }
        }),
        EditorView.theme({
          '.jinja-expression': {
            color: 'var(--Green-COZColorGreen7, #00A136)',
          },
        }),
      ]),
    [injector],
  );

  return null;
}

export default JinjaHighlight;
