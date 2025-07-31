// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package data

import (
	"context"
	"testing"

	"github.com/bytedance/gg/gptr"
	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/data/domain/dataset"
	"github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/entity"
)

func TestConvert2DatasetOrderBys(t *testing.T) {
	tests := []struct {
		name     string
		orderBys []*entity.OrderBy
		want     []*dataset.OrderBy
	}{
		{
			name:     "empty order bys",
			orderBys: nil,
			want:     nil,
		},
		{
			name: "single order by",
			orderBys: []*entity.OrderBy{
				{
					Field: gptr.Of("field1"),
					IsAsc: gptr.Of(true),
				},
			},
			want: []*dataset.OrderBy{
				{
					Field: gptr.Of("field1"),
					IsAsc: gptr.Of(true),
				},
			},
		},
		{
			name: "multiple order bys",
			orderBys: []*entity.OrderBy{
				{
					Field: gptr.Of("field1"),
					IsAsc: gptr.Of(true),
				},
				{
					Field: gptr.Of("field2"),
					IsAsc: gptr.Of(true),
				},
			},
			want: []*dataset.OrderBy{
				{
					Field: gptr.Of("field1"),
					IsAsc: gptr.Of(true),
				},
				{
					Field: gptr.Of("field2"),
					IsAsc: gptr.Of(true),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convert2DatasetOrderBys(context.Background(), tt.orderBys)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvert2DatasetMultiModalSpec(t *testing.T) {
	tests := []struct {
		name           string
		multiModalSpec *entity.MultiModalSpec
		want           *dataset.MultiModalSpec
	}{
		{
			name:           "nil spec",
			multiModalSpec: nil,
			want:           nil,
		},
		{
			name: "valid spec",
			multiModalSpec: &entity.MultiModalSpec{
				MaxFileCount:     10,
				MaxFileSize:      1024,
				SupportedFormats: []string{"jpg", "png"},
			},
			want: &dataset.MultiModalSpec{
				MaxFileCount:     gptr.Of(int64(10)),
				MaxFileSize:      gptr.Of(int64(1024)),
				SupportedFormats: []string{"jpg", "png"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convert2DatasetMultiModalSpec(context.Background(), tt.multiModalSpec)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvert2DatasetFieldSchemas(t *testing.T) {
	tests := []struct {
		name    string
		schemas []*entity.FieldSchema
		want    []*dataset.FieldSchema
		wantErr bool
	}{
		{
			name:    "empty schemas",
			schemas: nil,
			want:    nil,
			wantErr: false,
		},
		{
			name: "valid schemas",
			schemas: []*entity.FieldSchema{
				{
					Key:                  "test_key",
					Name:                 "Test Name",
					Description:          "Test Description",
					ContentType:          entity.ContentTypeText,
					DefaultDisplayFormat: entity.FieldDisplayFormat_PlainText,
					Status:               entity.FieldStatus_Available,
					TextSchema:           "test_schema",
					Hidden:               false,
				},
			},
			want: []*dataset.FieldSchema{
				{
					Key:           gptr.Of("test_key"),
					Name:          gptr.Of("Test Name"),
					Description:   gptr.Of("Test Description"),
					ContentType:   gptr.Of(dataset.ContentType_Text),
					DefaultFormat: gptr.Of(dataset.FieldDisplayFormat_PlainText),
					Status:        gptr.Of(dataset.FieldStatus_Available),
					TextSchema:    gptr.Of("test_schema"),
					Hidden:        gptr.Of(false),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid content type",
			schemas: []*entity.FieldSchema{
				{
					Key:         "test_key",
					ContentType: "invalid_type",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convert2DatasetFieldSchemas(context.Background(), tt.schemas)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestConvert2DatasetData(t *testing.T) {
	tests := []struct {
		name    string
		turns   []*entity.Turn
		want    []*dataset.FieldData
		wantErr bool
	}{
		{
			name:    "nil turns",
			turns:   nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty turns",
			turns:   []*entity.Turn{},
			want:    nil,
			wantErr: false,
		},
		{
			name: "valid turns",
			turns: []*entity.Turn{
				{
					FieldDataList: []*entity.FieldData{
						{
							Key:  "test_key",
							Name: "Test Name",
							Content: &entity.Content{
								ContentType: gptr.Of(entity.ContentTypeText),
								Format:      gptr.Of(entity.FieldDisplayFormat_PlainText),
								Text:        gptr.Of("test content"),
							},
						},
					},
				},
			},
			want: []*dataset.FieldData{
				{
					Key:         gptr.Of("test_key"),
					Name:        gptr.Of("Test Name"),
					ContentType: gptr.Of(dataset.ContentType_Text),
					Format:      gptr.Of(dataset.FieldDisplayFormat_PlainText),
					Content:     gptr.Of("test content"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convert2DatasetData(context.Background(), tt.turns)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestConvert2DatasetFieldData(t *testing.T) {
	tests := []struct {
		name      string
		fieldData *entity.FieldData
		want      *dataset.FieldData
		wantErr   bool
	}{
		{
			name:      "nil field data",
			fieldData: nil,
			want:      nil,
			wantErr:   false,
		},
		{
			name: "valid field data",
			fieldData: &entity.FieldData{
				Key:  "test_key",
				Name: "Test Name",
				Content: &entity.Content{
					ContentType: gptr.Of(entity.ContentTypeText),
					Format:      gptr.Of(entity.FieldDisplayFormat_PlainText),
					Text:        gptr.Of("test content"),
				},
			},
			want: &dataset.FieldData{
				Key:         gptr.Of("test_key"),
				Name:        gptr.Of("Test Name"),
				ContentType: gptr.Of(dataset.ContentType_Text),
				Format:      gptr.Of(dataset.FieldDisplayFormat_PlainText),
				Content:     gptr.Of("test content"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convert2DatasetFieldData(context.Background(), tt.fieldData)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestConvert2DatasetItem(t *testing.T) {
	tests := []struct {
		name    string
		item    *entity.EvaluationSetItem
		want    *dataset.DatasetItem
		wantErr bool
	}{
		{
			name:    "nil item",
			item:    nil,
			want:    nil,
			wantErr: false,
		},
		{
			name: "valid item",
			item: &entity.EvaluationSetItem{
				ID:              1,
				AppID:           int32(2),
				SpaceID:         3,
				EvaluationSetID: 4,
				SchemaID:        5,
				ItemID:          1,
				ItemKey:         "key1",
				Turns: []*entity.Turn{
					{
						FieldDataList: []*entity.FieldData{
							{
								Key:  "test_key",
								Name: "Test Name",
								Content: &entity.Content{
									ContentType: gptr.Of(entity.ContentTypeText),
									Format:      gptr.Of(entity.FieldDisplayFormat_PlainText),
									Text:        gptr.Of("test content"),
								},
							},
						},
					},
				},
			},
			want: &dataset.DatasetItem{
				ID:        gptr.Of(int64(1)),
				AppID:     gptr.Of(int32(2)),
				SpaceID:   gptr.Of(int64(3)),
				DatasetID: gptr.Of(int64(4)),
				SchemaID:  gptr.Of(int64(5)),
				ItemID:    gptr.Of(int64(1)),
				ItemKey:   gptr.Of("key1"),
				Data: []*dataset.FieldData{
					{
						Key:         gptr.Of("test_key"),
						Name:        gptr.Of("Test Name"),
						ContentType: gptr.Of(dataset.ContentType_Text),
						Format:      gptr.Of(dataset.FieldDisplayFormat_PlainText),
						Content:     gptr.Of("test content"),
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convert2DatasetItem(context.Background(), tt.item)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestConvert2EvaluationSetSpec(t *testing.T) {
	tests := []struct {
		name string
		spec *dataset.DatasetSpec
		want *entity.DatasetSpec
	}{
		{
			name: "nil spec",
			spec: nil,
			want: nil,
		},
		{
			name: "valid spec",
			spec: &dataset.DatasetSpec{
				MaxFieldCount: gptr.Of(int32(10)),
				MaxItemCount:  gptr.Of(int64(100)),
				MaxItemSize:   gptr.Of(int64(1024)),
			},
			want: &entity.DatasetSpec{
				MaxFieldCount: 10,
				MaxItemCount:  100,
				MaxItemSize:   1024,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convert2EvaluationSetSpec(context.Background(), tt.spec)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvert2DatasetFeatures(t *testing.T) {
	tests := []struct {
		name     string
		features *dataset.DatasetFeatures
		want     *entity.DatasetFeatures
	}{
		{
			name:     "nil features",
			features: nil,
			want:     nil,
		},
		{
			name: "valid features",
			features: &dataset.DatasetFeatures{
				EditSchema:   gptr.Of(true),
				RepeatedData: gptr.Of(true),
				MultiModal:   gptr.Of(true),
			},
			want: &entity.DatasetFeatures{
				EditSchema:   true,
				RepeatedData: true,
				MultiModal:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convert2DatasetFeatures(context.Background(), tt.features)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvert2EvaluationSetMultiModalSpec(t *testing.T) {
	tests := []struct {
		name           string
		multiModalSpec *dataset.MultiModalSpec
		want           *entity.MultiModalSpec
	}{
		{
			name:           "nil spec",
			multiModalSpec: nil,
			want:           nil,
		},
		{
			name: "valid spec",
			multiModalSpec: &dataset.MultiModalSpec{
				MaxFileCount:     gptr.Of(int64(10)),
				MaxFileSize:      gptr.Of(int64(1024)),
				SupportedFormats: []string{"jpg", "png"},
			},
			want: &entity.MultiModalSpec{
				MaxFileCount:     10,
				MaxFileSize:      1024,
				SupportedFormats: []string{"jpg", "png"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convert2EvaluationSetMultiModalSpec(context.Background(), tt.multiModalSpec)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvert2EvaluationSetFieldSchemas(t *testing.T) {
	tests := []struct {
		name    string
		schemas []*dataset.FieldSchema
		want    []*entity.FieldSchema
	}{
		{
			name:    "empty schemas",
			schemas: nil,
			want:    nil,
		},
		{
			name: "valid schemas",
			schemas: []*dataset.FieldSchema{
				{
					Key:         gptr.Of("test_key"),
					Name:        gptr.Of("Test Name"),
					Description: gptr.Of("Test Description"),
					TextSchema:  gptr.Of("test_schema"),
					Hidden:      gptr.Of(false),
					ContentType: gptr.Of(dataset.ContentType_Text),
					MultiModelSpec: &dataset.MultiModalSpec{
						MaxFileCount:     gptr.Of(int64(10)),
						MaxFileSize:      gptr.Of(int64(1024)),
						SupportedFormats: []string{"jpg", "png"},
					},
				},
			},
			want: []*entity.FieldSchema{
				{
					Key:         "test_key",
					Name:        "Test Name",
					Description: "Test Description",
					TextSchema:  "test_schema",
					Hidden:      false,
					ContentType: entity.ContentTypeText,
					MultiModelSpec: &entity.MultiModalSpec{
						MaxFileCount:     10,
						MaxFileSize:      1024,
						SupportedFormats: []string{"jpg", "png"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convert2EvaluationSetFieldSchemas(context.Background(), tt.schemas)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvert2EvaluationSetFieldSchema(t *testing.T) {
	tests := []struct {
		name   string
		schema *dataset.FieldSchema
		want   *entity.FieldSchema
	}{
		{
			name:   "nil schema",
			schema: nil,
			want:   nil,
		},
		{
			name: "valid schema",
			schema: &dataset.FieldSchema{
				Key:         gptr.Of("test_key"),
				Name:        gptr.Of("Test Name"),
				Description: gptr.Of("Test Description"),
				TextSchema:  gptr.Of("test_schema"),
				Hidden:      gptr.Of(false),
				ContentType: gptr.Of(dataset.ContentType_Text),
				MultiModelSpec: &dataset.MultiModalSpec{
					MaxFileCount:     gptr.Of(int64(10)),
					MaxFileSize:      gptr.Of(int64(1024)),
					SupportedFormats: []string{"jpg", "png"},
				},
			},
			want: &entity.FieldSchema{
				Key:         "test_key",
				Name:        "Test Name",
				Description: "Test Description",
				TextSchema:  "test_schema",
				Hidden:      false,
				ContentType: entity.ContentTypeText,
				MultiModelSpec: &entity.MultiModalSpec{
					MaxFileCount:     10,
					MaxFileSize:      1024,
					SupportedFormats: []string{"jpg", "png"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convert2EvaluationSetFieldSchema(context.Background(), tt.schema)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvert2EvaluationSetSchema(t *testing.T) {
	tests := []struct {
		name   string
		schema *dataset.DatasetSchema
		want   *entity.EvaluationSetSchema
	}{
		{
			name:   "nil schema",
			schema: nil,
			want:   nil,
		},
		{
			name: "valid schema",
			schema: &dataset.DatasetSchema{
				ID:        gptr.Of(int64(1)),
				AppID:     gptr.Of(int32(2)),
				SpaceID:   gptr.Of(int64(3)),
				DatasetID: gptr.Of(int64(4)),
				Fields: []*dataset.FieldSchema{
					{
						Key:         gptr.Of("test_key"),
						Name:        gptr.Of("Test Name"),
						ContentType: gptr.Of(dataset.ContentType_Text),
					},
				},
				CreatedAt: gptr.Of(int64(1234567890)),
				UpdatedAt: gptr.Of(int64(1234567890)),
				CreatedBy: gptr.Of("user1"),
				UpdatedBy: gptr.Of("user2"),
			},
			want: &entity.EvaluationSetSchema{
				ID:              1,
				AppID:           2,
				SpaceID:         3,
				EvaluationSetID: 4,
				FieldSchemas: []*entity.FieldSchema{
					{
						Key:         "test_key",
						Name:        "Test Name",
						ContentType: entity.ContentTypeText,
					},
				},
				BaseInfo: &entity.BaseInfo{
					CreatedAt: gptr.Of(int64(1234567890)),
					UpdatedAt: gptr.Of(int64(1234567890)),
					CreatedBy: &entity.UserInfo{UserID: gptr.Of("user1")},
					UpdatedBy: &entity.UserInfo{UserID: gptr.Of("user2")},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convert2EvaluationSetSchema(context.Background(), tt.schema)
			assert.Equal(t, tt.want, got)
		})
	}
}
