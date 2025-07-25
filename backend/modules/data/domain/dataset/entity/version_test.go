// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"testing"
	"time"
)

func TestDatasetVersion_GetID(t *testing.T) {
	type fields struct {
		ID               int64
		AppID            int32
		SpaceID          int64
		DatasetID        int64
		SchemaID         int64
		DatasetBrief     *Dataset
		Version          string
		VersionNum       int64
		Description      *string
		ItemCount        int64
		SnapshotStatus   SnapshotStatus
		SnapshotProgress *SnapshotProgress
		UpdateVersion    int64
		CreatedBy        string
		CreatedAt        time.Time
		DisabledAt       *time.Time
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
			s := &DatasetVersion{
				ID:               tt.fields.ID,
				AppID:            tt.fields.AppID,
				SpaceID:          tt.fields.SpaceID,
				DatasetID:        tt.fields.DatasetID,
				SchemaID:         tt.fields.SchemaID,
				DatasetBrief:     tt.fields.DatasetBrief,
				Version:          tt.fields.Version,
				VersionNum:       tt.fields.VersionNum,
				Description:      tt.fields.Description,
				ItemCount:        tt.fields.ItemCount,
				SnapshotStatus:   tt.fields.SnapshotStatus,
				SnapshotProgress: tt.fields.SnapshotProgress,
				UpdateVersion:    tt.fields.UpdateVersion,
				CreatedBy:        tt.fields.CreatedBy,
				CreatedAt:        tt.fields.CreatedAt,
				DisabledAt:       tt.fields.DisabledAt,
			}
			if got := s.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatasetVersion_SetID(t *testing.T) {
	type fields struct {
		ID               int64
		AppID            int32
		SpaceID          int64
		DatasetID        int64
		SchemaID         int64
		DatasetBrief     *Dataset
		Version          string
		VersionNum       int64
		Description      *string
		ItemCount        int64
		SnapshotStatus   SnapshotStatus
		SnapshotProgress *SnapshotProgress
		UpdateVersion    int64
		CreatedBy        string
		CreatedAt        time.Time
		DisabledAt       *time.Time
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
			s := &DatasetVersion{
				ID:               tt.fields.ID,
				AppID:            tt.fields.AppID,
				SpaceID:          tt.fields.SpaceID,
				DatasetID:        tt.fields.DatasetID,
				SchemaID:         tt.fields.SchemaID,
				DatasetBrief:     tt.fields.DatasetBrief,
				Version:          tt.fields.Version,
				VersionNum:       tt.fields.VersionNum,
				Description:      tt.fields.Description,
				ItemCount:        tt.fields.ItemCount,
				SnapshotStatus:   tt.fields.SnapshotStatus,
				SnapshotProgress: tt.fields.SnapshotProgress,
				UpdateVersion:    tt.fields.UpdateVersion,
				CreatedBy:        tt.fields.CreatedBy,
				CreatedAt:        tt.fields.CreatedAt,
				DisabledAt:       tt.fields.DisabledAt,
			}
			s.SetID(tt.args.id)
			if s.ID != tt.args.id {
				t.Errorf("SetID() 未正确设置 ID，期望 %v，实际 %v", tt.args.id, s.ID)
			}
		})
	}
}

func TestSnapshotStatus_IsFinished(t *testing.T) {
	tests := []struct {
		name string
		ss   SnapshotStatus
		want bool
	}{
		{
			name: "测试未开始状态返回 false",
			ss:   SnapshotStatusUnstarted,
			want: false,
		},
		{
			name: "测试进行中状态返回 false",
			ss:   SnapshotStatusInProgress,
			want: false,
		},
		{
			name: "测试完成状态返回 true",
			ss:   SnapshotStatusCompleted,
			want: true,
		},
		{
			name: "测试失败状态返回 true",
			ss:   SnapshotStatusFailed,
			want: true,
		},
		{
			name: "测试未知状态返回 false",
			ss:   SnapshotStatusUnknown,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.IsFinished(); got != tt.want {
				t.Errorf("IsFinished() = %v, want %v", got, tt.want)
			}
		})
	}
}
