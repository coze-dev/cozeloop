// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
)

func TestSanitizeInputItem(t *testing.T) {
	// 创建一个示例 DatasetWithSchema
	ds := &DatasetWithSchema{
		Dataset: &entity.Dataset{
			ID:            1,
			Name:          "TestDataset",
			LatestVersion: "1.0",
			Features: &entity.DatasetFeatures{
				RepeatedData: false,
			},
		},
		Schema: &entity.DatasetSchema{
			ID: 1,
			Fields: []*entity.FieldSchema{
				{
					Key: "name",
				},
				{
					Key: "age",
				},
			},
		},
	}

	// 创建示例 items
	items := []*entity.Item{
		{
			ID:     1,
			ItemID: 1,
			Data: []*entity.FieldData{
				{
					Key:     "name",
					Content: "John Doe",
				},
				{
					Key:     "age",
					Content: "30",
				},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:     2,
			ItemID: 2,
			Data: []*entity.FieldData{
				{
					Key:     "name",
					Content: "Jane Smith",
				},
				{
					Key:     "age",
					Content: "25",
				},
			},
			CreatedAt: time.Now(),
		},
	}

	// 预期结果，这里假设 SanitizeInputItem 不会修改数据
	expected := []*entity.Item{
		{
			ID:     1,
			ItemID: 1,
			Data: []*entity.FieldData{
				{
					Key:     "name",
					Content: "John Doe",
				},
				{
					Key:     "age",
					Content: "30",
				},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:     2,
			ItemID: 2,
			Data: []*entity.FieldData{
				{
					Key:     "name",
					Content: "Jane Smith",
				},
				{
					Key:     "age",
					Content: "25",
				},
			},
			CreatedAt: time.Now(),
		},
	}

	// 调用函数
	SanitizeInputItem(ds, items...)

	// 比较结果
	for i, item := range items {
		if item.ID != expected[i].ID || item.ItemID != expected[i].ItemID {
			t.Errorf("Test case %d failed: ID or ItemID mismatch", i+1)
		}
		for j, field := range item.Data {
			if field.Key != expected[i].Data[j].Key || field.Content != expected[i].Data[j].Content {
				t.Errorf("Test case %d, field %d failed: Key or Content mismatch", i+1, j+1)
			}
		}
	}
}

func TestSanitizeFieldData(t *testing.T) {
	tests := []struct {
		name      string
		inputFD   *entity.FieldData
		walkLevel int
		wantFD    *entity.FieldData
	}{
		{
			name: "ContentTypeText - should trim spaces",
			inputFD: &entity.FieldData{
				ContentType: entity.ContentTypeText,
				Content:     "  hello world  ",
				Parts:       []*entity.FieldData{{Content: "should be removed"}},
			},
			walkLevel: 0,
			wantFD: &entity.FieldData{
				ContentType: entity.ContentTypeText,
				Content:     "hello world",
				Parts:       nil,
			},
		},
		{
			name: "ContentTypeImage - should clear content and parts",
			inputFD: &entity.FieldData{
				ContentType: entity.ContentTypeImage,
				Content:     "some content",
				Parts:       []*entity.FieldData{{Content: "some part"}},
			},
			walkLevel: 0,
			wantFD: &entity.FieldData{
				ContentType: entity.ContentTypeImage,
				Content:     "",
				Parts:       nil,
			},
		},
		{
			name: "ContentTypeAudio - should clear content and parts",
			inputFD: &entity.FieldData{
				ContentType: entity.ContentTypeAudio,
				Content:     "some content",
				Parts:       []*entity.FieldData{{Content: "some part"}},
			},
			walkLevel: 0,
			wantFD: &entity.FieldData{
				ContentType: entity.ContentTypeAudio,
				Content:     "",
				Parts:       nil,
			},
		},
		{
			name: "ContentTypeVideo - should clear content and parts",
			inputFD: &entity.FieldData{
				ContentType: entity.ContentTypeVideo,
				Content:     "some content",
				Parts:       []*entity.FieldData{{Content: "some part"}},
			},
			walkLevel: 0,
			wantFD: &entity.FieldData{
				ContentType: entity.ContentTypeVideo,
				Content:     "",
				Parts:       nil,
			},
		},

		{
			name: "ContentTypeMultiPart - walkLevel 1 should process parts",
			inputFD: &entity.FieldData{
				ContentType: entity.ContentTypeMultiPart,
				Content:     "some content",
				Parts: []*entity.FieldData{
					{Content: "  part1  ", ContentType: entity.ContentTypeText},
					{Content: "", ContentType: entity.ContentTypeText}, // This should be filtered out
					{Content: "  part2  ", ContentType: entity.ContentTypeText},
				},
				Attachments: []*entity.ObjectStorage{{URL: "some-url"}},
			},
			walkLevel: 1,
			wantFD: &entity.FieldData{
				ContentType: entity.ContentTypeMultiPart,
				Content:     "",
				Parts: []*entity.FieldData{
					{Content: "part1", ContentType: entity.ContentTypeText},
					{Content: "part2", ContentType: entity.ContentTypeText},
				},
				Attachments: nil,
			},
		},
		{
			name: "Default case - should clear everything",
			inputFD: &entity.FieldData{
				ContentType: "unknown_type",
				Content:     "some content",
				Parts:       []*entity.FieldData{{Content: "some part"}},
				Attachments: []*entity.ObjectStorage{{URL: "some-url"}},
			},
			walkLevel: 0,
			wantFD: &entity.FieldData{
				ContentType: "unknown_type",
				Content:     "",
				Parts:       nil,
				Attachments: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sanitizeFieldData(tt.inputFD, tt.walkLevel)
			assert.Equal(t, tt.wantFD, tt.inputFD)
		})
	}
}

func TestCastFieldData(t *testing.T) {
	tests := []struct {
		name     string
		schema   *entity.FieldSchema
		data     *entity.FieldData
		expected string
	}{
		{
			name: "Non-text content type should not modify content",
			schema: &entity.FieldSchema{
				ContentType: entity.ContentTypeImage,
				TextSchema:  &entity.JSONSchema{},
			},
			data: &entity.FieldData{
				Content: "TRUE",
			},
			expected: "TRUE",
		},
		{
			name: "Empty content should not be modified",
			schema: &entity.FieldSchema{
				ContentType: entity.ContentTypeText,
				TextSchema:  &entity.JSONSchema{},
			},
			data: &entity.FieldData{
				Content: "",
			},
			expected: "",
		},
		{
			name: "Nil TextSchema should not modify content",
			schema: &entity.FieldSchema{
				ContentType: entity.ContentTypeText,
			},
			data: &entity.FieldData{
				Content: "TRUE",
			},
			expected: "TRUE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			castFieldData(tt.schema, tt.data)
			assert.Equal(t, tt.expected, tt.data.Content)
		})
	}
}

func TestValidateItems(t *testing.T) {
	tests := []struct {
		name          string
		ds            *DatasetWithSchema
		items         []*entity.Item
		wantGoodCount int
		wantBadTypes  []entity.ItemErrorType
	}{
		{
			name: "Valid items within size limit",
			ds: &DatasetWithSchema{
				Dataset: &entity.Dataset{
					ID:   1,
					Name: "test_dataset",
					Spec: &entity.DatasetSpec{
						MaxItemSize: 1000,
					},
				},
				Schema: &entity.DatasetSchema{
					Fields: []*entity.FieldSchema{
						{
							Key:         "field1",
							Name:        "Field 1",
							ContentType: entity.ContentTypeText,
						},
					},
				},
			},
			items: []*entity.Item{
				{
					ID: 1,
					Data: []*entity.FieldData{
						{
							Key:         "field1",
							Content:     "test content",
							ContentType: entity.ContentTypeText,
						},
					},
				},
			},
			wantGoodCount: 1,
			wantBadTypes:  nil,
		},
		{
			name: "Item exceeds size limit",
			ds: &DatasetWithSchema{
				Dataset: &entity.Dataset{
					ID:   2,
					Name: "test_dataset",
					Spec: &entity.DatasetSpec{
						MaxItemSize: 10, // Very small size limit
					},
				},
				Schema: &entity.DatasetSchema{
					Fields: []*entity.FieldSchema{
						{
							Key:         "field1",
							Name:        "Field 1",
							ContentType: entity.ContentTypeText,
						},
					},
				},
			},
			items: []*entity.Item{
				{
					ID: 1,
					Data: []*entity.FieldData{
						{
							Key:         "field1",
							Content:     "this content is too long for the size limit",
							ContentType: entity.ContentTypeText,
						},
					},
				},
			},
			wantGoodCount: 0,
			wantBadTypes:  []entity.ItemErrorType{entity.ItemErrorType_ExceedMaxItemSize},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			good, bad := ValidateItems(tt.ds, tt.items)

			// Check number of good items
			assert.Equal(t, tt.wantGoodCount, len(good), "number of good items mismatch")

			// Check error types
			if tt.wantBadTypes == nil {
				assert.Empty(t, bad, "expected no bad items")
			} else {
				assert.Equal(t, len(tt.wantBadTypes), len(bad), "number of bad item groups mismatch")
				for i, wantType := range tt.wantBadTypes {
					assert.Equal(t, wantType, *bad[i].Type, "error type mismatch")
				}
			}
		})
	}
}

func TestValidateItem(t *testing.T) {
	tests := []struct {
		name        string
		ds          *DatasetWithSchema
		item        *entity.Item
		wantErrCode int32
	}{
		{
			name: "Valid item",
			ds: &DatasetWithSchema{
				Dataset: &entity.Dataset{
					ID:   1,
					Name: "test_dataset",
					Spec: &entity.DatasetSpec{
						MaxItemSize: 1000,
					},
				},
				Schema: &entity.DatasetSchema{
					Fields: []*entity.FieldSchema{
						{
							Key:         "field1",
							Name:        "Field 1",
							ContentType: entity.ContentTypeText,
						},
					},
				},
			},
			item: &entity.Item{
				ID: 1,
				Data: []*entity.FieldData{
					{
						Key:         "field1",
						Content:     "test content",
						ContentType: entity.ContentTypeText,
					},
				},
			},
			wantErrCode: 0, // no error
		},
		{
			name: "Item exceeds size limit",
			ds: &DatasetWithSchema{
				Dataset: &entity.Dataset{
					ID:   2,
					Name: "test_dataset",
					Spec: &entity.DatasetSpec{
						MaxItemSize: 10, // Very small size limit
					},
				},
				Schema: &entity.DatasetSchema{
					Fields: []*entity.FieldSchema{
						{
							Key:         "field1",
							Name:        "Field 1",
							ContentType: entity.ContentTypeText,
						},
					},
				},
			},
			item: &entity.Item{
				ID: 1,
				Data: []*entity.FieldData{
					{
						Key:         "field1",
						Content:     "this content is too long for the size limit",
						ContentType: entity.ContentTypeText,
					},
				},
			},
			wantErrCode: errno.ItemDataSizeExceededCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateItem(tt.ds, tt.item)
			if tt.wantErrCode == 0 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestSanitizeOutputItem(t *testing.T) {
	tests := []struct {
		name     string
		schema   *entity.DatasetSchema
		items    []*entity.Item
		expected []*entity.Item
	}{
		{
			name: "Basic field filtering and property filling",
			schema: &entity.DatasetSchema{
				Fields: []*entity.FieldSchema{
					{
						Key:           "field1",
						Name:          "Field One",
						ContentType:   entity.ContentTypeText,
						DefaultFormat: "plain",
					},
					{
						Key:           "field2",
						Name:          "Field Two",
						ContentType:   entity.ContentTypeImage,
						DefaultFormat: "url",
					},
				},
			},
			items: []*entity.Item{
				{
					ID: 1,
					Data: []*entity.FieldData{
						{
							Key:         "field1",
							Content:     "test content",
							ContentType: "",
							Format:      "",
						},
						{
							Key:         "field2",
							Content:     "image.jpg",
							ContentType: "",
							Format:      "",
						},
						{
							Key:         "field3", // This field should be filtered out
							Content:     "should not appear",
							ContentType: entity.ContentTypeText,
						},
					},
				},
			},
			expected: []*entity.Item{
				{
					ID: 1,
					Data: []*entity.FieldData{
						{
							Key:         "field1",
							Name:        "Field One",
							Content:     "test content",
							ContentType: entity.ContentTypeText,
							Format:      "plain",
						},
						{
							Key:         "field2",
							Name:        "Field Two",
							Content:     "image.jpg",
							ContentType: entity.ContentTypeImage,
							Format:      "url",
						},
					},
				},
			},
		},
		{
			name: "Handle repeated data",
			schema: &entity.DatasetSchema{
				Fields: []*entity.FieldSchema{
					{
						Key:           "field1",
						Name:          "Field One",
						ContentType:   entity.ContentTypeText,
						DefaultFormat: "plain",
					},
				},
			},
			items: []*entity.Item{
				{
					ID: 1,
					RepeatedData: []*entity.ItemData{
						{
							Data: []*entity.FieldData{
								{
									Key:         "field1",
									Content:     "repeated content 1",
									ContentType: "",
									Format:      "",
								},
								{
									Key:         "field2", // Should be filtered out
									Content:     "should not appear",
									ContentType: entity.ContentTypeText,
								},
							},
						},
					},
				},
			},
			expected: []*entity.Item{
				{
					ID: 1,
					RepeatedData: []*entity.ItemData{
						{
							Data: []*entity.FieldData{
								{
									Key:         "field1",
									Name:        "Field One",
									Content:     "repeated content 1",
									ContentType: entity.ContentTypeText,
									Format:      "plain",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Empty items list",
			schema: &entity.DatasetSchema{
				Fields: []*entity.FieldSchema{
					{
						Key:           "field1",
						Name:          "Field One",
						ContentType:   entity.ContentTypeText,
						DefaultFormat: "plain",
					},
				},
			},
			items:    []*entity.Item{},
			expected: []*entity.Item{},
		},
		{
			name: "Item with no data",
			schema: &entity.DatasetSchema{
				Fields: []*entity.FieldSchema{
					{
						Key:           "field1",
						Name:          "Field One",
						ContentType:   entity.ContentTypeText,
						DefaultFormat: "plain",
					},
				},
			},
			items: []*entity.Item{
				{
					ID:   1,
					Data: []*entity.FieldData{},
				},
			},
			expected: []*entity.Item{
				{
					ID:   1,
					Data: []*entity.FieldData{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a deep copy of items to avoid modifying the test case data
			items := make([]*entity.Item, len(tt.items))
			for i, item := range tt.items {
				items[i] = &entity.Item{
					ID:           item.ID,
					Data:         make([]*entity.FieldData, len(item.Data)),
					RepeatedData: make([]*entity.ItemData, len(item.RepeatedData)),
				}
				for j, fd := range item.Data {
					items[i].Data[j] = &entity.FieldData{
						Key:         fd.Key,
						Name:        fd.Name,
						Content:     fd.Content,
						ContentType: fd.ContentType,
						Format:      fd.Format,
					}
				}
				for j, rd := range item.RepeatedData {
					items[i].RepeatedData[j] = &entity.ItemData{
						Data: make([]*entity.FieldData, len(rd.Data)),
					}
					for k, fd := range rd.Data {
						items[i].RepeatedData[j].Data[k] = &entity.FieldData{
							Key:         fd.Key,
							Name:        fd.Name,
							Content:     fd.Content,
							ContentType: fd.ContentType,
							Format:      fd.Format,
						}
					}
				}
			}

			SanitizeOutputItem(tt.schema, items)

			assert.Equal(t, len(tt.expected), len(items), "items length mismatch")
			for i, expectedItem := range tt.expected {
				assert.Equal(t, expectedItem.ID, items[i].ID, "item ID mismatch")

				// Check regular data
				assert.Equal(t, len(expectedItem.Data), len(items[i].Data), "data length mismatch")
				for j, expectedFD := range expectedItem.Data {
					assert.Equal(t, expectedFD.Key, items[i].Data[j].Key, "field key mismatch")
					assert.Equal(t, expectedFD.Name, items[i].Data[j].Name, "field name mismatch")
					assert.Equal(t, expectedFD.Content, items[i].Data[j].Content, "field content mismatch")
					assert.Equal(t, expectedFD.ContentType, items[i].Data[j].ContentType, "field content type mismatch")
					assert.Equal(t, expectedFD.Format, items[i].Data[j].Format, "field format mismatch")
				}

				// Check repeated data
				assert.Equal(t, len(expectedItem.RepeatedData), len(items[i].RepeatedData), "repeated data length mismatch")
				for j, expectedRD := range expectedItem.RepeatedData {
					assert.Equal(t, len(expectedRD.Data), len(items[i].RepeatedData[j].Data), "repeated data fields length mismatch")
					for k, expectedFD := range expectedRD.Data {
						assert.Equal(t, expectedFD.Key, items[i].RepeatedData[j].Data[k].Key, "repeated field key mismatch")
						assert.Equal(t, expectedFD.Name, items[i].RepeatedData[j].Data[k].Name, "repeated field name mismatch")
						assert.Equal(t, expectedFD.Content, items[i].RepeatedData[j].Data[k].Content, "repeated field content mismatch")
						assert.Equal(t, expectedFD.ContentType, items[i].RepeatedData[j].Data[k].ContentType, "repeated field content type mismatch")
						assert.Equal(t, expectedFD.Format, items[i].RepeatedData[j].Data[k].Format, "repeated field format mismatch")
					}
				}
			}
		})
	}
}
