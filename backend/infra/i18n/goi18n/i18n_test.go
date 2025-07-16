// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package goi18n

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeTestLangFile(t *testing.T, dir, lang, content string) {
	t.Helper()
	file := filepath.Join(dir, lang+".yaml")
	err := os.WriteFile(file, []byte(content), 0o644)
	require.NoError(t, err)
}

func TestTranslater(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestLangFile(t, tmpDir, "en-US", `
- id: hello
  translation: "Hello"
`)
	writeTestLangFile(t, tmpDir, "zh-CN", `
- id: hello
  translation: "你好"
`)

	trans, err := NewTranslater(tmpDir)
	require.NoError(t, err)
	require.NotNil(t, trans)

	ctx := context.Background()

	msg, err := trans.Translate(ctx, "hello", "en-US")
	require.NoError(t, err)
	require.Equal(t, "Hello", msg)

	msg, err = trans.Translate(ctx, "hello", "zh-CN")
	require.NoError(t, err)
	require.Equal(t, "你好", msg)

	msg, err = trans.Translate(ctx, "not-exist", "en-US")
	require.Error(t, err)
	require.Equal(t, "", msg)

	msg, err = trans.Translate(ctx, "hello", "ja-JP")
	require.Error(t, err)
	require.Equal(t, "", msg)

	require.Equal(t, "Hello", trans.MustTranslate(ctx, "hello", "en-US"))
	require.Equal(t, "你好", trans.MustTranslate(ctx, "hello", "zh-CN"))

	require.Equal(t, "", trans.MustTranslate(ctx, "not-exist", "en-US"))
	require.Equal(t, "", trans.MustTranslate(ctx, "hello", "ja-JP"))
}
