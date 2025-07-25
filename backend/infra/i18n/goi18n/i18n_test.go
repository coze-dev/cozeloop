// Copyright (c) 2025 coze-dev Authors
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

func TestTranslater_TableDriven(t *testing.T) {
	tmpDir := t.TempDir()
	writeTestLangFile(t, tmpDir, "en-US", `
- id: hello
  translation: "Hello"
`)
	writeTestLangFile(t, tmpDir, "zh-CN", `
- id: hello
  translation: "你好"
`)

	testCases := []struct {
		name        string
		prepare     func(dir string)
		langDir     string
		key         string
		lang        string
		expectMsg   string
		expectErr   bool
		newTransErr bool
	}{
		{
			name:      "正常英文",
			langDir:   tmpDir,
			key:       "hello",
			lang:      "en-US",
			expectMsg: "Hello",
			prepare:   nil,
		},
		{
			name:      "正常中文",
			langDir:   tmpDir,
			key:       "hello",
			lang:      "zh-CN",
			expectMsg: "你好",
			prepare:   nil,
		},
		{
			name:      "key 不存在",
			langDir:   tmpDir,
			key:       "not-exist",
			lang:      "en-US",
			expectMsg: "",
			expectErr: true,
			prepare:   nil,
		},
		{
			name:      "语言不支持",
			langDir:   tmpDir,
			key:       "hello",
			lang:      "ja-JP",
			expectMsg: "",
			expectErr: true,
			prepare:   nil,
		},
		{
			name:      "非法语言格式",
			langDir:   tmpDir,
			key:       "hello",
			lang:      "not-a-lang",
			expectMsg: "",
			expectErr: true,
			prepare:   nil,
		},
		{
			name:        "NewTranslater: 目录不存在",
			langDir:     filepath.Join(tmpDir, "not-exist-dir"),
			key:         "hello",
			lang:        "en-US",
			expectMsg:   "",
			newTransErr: true,
			prepare:     nil,
		},
		{
			name:    "NewTranslater: 语言文件名非法，Parse 失败",
			langDir: tmpDir,
			key:     "hello",
			lang:    "en-US",
			prepare: func(dir string) {
				writeTestLangFile(t, dir, "badlang", `- id: hello\n  translation: 'bad'`)
			},
			expectMsg: "Hello",
		},
		{
			name:    "NewTranslater: 加载文件失败",
			langDir: tmpDir,
			key:     "hello",
			lang:    "en-US",
			prepare: func(dir string) {
				// 写一个非法 yaml 文件，导致 LoadMessageFile 失败
				writeTestLangFile(t, dir, "fr-FR", `not a yaml`)
			},
			newTransErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.prepare != nil {
				tc.prepare(tc.langDir)
			}
			trans, err := NewTranslater(tc.langDir)
			if tc.newTransErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, trans)

			ctx := context.Background()
			msg, err := trans.Translate(ctx, tc.key, tc.lang)
			if tc.expectErr {
				require.Error(t, err)
				require.Equal(t, "", msg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectMsg, msg)
			}
		})
	}
}

func TestTranslater_MustTranslate(t *testing.T) {
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
	ctx := context.Background()
	require.Equal(t, "Hello", trans.MustTranslate(ctx, "hello", "en-US"))
	require.Equal(t, "你好", trans.MustTranslate(ctx, "hello", "zh-CN"))
	require.Equal(t, "", trans.MustTranslate(ctx, "not-exist", "en-US"))
	require.Equal(t, "", trans.MustTranslate(ctx, "hello", "ja-JP"))
}
