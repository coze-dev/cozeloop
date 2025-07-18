// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @typescript-eslint/no-explicit-any */
import { forwardRef, useImperativeHandle, useMemo, useRef } from 'react';

import { I18n } from '@cozeloop/i18n-adapter';
import { type VariableDef } from '@cozeloop/api-schema/prompt';
import {
  EditorProvider,
  Renderer,
  Placeholder,
} from '@coze-editor/editor/react';
import preset from '@coze-editor/editor/preset-prompt';
import { EditorView } from '@codemirror/view';
import { type Extension } from '@codemirror/state';

import Variable from './extensions/variable';
import Validation from './extensions/validation';
import MarkdownHighlight from './extensions/markdown';
import LanguageSupport from './extensions/language-support';
import JinjaHighlight from './extensions/jinja';
import { goExtension } from './extensions/go-template';

export interface PromptBasicEditorProps {
  defaultValue?: string;
  height?: number;
  minHeight?: number;
  maxHeight?: number;
  fontSize?: number;
  variables?: VariableDef[];
  forbidVariables?: boolean;
  linePlaceholder?: string;
  forbidJinjaHighlight?: boolean;
  readOnly?: boolean;
  customExtensions?: Extension[];
  autoScrollToBottom?: boolean;
  isGoTemplate?: boolean;
  onChange?: (value: string) => void;
  onBlur?: () => void;
  onFocus?: () => void;
  children?: React.ReactNode;
}

export interface PromptBasicEditorRef {
  setEditorValue: (value?: string) => void;
  insertText?: (text: string) => void;
}

const extensions = [
  EditorView.theme({
    '.cm-gutters': {
      backgroundColor: 'transparent',
      borderRight: 'none',
    },
    '.cm-scroller': {
      paddingLeft: '10px',
      paddingRight: '6px !important',
    },
  }),
];

export const PromptBasicEditor = forwardRef<
  PromptBasicEditorRef,
  PromptBasicEditorProps
>(
  (
    {
      defaultValue,
      onChange,
      variables,
      height,
      minHeight,
      maxHeight,
      fontSize = 13,
      forbidJinjaHighlight,
      forbidVariables,
      readOnly,
      linePlaceholder = I18n.t('please_input_with_vars'),
      customExtensions,
      autoScrollToBottom,
      onBlur,
      isGoTemplate,
      onFocus,
      children,
    }: PromptBasicEditorProps,
    ref,
  ) => {
    const editorRef = useRef<any>(null);

    useImperativeHandle(ref, () => ({
      setEditorValue: (value?: string) => {
        const editor = editorRef.current;
        if (!editor) {
          return;
        }
        editor?.setValue?.(value);
      },
      insertText: (text: string) => {
        const editor = editorRef.current;
        if (!editor) {
          return;
        }
        const range = editor.getSelection();
        if (!range) {
          return;
        }
        editor.replaceText({
          ...range,
          text,
          cursorOffset: 0,
        });
      },
    }));

    const newExtensions = useMemo(() => {
      const xExtensions = customExtensions || extensions;
      if (isGoTemplate) {
        return [...xExtensions, goExtension];
      }
      return xExtensions;
    }, [customExtensions, extensions, isGoTemplate]);

    return (
      <EditorProvider>
        <Renderer
          plugins={preset}
          defaultValue={defaultValue}
          options={{
            editable: !readOnly,
            readOnly,
            height,
            minHeight: minHeight || height,
            maxHeight: maxHeight || height,
            fontSize,
          }}
          onChange={e => onChange?.(e.value)}
          onFocus={onFocus}
          onBlur={onBlur}
          extensions={newExtensions}
          didMount={editor => {
            editorRef.current = editor;
            if (autoScrollToBottom) {
              editor.$view.dispatch({
                effects: EditorView.scrollIntoView(
                  editor.$view.state.doc.length,
                ),
              });
            }
          }}
        />

        {/* 输入 { 唤起变量选择 */}
        {!forbidVariables && <Variable variables={variables || []} />}

        <LanguageSupport />
        {/* Jinja 语法高亮 */}
        {!forbidJinjaHighlight && (
          <>
            <Validation />
            <JinjaHighlight />
          </>
        )}

        {/* Markdown 语法高亮 */}
        <MarkdownHighlight />

        {/* 激活行为空时的占位提示 */}

        <Placeholder>{linePlaceholder}</Placeholder>
        {children}
      </EditorProvider>
    );
  },
);
