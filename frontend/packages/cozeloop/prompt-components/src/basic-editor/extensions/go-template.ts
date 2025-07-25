// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { createOnigurumaEngine } from 'shiki/engine/oniguruma';
import { createHighlighterCore } from 'shiki/core';
import shiki from 'codemirror-shiki';

const highlighter = createHighlighterCore({
  langs: [import('./go-syntax').then(mod => mod.tmLanguage)],
  themes: [import('@shikijs/themes/one-light')],
  engine: createOnigurumaEngine(import('shiki/wasm')),
});

export const goExtension = shiki({
  highlighter,
  language: 'go-template',
  theme: 'one-light',
});
