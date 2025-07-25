// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"testing"
	"time"
)

func TestItemSnapshot_GetID(t *testing.T) {
	type fields struct {
		ID        int64
		VersionID int64
		Snapshot  *Item
		CreatedAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "测试 GetID 返回正确的 ID",
			fields: fields{
				ID: 123,
			},
			want: 123,
		},
		{
			name: "测试 GetID 返回 0 当 ID 为 0 时",
			fields: fields{
				ID: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ItemSnapshot{
				ID:        tt.fields.ID,
				VersionID: tt.fields.VersionID,
				Snapshot:  tt.fields.Snapshot,
				CreatedAt: tt.fields.CreatedAt,
			}
			if got := i.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemSnapshot_SetID(t *testing.T) {
	type fields struct {
		ID        int64
		VersionID int64
		Snapshot  *Item
		CreatedAt time.Time
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "测试 SetID 能正确设置 ID",
			fields: fields{
				ID: 0,
			},
			args: args{
				id: 456,
			},
		},
		{
			name: "测试 SetID 能覆盖原有 ID",
			fields: fields{
				ID: 123,
			},
			args: args{
				id: 789,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ItemSnapshot{
				ID:        tt.fields.ID,
				VersionID: tt.fields.VersionID,
				Snapshot:  tt.fields.Snapshot,
				CreatedAt: tt.fields.CreatedAt,
			}
			i.SetID(tt.args.id)
			if i.ID != tt.args.id {
				t.Errorf("SetID() 未正确设置 ID，期望 %v，实际 %v", tt.args.id, i.ID)
			}
		})
	}
}
