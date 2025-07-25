// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	_ "embed"

	"github.com/pkg/errors"
)

//go:embed embed/message.json
var schemaMessage string

var builtinSchemas = map[SchemaKey]*JSONSchema{}

func init() {
	if err := initSchemas(); err != nil {
		panic(err)
	}
}

func initSchemas() error {
	for key, raw := range map[SchemaKey]string{
		SchemaKeyString:  `{"type": "string"}`,
		SchemaKeyInteger: `{"type": "integer"}`,
		SchemaKeyFloat:   `{"type": "number"}`,
		SchemaKeyBool:    `{"type": "boolean"}`,
		SchemaKeyMessage: schemaMessage,
	} {
		schema, err := NewJSONSchema(raw)
		if err != nil {
			return errors.Wrapf(err, "loading schema of %s", key)
		}
		builtinSchemas[key] = schema
	}
	return nil
}
