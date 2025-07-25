// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
)

func TestNewSchemaOfDataset(t *testing.T) {
	tests := []struct {
		name     string
		dataset  *entity.Dataset
		fields   []*entity.FieldSchema
		expected *entity.DatasetSchema
	}{
		{
			name: "Create schema with immutable=false when Features.EditSchema=true",
			dataset: &entity.Dataset{
				ID:        1,
				AppID:     100,
				SpaceID:   1000,
				SchemaID:  2000,
				CreatedBy: "user1",
				Features: &entity.DatasetFeatures{
					EditSchema: true,
				},
			},
			fields: []*entity.FieldSchema{
				{
					Key:         "field1",
					Name:        "Field 1",
					ContentType: entity.ContentTypeText,
				},
			},
			expected: &entity.DatasetSchema{
				ID:        2000,
				AppID:     100,
				SpaceID:   1000,
				DatasetID: 1,
				Fields: []*entity.FieldSchema{
					{
						Key:         "field1",
						Name:        "Field 1",
						ContentType: entity.ContentTypeText,
					},
				},
				Immutable: false,
				CreatedBy: "user1",
				UpdatedBy: "user1",
			},
		},
		{
			name: "Create schema with immutable=true when Features.EditSchema=false",
			dataset: &entity.Dataset{
				ID:        2,
				AppID:     200,
				SpaceID:   2000,
				SchemaID:  3000,
				CreatedBy: "user2",
				Features: &entity.DatasetFeatures{
					EditSchema: false,
				},
			},
			fields: []*entity.FieldSchema{
				{
					Key:         "field1",
					Name:        "Field 1",
					ContentType: entity.ContentTypeText,
				},
			},
			expected: &entity.DatasetSchema{
				ID:        3000,
				AppID:     200,
				SpaceID:   2000,
				DatasetID: 2,
				Fields: []*entity.FieldSchema{
					{
						Key:         "field1",
						Name:        "Field 1",
						ContentType: entity.ContentTypeText,
					},
				},
				Immutable: true,
				CreatedBy: "user2",
				UpdatedBy: "user2",
			},
		},
		{
			name: "Create schema with immutable=false when Features is nil",
			dataset: &entity.Dataset{
				ID:        3,
				AppID:     300,
				SpaceID:   3000,
				SchemaID:  4000,
				CreatedBy: "user3",
				Features:  nil,
			},
			fields: []*entity.FieldSchema{
				{
					Key:         "field1",
					Name:        "Field 1",
					ContentType: entity.ContentTypeText,
				},
			},
			expected: &entity.DatasetSchema{
				ID:        4000,
				AppID:     300,
				SpaceID:   3000,
				DatasetID: 3,
				Fields: []*entity.FieldSchema{
					{
						Key:         "field1",
						Name:        "Field 1",
						ContentType: entity.ContentTypeText,
					},
				},
				Immutable: false,
				CreatedBy: "user3",
				UpdatedBy: "user3",
			},
		},
		{
			name: "Create schema with multiple fields",
			dataset: &entity.Dataset{
				ID:        4,
				AppID:     400,
				SpaceID:   4000,
				SchemaID:  5000,
				CreatedBy: "user4",
				Features: &entity.DatasetFeatures{
					EditSchema: true,
				},
			},
			fields: []*entity.FieldSchema{
				{
					Key:         "field1",
					Name:        "Field 1",
					ContentType: entity.ContentTypeText,
				},
				{
					Key:         "field2",
					Name:        "Field 2",
					ContentType: entity.ContentTypeImage,
				},
			},
			expected: &entity.DatasetSchema{
				ID:        5000,
				AppID:     400,
				SpaceID:   4000,
				DatasetID: 4,
				Fields: []*entity.FieldSchema{
					{
						Key:         "field1",
						Name:        "Field 1",
						ContentType: entity.ContentTypeText,
					},
					{
						Key:         "field2",
						Name:        "Field 2",
						ContentType: entity.ContentTypeImage,
					},
				},
				Immutable: false,
				CreatedBy: "user4",
				UpdatedBy: "user4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newSchemaOfDataset(tt.dataset, tt.fields)

			// Verify all fields match expected values
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.AppID, result.AppID)
			assert.Equal(t, tt.expected.SpaceID, result.SpaceID)
			assert.Equal(t, tt.expected.DatasetID, result.DatasetID)
			assert.Equal(t, tt.expected.Immutable, result.Immutable)
			assert.Equal(t, tt.expected.CreatedBy, result.CreatedBy)
			assert.Equal(t, tt.expected.UpdatedBy, result.UpdatedBy)

			// Verify fields
			assert.Equal(t, len(tt.expected.Fields), len(result.Fields))
			for i, expectedField := range tt.expected.Fields {
				assert.Equal(t, expectedField.Key, result.Fields[i].Key)
				assert.Equal(t, expectedField.Name, result.Fields[i].Name)
				assert.Equal(t, expectedField.ContentType, result.Fields[i].ContentType)
			}
		})
	}
}

func TestGenFieldKeys(t *testing.T) {
	tests := []struct {
		name   string
		fields []*entity.FieldSchema
		want   map[string]string // map[name]key
	}{
		{
			name:   "空字段列表",
			fields: []*entity.FieldSchema{},
			want:   map[string]string{},
		},
		{
			name: "已有key的字段",
			fields: []*entity.FieldSchema{
				{
					Name: "Field One",
					Key:  "existing_key",
				},
			},
			want: map[string]string{
				"Field One": "existing_key",
			},
		},
		{
			name: "需要生成key的字段",
			fields: []*entity.FieldSchema{
				{
					Name: "validName",
					Key:  "",
				},
			},
			want: map[string]string{
				"validName": "validName",
			},
		},
		{
			name: "特殊字符的字段名",
			fields: []*entity.FieldSchema{
				{
					Name: "Field With Space!",
					Key:  "",
				},
			},
			want: map[string]string{
				"Field With Space!": "key_1",
			},
		},
		{
			name: "key冲突的字段",
			fields: []*entity.FieldSchema{
				{
					Name: "test",
					Key:  "test",
				},
				{
					Name: "test2",
					Key:  "",
				},
			},
			want: map[string]string{
				"test":  "test",
				"test2": "test2",
			},
		},
		{
			name: "多个key冲突的字段",
			fields: []*entity.FieldSchema{
				{
					Name: "test",
					Key:  "test",
				},
				{
					Name: "test2",
					Key:  "",
				},
				{
					Name: "test3",
					Key:  "",
				},
			},
			want: map[string]string{
				"test":  "test",
				"test2": "test2",
				"test3": "test3",
			},
		},
		{
			name: "保留字key",
			fields: []*entity.FieldSchema{
				{
					Name: "key",
					Key:  "",
				},
			},
			want: map[string]string{
				"key": "key_1",
			},
		},
		{
			name: "混合场景",
			fields: []*entity.FieldSchema{
				{
					Name: "validName1",
					Key:  "existing_key1",
				},
				{
					Name: "validName2",
					Key:  "",
				},
				{
					Name: "Invalid Name!",
					Key:  "",
				},
				{
					Name: "validName2",
					Key:  "",
				},
			},
			want: map[string]string{
				"validName1":    "existing_key1",
				"validName2":    "validName2_1",
				"Invalid Name!": "key_1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用被测试的函数
			genFieldKeys(tt.fields)

			// 验证结果
			got := make(map[string]string)
			for _, field := range tt.fields {
				got[field.Name] = field.Key
			}

			assert.Equal(t, tt.want, got)

			// 验证生成的key是否都符合标识符规则
			for _, field := range tt.fields {
				assert.True(t, identityRegexp.MatchString(field.Key), "生成的key '%s' 不符合标识符规则", field.Key)
			}

			// 验证key的唯一性
			keys := make(map[string]bool)
			for _, field := range tt.fields {
				assert.False(t, keys[field.Key], "key '%s' 重复", field.Key)
				keys[field.Key] = true
			}
		})
	}
}

func TestMergeSchema(t *testing.T) {
	tests := []struct {
		name      string
		ds        *entity.Dataset
		preSchema *entity.DatasetSchema
		fields    []*entity.FieldSchema
		want      []*entity.FieldSchema
		wantErr   bool
	}{
		{
			name: "成功合并字段 - 新增字段",
			ds: &entity.Dataset{
				ID: 1,
				Spec: &entity.DatasetSpec{
					MaxFieldCount: 10,
				},
				Features: &entity.DatasetFeatures{
					MultiModal: true,
				},
			},
			preSchema: &entity.DatasetSchema{
				Fields: []*entity.FieldSchema{
					{
						Key:         "field1",
						Name:        "Field 1",
						ContentType: entity.ContentTypeText,
						Status:      entity.FieldStatusAvailable,
					},
				},
			},
			fields: []*entity.FieldSchema{
				{
					Key:         "field1",
					Name:        "Field 1",
					ContentType: entity.ContentTypeText,
					Status:      entity.FieldStatusAvailable,
				},
				{
					Key:         "",
					Name:        "Field 2",
					ContentType: entity.ContentTypeText,
					Status:      entity.FieldStatusAvailable,
				},
			},
			want: []*entity.FieldSchema{
				{
					Key:         "field1",
					Name:        "Field 1",
					ContentType: entity.ContentTypeText,
					Status:      entity.FieldStatusAvailable,
				},
				{
					Key:         "key_1",
					Name:        "Field 2",
					ContentType: entity.ContentTypeText,
					Status:      entity.FieldStatusAvailable,
				},
			},
			wantErr: false,
		},
		{
			name: "成功合并字段 - 删除字段",
			ds: &entity.Dataset{
				ID: 1,
				Spec: &entity.DatasetSpec{
					MaxFieldCount: 10,
				},
				Features: &entity.DatasetFeatures{
					MultiModal: true,
				},
			},
			preSchema: &entity.DatasetSchema{
				Fields: []*entity.FieldSchema{
					{
						Key:         "field1",
						Name:        "Field 1",
						ContentType: entity.ContentTypeText,
						Status:      entity.FieldStatusAvailable,
					},
					{
						Key:         "field2",
						Name:        "Field 2",
						ContentType: entity.ContentTypeText,
						Status:      entity.FieldStatusAvailable,
					},
				},
			},
			fields: []*entity.FieldSchema{
				{
					Key:         "field1",
					Name:        "Field 1",
					ContentType: entity.ContentTypeText,
					Status:      entity.FieldStatusAvailable,
				},
			},
			want: []*entity.FieldSchema{
				{
					Key:         "field1",
					Name:        "Field 1",
					ContentType: entity.ContentTypeText,
					Status:      entity.FieldStatusAvailable,
				},
				{
					Key:         "field2",
					Name:        "Field 2",
					ContentType: entity.ContentTypeText,
					Status:      entity.FieldStatusDeleted,
				},
			},
			wantErr: false,
		},
		{
			name: "错误 - 超出最大字段数",
			ds: &entity.Dataset{
				ID: 1,
				Spec: &entity.DatasetSpec{
					MaxFieldCount: 1,
				},
				Features: &entity.DatasetFeatures{
					MultiModal: true,
				},
			},
			preSchema: &entity.DatasetSchema{
				Fields: []*entity.FieldSchema{},
			},
			fields: []*entity.FieldSchema{
				{
					Key:         "field1",
					Name:        "Field 1",
					ContentType: entity.ContentTypeText,
					Status:      entity.FieldStatusAvailable,
				},
				{
					Key:         "field2",
					Name:        "Field 2",
					ContentType: entity.ContentTypeText,
					Status:      entity.FieldStatusAvailable,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "错误 - 没有可用字段",
			ds: &entity.Dataset{
				ID: 1,
				Spec: &entity.DatasetSpec{
					MaxFieldCount: 10,
				},
				Features: &entity.DatasetFeatures{
					MultiModal: true,
				},
			},
			preSchema: &entity.DatasetSchema{
				Fields: []*entity.FieldSchema{},
			},
			fields: []*entity.FieldSchema{
				{
					Key:         "field1",
					Name:        "Field 1",
					ContentType: entity.ContentTypeText,
					Status:      entity.FieldStatusDeleted,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "错误 - 字段名为空",
			ds: &entity.Dataset{
				ID: 1,
				Spec: &entity.DatasetSpec{
					MaxFieldCount: 10,
				},
				Features: &entity.DatasetFeatures{
					MultiModal: true,
				},
			},
			preSchema: &entity.DatasetSchema{
				Fields: []*entity.FieldSchema{},
			},
			fields: []*entity.FieldSchema{
				{
					Key:         "field1",
					Name:        "",
					ContentType: entity.ContentTypeText,
					Status:      entity.FieldStatusAvailable,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "错误 - 非多模态数据集包含多模态字段",
			ds: &entity.Dataset{
				ID: 1,
				Spec: &entity.DatasetSpec{
					MaxFieldCount: 10,
				},
				Features: &entity.DatasetFeatures{
					MultiModal: false,
				},
			},
			preSchema: &entity.DatasetSchema{
				Fields: []*entity.FieldSchema{},
			},
			fields: []*entity.FieldSchema{
				{
					Key:         "field1",
					Name:        "Field 1",
					ContentType: entity.ContentTypeImage,
					Status:      entity.FieldStatusAvailable,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeSchema(tt.ds, tt.preSchema, tt.fields)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.want), len(got))

			// 验证字段内容
			for i, wantField := range tt.want {
				assert.Equal(t, wantField.Key, got[i].Key)
				assert.Equal(t, wantField.Name, got[i].Name)
				assert.Equal(t, wantField.ContentType, got[i].ContentType)
				assert.Equal(t, wantField.Status, got[i].Status)
			}

			// 验证字段key的唯一性
			keys := make(map[string]bool)
			for _, field := range got {
				assert.False(t, keys[field.Key], "key '%s' 重复", field.Key)
				keys[field.Key] = true
			}

			// 验证所有key都符合标识符规则
			for _, field := range got {
				assert.True(t, identityRegexp.MatchString(field.Key), "key '%s' 不符合标识符规则", field.Key)
			}
		})
	}
}
