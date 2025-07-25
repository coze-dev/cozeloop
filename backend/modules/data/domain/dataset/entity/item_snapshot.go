// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import "time"

type ItemSnapshot struct {
	ID        int64
	VersionID int64

	Snapshot *Item

	CreatedAt time.Time
}

func (i *ItemSnapshot) SetID(id int64) { i.ID = id }

func (i *ItemSnapshot) GetID() int64 { return i.ID }

type SnapshotProgress struct {
	Cursor string
}
