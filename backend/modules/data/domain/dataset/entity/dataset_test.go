// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"testing"
	"time"

	"github.com/bytedance/gg/gptr"
)

func TestContentType_IsMultiModal(t *testing.T) {
	tests := []struct {
		name string
		ct   ContentType
		want bool
	}{
		{ct: ContentTypeText, want: false},
		{ct: ContentTypeImage, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.IsMultiModal(); got != tt.want {
				t.Errorf("IsMultiModal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatasetOperation_String(t *testing.T) {
	type fields struct {
		ID   string
		Type DatasetOpType
		TS   time.Time
		TTL  time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{ID: "123"},
			want:   "{id=123, type=, ts=0001-01-01 00:00:00 +0000 UTC, ttl=0s}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &DatasetOperation{
				ID:   tt.fields.ID,
				Type: tt.fields.Type,
				TS:   tt.fields.TS,
				TTL:  tt.fields.TTL,
			}
			if got := o.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataset_CanWriteItem(t *testing.T) {
	type fields struct {
		ID             int64
		AppID          int32
		SpaceID        int64
		SchemaID       int64
		Name           string
		Description    *string
		Category       DatasetCategory
		BizCategory    string
		Status         DatasetStatus
		SecurityLevel  SecurityLevel
		Visibility     DatasetVisibility
		Spec           *DatasetSpec
		Features       *DatasetFeatures
		LatestVersion  string
		NextVersionNum int64
		LastOperation  DatasetOpType
		CreatedBy      string
		CreatedAt      time.Time
		UpdatedBy      string
		UpdatedAt      time.Time
		ExpiredAt      *time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{fields: fields{Status: DatasetStatusAvailable}, want: true},
		{fields: fields{Status: DatasetStatusDeleted}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dataset{
				ID:             tt.fields.ID,
				AppID:          tt.fields.AppID,
				SpaceID:        tt.fields.SpaceID,
				SchemaID:       tt.fields.SchemaID,
				Name:           tt.fields.Name,
				Description:    tt.fields.Description,
				Category:       tt.fields.Category,
				BizCategory:    tt.fields.BizCategory,
				Status:         tt.fields.Status,
				SecurityLevel:  tt.fields.SecurityLevel,
				Visibility:     tt.fields.Visibility,
				Spec:           tt.fields.Spec,
				Features:       tt.fields.Features,
				LatestVersion:  tt.fields.LatestVersion,
				NextVersionNum: tt.fields.NextVersionNum,
				LastOperation:  tt.fields.LastOperation,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				UpdatedBy:      tt.fields.UpdatedBy,
				UpdatedAt:      tt.fields.UpdatedAt,
				ExpiredAt:      tt.fields.ExpiredAt,
			}
			if got := d.CanWriteItem(); got != tt.want {
				t.Errorf("CanWriteItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataset_GetDescription(t *testing.T) {
	type fields struct {
		ID             int64
		AppID          int32
		SpaceID        int64
		SchemaID       int64
		Name           string
		Description    *string
		Category       DatasetCategory
		BizCategory    string
		Status         DatasetStatus
		SecurityLevel  SecurityLevel
		Visibility     DatasetVisibility
		Spec           *DatasetSpec
		Features       *DatasetFeatures
		LatestVersion  string
		NextVersionNum int64
		LastOperation  DatasetOpType
		CreatedBy      string
		CreatedAt      time.Time
		UpdatedBy      string
		UpdatedAt      time.Time
		ExpiredAt      *time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{Description: gptr.Of("test")},
			want:   "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dataset{
				ID:             tt.fields.ID,
				AppID:          tt.fields.AppID,
				SpaceID:        tt.fields.SpaceID,
				SchemaID:       tt.fields.SchemaID,
				Name:           tt.fields.Name,
				Description:    tt.fields.Description,
				Category:       tt.fields.Category,
				BizCategory:    tt.fields.BizCategory,
				Status:         tt.fields.Status,
				SecurityLevel:  tt.fields.SecurityLevel,
				Visibility:     tt.fields.Visibility,
				Spec:           tt.fields.Spec,
				Features:       tt.fields.Features,
				LatestVersion:  tt.fields.LatestVersion,
				NextVersionNum: tt.fields.NextVersionNum,
				LastOperation:  tt.fields.LastOperation,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				UpdatedBy:      tt.fields.UpdatedBy,
				UpdatedAt:      tt.fields.UpdatedAt,
				ExpiredAt:      tt.fields.ExpiredAt,
			}
			if got := d.GetDescription(); got != tt.want {
				t.Errorf("GetDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataset_GetID(t *testing.T) {
	type fields struct {
		ID             int64
		AppID          int32
		SpaceID        int64
		SchemaID       int64
		Name           string
		Description    *string
		Category       DatasetCategory
		BizCategory    string
		Status         DatasetStatus
		SecurityLevel  SecurityLevel
		Visibility     DatasetVisibility
		Spec           *DatasetSpec
		Features       *DatasetFeatures
		LatestVersion  string
		NextVersionNum int64
		LastOperation  DatasetOpType
		CreatedBy      string
		CreatedAt      time.Time
		UpdatedBy      string
		UpdatedAt      time.Time
		ExpiredAt      *time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{fields: fields{ID: 123}, want: 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dataset{
				ID:             tt.fields.ID,
				AppID:          tt.fields.AppID,
				SpaceID:        tt.fields.SpaceID,
				SchemaID:       tt.fields.SchemaID,
				Name:           tt.fields.Name,
				Description:    tt.fields.Description,
				Category:       tt.fields.Category,
				BizCategory:    tt.fields.BizCategory,
				Status:         tt.fields.Status,
				SecurityLevel:  tt.fields.SecurityLevel,
				Visibility:     tt.fields.Visibility,
				Spec:           tt.fields.Spec,
				Features:       tt.fields.Features,
				LatestVersion:  tt.fields.LatestVersion,
				NextVersionNum: tt.fields.NextVersionNum,
				LastOperation:  tt.fields.LastOperation,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				UpdatedBy:      tt.fields.UpdatedBy,
				UpdatedAt:      tt.fields.UpdatedAt,
				ExpiredAt:      tt.fields.ExpiredAt,
			}
			if got := d.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataset_IsChangeUncommitted(t *testing.T) {
	type fields struct {
		ID             int64
		AppID          int32
		SpaceID        int64
		SchemaID       int64
		Name           string
		Description    *string
		Category       DatasetCategory
		BizCategory    string
		Status         DatasetStatus
		SecurityLevel  SecurityLevel
		Visibility     DatasetVisibility
		Spec           *DatasetSpec
		Features       *DatasetFeatures
		LatestVersion  string
		NextVersionNum int64
		LastOperation  DatasetOpType
		CreatedBy      string
		CreatedAt      time.Time
		UpdatedBy      string
		UpdatedAt      time.Time
		ExpiredAt      *time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{fields: fields{LastOperation: DatasetOpTypeCreateVersion}, want: false},
		{fields: fields{LastOperation: "other"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dataset{
				ID:             tt.fields.ID,
				AppID:          tt.fields.AppID,
				SpaceID:        tt.fields.SpaceID,
				SchemaID:       tt.fields.SchemaID,
				Name:           tt.fields.Name,
				Description:    tt.fields.Description,
				Category:       tt.fields.Category,
				BizCategory:    tt.fields.BizCategory,
				Status:         tt.fields.Status,
				SecurityLevel:  tt.fields.SecurityLevel,
				Visibility:     tt.fields.Visibility,
				Spec:           tt.fields.Spec,
				Features:       tt.fields.Features,
				LatestVersion:  tt.fields.LatestVersion,
				NextVersionNum: tt.fields.NextVersionNum,
				LastOperation:  tt.fields.LastOperation,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				UpdatedBy:      tt.fields.UpdatedBy,
				UpdatedAt:      tt.fields.UpdatedAt,
				ExpiredAt:      tt.fields.ExpiredAt,
			}
			if got := d.IsChangeUncommitted(); got != tt.want {
				t.Errorf("IsChangeUncommitted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataset_SetID(t *testing.T) {
	type fields struct {
		ID             int64
		AppID          int32
		SpaceID        int64
		SchemaID       int64
		Name           string
		Description    *string
		Category       DatasetCategory
		BizCategory    string
		Status         DatasetStatus
		SecurityLevel  SecurityLevel
		Visibility     DatasetVisibility
		Spec           *DatasetSpec
		Features       *DatasetFeatures
		LatestVersion  string
		NextVersionNum int64
		LastOperation  DatasetOpType
		CreatedBy      string
		CreatedAt      time.Time
		UpdatedBy      string
		UpdatedAt      time.Time
		ExpiredAt      *time.Time
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{args: args{id: 123}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dataset{
				ID:             tt.fields.ID,
				AppID:          tt.fields.AppID,
				SpaceID:        tt.fields.SpaceID,
				SchemaID:       tt.fields.SchemaID,
				Name:           tt.fields.Name,
				Description:    tt.fields.Description,
				Category:       tt.fields.Category,
				BizCategory:    tt.fields.BizCategory,
				Status:         tt.fields.Status,
				SecurityLevel:  tt.fields.SecurityLevel,
				Visibility:     tt.fields.Visibility,
				Spec:           tt.fields.Spec,
				Features:       tt.fields.Features,
				LatestVersion:  tt.fields.LatestVersion,
				NextVersionNum: tt.fields.NextVersionNum,
				LastOperation:  tt.fields.LastOperation,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				UpdatedBy:      tt.fields.UpdatedBy,
				UpdatedAt:      tt.fields.UpdatedAt,
				ExpiredAt:      tt.fields.ExpiredAt,
			}
			d.SetID(tt.args.id)
		})
	}
}
