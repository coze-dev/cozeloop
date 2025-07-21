// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package i18n

import (
	"context"
)

//go:generate mockgen -destination=mocks/i18n.go -package=mocks . ITranslater
type ITranslater interface {
	Translate(ctx context.Context, key string, lang string) (string, error)
	MustTranslate(ctx context.Context, key string, lang string) string
}
