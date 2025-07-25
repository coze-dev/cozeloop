// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"reflect"
	"testing"
	"time"

	"github.com/bytedance/gg/gptr"
)

func TestFieldData_DataBytes(t *testing.T) {
	type fields struct {
		Key         string
		Name        string
		ContentType ContentType
		Format      FieldDisplayFormat
		Content     string
		Attachments []*ObjectStorage
		Parts       []*FieldData
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Test empty content",
			fields: fields{
				Content: "",
			},
			want: 0,
		},
		{
			name: "Test non-empty content",
			fields: fields{
				Content: "test",
			},
			want: 4,
		},
		{
			name: "Test attachments and parts",
			fields: fields{
				Content: "test",
				Attachments: []*ObjectStorage{
					{
						URL: "test",
					},
				},
				Parts: []*FieldData{
					{
						Name: "test",
					},
				},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FieldData{
				Key:         tt.fields.Key,
				Name:        tt.fields.Name,
				ContentType: tt.fields.ContentType,
				Format:      tt.fields.Format,
				Content:     tt.fields.Content,
				Attachments: tt.fields.Attachments,
				Parts:       tt.fields.Parts,
			}
			if got := f.DataBytes(); got != tt.want {
				t.Errorf("DataBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldData_DataRunes(t *testing.T) {
	type fields struct {
		Key         string
		Name        string
		ContentType ContentType
		Format      FieldDisplayFormat
		Content     string
		Attachments []*ObjectStorage
		Parts       []*FieldData
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Test empty content",
			fields: fields{
				Content: "",
			},
			want: 0,
		},
		{
			name: "Test non-empty content",
			fields: fields{
				Content: "test",
			},
			want: 4,
		},
		{
			name: "Test attachments and parts",
			fields: fields{
				Content: "test",
				Attachments: []*ObjectStorage{
					{
						URL: "test",
					},
				},
				Parts: []*FieldData{
					{
						Name: "test",
					},
				},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FieldData{
				Key:         tt.fields.Key,
				Name:        tt.fields.Name,
				ContentType: tt.fields.ContentType,
				Format:      tt.fields.Format,
				Content:     tt.fields.Content,
				Attachments: tt.fields.Attachments,
				Parts:       tt.fields.Parts,
			}
			if got := f.DataRunes(); got != tt.want {
				t.Errorf("DataRunes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemErrorDetailToString(t *testing.T) {
	type args struct {
		d *ItemErrorDetail
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test nil input",
			args: args{
				d: &ItemErrorDetail{Message: gptr.Of("test")},
			},
			want: "test, index=0",
		},
		// 实际需要根据 ItemErrorDetailToString 实现添加更多测试用例
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ItemErrorDetailToString(tt.args.d); got != tt.want {
				t.Errorf("ItemErrorDetailToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItem_AllData(t *testing.T) {
	type fields struct {
		ID             int64
		AppID          int32
		SpaceID        int64
		DatasetID      int64
		SchemaID       int64
		ItemID         int64
		ItemKey        string
		Data           []*FieldData
		RepeatedData   []*ItemData
		DataProperties *ItemDataProperties
		AddVN          int64
		DelVN          int64
		CreatedBy      string
		CreatedAt      time.Time
		UpdatedBy      string
		UpdatedAt      time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   [][]*FieldData
	}{
		{
			name: "Test empty data",
			fields: fields{
				Data:         nil,
				RepeatedData: nil,
			},
			want: [][]*FieldData{nil},
		},
		// 实际需要根据 Item.AllData 实现添加更多测试用例
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Item{
				ID:             tt.fields.ID,
				AppID:          tt.fields.AppID,
				SpaceID:        tt.fields.SpaceID,
				DatasetID:      tt.fields.DatasetID,
				SchemaID:       tt.fields.SchemaID,
				ItemID:         tt.fields.ItemID,
				ItemKey:        tt.fields.ItemKey,
				Data:           tt.fields.Data,
				RepeatedData:   tt.fields.RepeatedData,
				DataProperties: tt.fields.DataProperties,
				AddVN:          tt.fields.AddVN,
				DelVN:          tt.fields.DelVN,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				UpdatedBy:      tt.fields.UpdatedBy,
				UpdatedAt:      tt.fields.UpdatedAt,
			}
			if got := i.AllData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItem_BuildProperties(t *testing.T) {
	type fields struct {
		ID             int64
		AppID          int32
		SpaceID        int64
		DatasetID      int64
		SchemaID       int64
		ItemID         int64
		ItemKey        string
		Data           []*FieldData
		RepeatedData   []*ItemData
		DataProperties *ItemDataProperties
		AddVN          int64
		DelVN          int64
		CreatedBy      string
		CreatedAt      time.Time
		UpdatedBy      string
		UpdatedAt      time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "Test BuildProperties",
			fields: fields{
				// 初始化字段
			},
		},
		// 实际需要根据 Item.BuildProperties 实现添加更多测试用例
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Item{
				ID:             tt.fields.ID,
				AppID:          tt.fields.AppID,
				SpaceID:        tt.fields.SpaceID,
				DatasetID:      tt.fields.DatasetID,
				SchemaID:       tt.fields.SchemaID,
				ItemID:         tt.fields.ItemID,
				ItemKey:        tt.fields.ItemKey,
				Data:           tt.fields.Data,
				RepeatedData:   tt.fields.RepeatedData,
				DataProperties: tt.fields.DataProperties,
				AddVN:          tt.fields.AddVN,
				DelVN:          tt.fields.DelVN,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				UpdatedBy:      tt.fields.UpdatedBy,
				UpdatedAt:      tt.fields.UpdatedAt,
			}
			i.BuildProperties()
			// 可以添加更多断言来验证 BuildProperties 的效果
		})
	}
}

func TestItem_ClearData(t *testing.T) {
	type fields struct {
		ID             int64
		AppID          int32
		SpaceID        int64
		DatasetID      int64
		SchemaID       int64
		ItemID         int64
		ItemKey        string
		Data           []*FieldData
		RepeatedData   []*ItemData
		DataProperties *ItemDataProperties
		AddVN          int64
		DelVN          int64
		CreatedBy      string
		CreatedAt      time.Time
		UpdatedBy      string
		UpdatedAt      time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test ClearData",
			fields: fields{
				Data:         []*FieldData{{}},
				RepeatedData: []*ItemData{{}},
			},
		},
		// 实际需要根据 Item.ClearData 实现添加更多测试用例
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Item{
				ID:             tt.fields.ID,
				AppID:          tt.fields.AppID,
				SpaceID:        tt.fields.SpaceID,
				DatasetID:      tt.fields.DatasetID,
				SchemaID:       tt.fields.SchemaID,
				ItemID:         tt.fields.ItemID,
				ItemKey:        tt.fields.ItemKey,
				Data:           tt.fields.Data,
				RepeatedData:   tt.fields.RepeatedData,
				DataProperties: tt.fields.DataProperties,
				AddVN:          tt.fields.AddVN,
				DelVN:          tt.fields.DelVN,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				UpdatedBy:      tt.fields.UpdatedBy,
				UpdatedAt:      tt.fields.UpdatedAt,
			}
			i.ClearData()
			// 可以添加更多断言来验证 ClearData 的效果
		})
	}
}

func TestItem_GetID(t *testing.T) {
	type fields struct {
		ID             int64
		AppID          int32
		SpaceID        int64
		DatasetID      int64
		SchemaID       int64
		ItemID         int64
		ItemKey        string
		Data           []*FieldData
		RepeatedData   []*ItemData
		DataProperties *ItemDataProperties
		AddVN          int64
		DelVN          int64
		CreatedBy      string
		CreatedAt      time.Time
		UpdatedBy      string
		UpdatedAt      time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "Test GetID",
			fields: fields{
				ID: 123,
			},
			want: 123,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Item{
				ID:             tt.fields.ID,
				AppID:          tt.fields.AppID,
				SpaceID:        tt.fields.SpaceID,
				DatasetID:      tt.fields.DatasetID,
				SchemaID:       tt.fields.SchemaID,
				ItemID:         tt.fields.ItemID,
				ItemKey:        tt.fields.ItemKey,
				Data:           tt.fields.Data,
				RepeatedData:   tt.fields.RepeatedData,
				DataProperties: tt.fields.DataProperties,
				AddVN:          tt.fields.AddVN,
				DelVN:          tt.fields.DelVN,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				UpdatedBy:      tt.fields.UpdatedBy,
				UpdatedAt:      tt.fields.UpdatedAt,
			}
			if got := i.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItem_GetOrBuildProperties(t *testing.T) {
	type fields struct {
		ID             int64
		AppID          int32
		SpaceID        int64
		DatasetID      int64
		SchemaID       int64
		ItemID         int64
		ItemKey        string
		Data           []*FieldData
		RepeatedData   []*ItemData
		DataProperties *ItemDataProperties
		AddVN          int64
		DelVN          int64
		CreatedBy      string
		CreatedAt      time.Time
		UpdatedBy      string
		UpdatedAt      time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   *ItemDataProperties
	}{
		{
			name: "Test GetOrBuildProperties with existing properties",
			fields: fields{
				DataProperties: &ItemDataProperties{},
			},
			want: &ItemDataProperties{},
		},
		// 实际需要根据 Item.GetOrBuildProperties 实现添加更多测试用例
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Item{
				ID:             tt.fields.ID,
				AppID:          tt.fields.AppID,
				SpaceID:        tt.fields.SpaceID,
				DatasetID:      tt.fields.DatasetID,
				SchemaID:       tt.fields.SchemaID,
				ItemID:         tt.fields.ItemID,
				ItemKey:        tt.fields.ItemKey,
				Data:           tt.fields.Data,
				RepeatedData:   tt.fields.RepeatedData,
				DataProperties: tt.fields.DataProperties,
				AddVN:          tt.fields.AddVN,
				DelVN:          tt.fields.DelVN,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				UpdatedBy:      tt.fields.UpdatedBy,
				UpdatedAt:      tt.fields.UpdatedAt,
			}
			if got := i.GetOrBuildProperties(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrBuildProperties() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItem_SetID(t *testing.T) {
	type fields struct {
		ID             int64
		AppID          int32
		SpaceID        int64
		DatasetID      int64
		SchemaID       int64
		ItemID         int64
		ItemKey        string
		Data           []*FieldData
		RepeatedData   []*ItemData
		DataProperties *ItemDataProperties
		AddVN          int64
		DelVN          int64
		CreatedBy      string
		CreatedAt      time.Time
		UpdatedBy      string
		UpdatedAt      time.Time
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
			name: "Test SetID",
			fields: fields{
				ID: 0,
			},
			args: args{
				id: 123,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Item{
				ID:             tt.fields.ID,
				AppID:          tt.fields.AppID,
				SpaceID:        tt.fields.SpaceID,
				DatasetID:      tt.fields.DatasetID,
				SchemaID:       tt.fields.SchemaID,
				ItemID:         tt.fields.ItemID,
				ItemKey:        tt.fields.ItemKey,
				Data:           tt.fields.Data,
				RepeatedData:   tt.fields.RepeatedData,
				DataProperties: tt.fields.DataProperties,
				AddVN:          tt.fields.AddVN,
				DelVN:          tt.fields.DelVN,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				UpdatedBy:      tt.fields.UpdatedBy,
				UpdatedAt:      tt.fields.UpdatedAt,
			}
			i.SetID(tt.args.id)
			if i.ID != tt.args.id {
				t.Errorf("SetID() failed, expected ID: %v, got: %v", tt.args.id, i.ID)
			}
		})
	}
}

func TestSanitizeItemErrorGroup(t *testing.T) {
	type args struct {
		eg           *ItemErrorGroup
		maxDetailCnt int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test SanitizeItemErrorGroup with nil input",
			args: args{
				eg: &ItemErrorGroup{
					Details: []*ItemErrorDetail{
						{Message: gptr.Of("test")},
					},
				},
				maxDetailCnt: 0,
			},
		},
		// 实际需要根据 SanitizeItemErrorGroup 实现添加更多测试用例
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SanitizeItemErrorGroup(tt.args.eg, tt.args.maxDetailCnt)
			// 可以添加更多断言来验证 SanitizeItemErrorGroup 的效果
		})
	}
}
