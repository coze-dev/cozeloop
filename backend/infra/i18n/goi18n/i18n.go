// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package goi18n

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"

	internationalization "github.com/coze-dev/cozeloop/backend/infra/i18n"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

func NewTranslater(langDir string) (internationalization.ITranslater, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	files, err := os.ReadDir(langDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read lang dir: %w", err)
	}

	localizers := make(map[language.Tag]*i18n.Localizer)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if !strings.HasSuffix(name, ".yaml") {
			continue
		}
		langStr := strings.TrimSuffix(name, ".yaml")
		langTag, err := language.Parse(langStr)
		if err != nil {
			continue
		}
		langFile := filepath.Join(langDir, name)
		_, err = bundle.LoadMessageFile(langFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load language file %s: %w", langFile, err)
		}
		localizers[langTag] = i18n.NewLocalizer(bundle, langTag.String())
	}

	return &translater{localizers: localizers}, nil
}

type translater struct {
	localizers map[language.Tag]*i18n.Localizer
}

func (t *translater) Translate(ctx context.Context, key string, lang string) (string, error) {
	langTag, err := language.Parse(lang)
	if err != nil {
		return "", fmt.Errorf("invalid language: %s", lang)
	}
	localizer, ok := t.localizers[langTag]
	if !ok {
		return "", fmt.Errorf("language %s not supported", lang)
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: key})
	if err != nil {
		return "", fmt.Errorf("i18n localize fail, key: %s, lang: %s, err: %w", key, lang, err)
	}
	return msg, nil
}

func (t *translater) MustTranslate(ctx context.Context, key string, lang string) string {
	msg, err := t.Translate(ctx, key, lang)
	if err != nil {
		logs.CtxWarn(ctx, "i18n translater fail, err: %v", err)
		return ""
	}
	return msg
}
