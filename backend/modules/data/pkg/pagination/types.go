// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package pagination

import (
	"encoding/base64"
	"time"

	"github.com/bytedance/sonic"
)

type Cursor struct {
	TimeValue *time.Time `json:"t,omitempty"` // updated_at or created_at of last page
	IDValue   *int64     `json:"id,omitempty"`
}

func (c *Cursor) Encode() (string, error) {
	if c.TimeValue == nil && c.IDValue == nil {
		return "", nil
	}

	b, err := sonic.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func decodeCursor(c string) (*Cursor, error) {
	b, err := base64.StdEncoding.DecodeString(c)
	if err != nil {
		return nil, err
	}
	var cursor *Cursor
	if err := sonic.Unmarshal(b, &cursor); err != nil {
		return nil, err
	}
	return cursor, nil
}

type PageResult struct {
	Total  int64
	Cursor string
}

type PreviousPage interface {
	GetPage() int32
	GetPageSize() int32
	GetCursor() string
}
