// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { EditorProvider, MergeViewRenderer } from '@coze-editor/editor/react';
import preset from '@coze-editor/editor/preset-prompt';
import { EditorView } from '@codemirror/view';
import { type Extension } from '@codemirror/state';

import MarkdownHighlight from './extensions/markdown';
import LanguageSupport from './extensions/language-support';
import JinjaHighlight from './extensions/jinja';

const extensions: Extension[] = [
  // diff theme
  EditorView.theme({
    '&.cm-merge-b .cm-changedLine': {
      background: 'transparent !important',
      paddingLeft: '12px !important',
    },
    '&.cm-merge-b .cm-changedText': {
      background:
        'linear-gradient(#22bb2266, #22bb2266) bottom/100% 100% no-repeat',
    },
    '&.cm-merge-a .cm-changedText, .cm-deletedChunk .cm-deletedText': {
      background:
        'linear-gradient(#ee443366, #ee443366) bottom/100% 100% no-repeat',
    },
  }),
];

interface PromptDiffEditorProps {
  oldValue?: string;
  newValue?: string;
  autoScrollToBottom?: boolean;
}

export function PromptDiffEditor({
  oldValue,
  newValue,
  autoScrollToBottom,
}: PromptDiffEditorProps) {
  return (
    <EditorProvider>
      <MergeViewRenderer
        plugins={preset}
        domProps={{
          style: {
            flex: 1,
            fontSize: 12,
          },
        }}
        mergeConfig={{
          gutter: false,
        }}
        a={{
          defaultValue: oldValue,
          extensions,
          options: {
            editable: false,
            readOnly: true,
          },
        }}
        b={{
          defaultValue: newValue,
          extensions,
          options: {
            editable: false,
            readOnly: true,
          },
        }}
        didMount={editor => {
          if (autoScrollToBottom) {
            editor.b.$view.dispatch({
              effects: EditorView.scrollIntoView(
                editor.b.$view.state.doc.length,
              ),
            });
          }
        }}
      />
      <LanguageSupport />
      <JinjaHighlight />
      <MarkdownHighlight />
    </EditorProvider>
  );
}
