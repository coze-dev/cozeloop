// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

func TestDatasetSchema_AvailableFields(t *testing.T) {
	type fields struct {
		ID            int64
		AppID         int32
		SpaceID       int64
		DatasetID     int64
		Fields        []*FieldSchema
		Immutable     bool
		CreatedBy     string
		CreatedAt     time.Time
		UpdatedBy     string
		UpdatedAt     time.Time
		UpdateVersion int64
	}
	tests := []struct {
		name   string
		fields fields
		want   []*FieldSchema
	}{
		{
			name: "Test available fields",
			fields: fields{
				Fields: []*FieldSchema{
					{Status: FieldStatusAvailable},
					{Status: FieldStatusDeleted},
				},
			},
			want: []*FieldSchema{
				{Status: FieldStatusAvailable},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatasetSchema{
				ID:            tt.fields.ID,
				AppID:         tt.fields.AppID,
				SpaceID:       tt.fields.SpaceID,
				DatasetID:     tt.fields.DatasetID,
				Fields:        tt.fields.Fields,
				Immutable:     tt.fields.Immutable,
				CreatedBy:     tt.fields.CreatedBy,
				CreatedAt:     tt.fields.CreatedAt,
				UpdatedBy:     tt.fields.UpdatedBy,
				UpdatedAt:     tt.fields.UpdatedAt,
				UpdateVersion: tt.fields.UpdateVersion,
			}
			if got := s.AvailableFields(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AvailableFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatasetSchema_GetID(t *testing.T) {
	type fields struct {
		ID            int64
		AppID         int32
		SpaceID       int64
		DatasetID     int64
		Fields        []*FieldSchema
		Immutable     bool
		CreatedBy     string
		CreatedAt     time.Time
		UpdatedBy     string
		UpdatedAt     time.Time
		UpdateVersion int64
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "Test get ID",
			fields: fields{
				ID: 123,
			},
			want: 123,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatasetSchema{
				ID:            tt.fields.ID,
				AppID:         tt.fields.AppID,
				SpaceID:       tt.fields.SpaceID,
				DatasetID:     tt.fields.DatasetID,
				Fields:        tt.fields.Fields,
				Immutable:     tt.fields.Immutable,
				CreatedBy:     tt.fields.CreatedBy,
				CreatedAt:     tt.fields.CreatedAt,
				UpdatedBy:     tt.fields.UpdatedBy,
				UpdatedAt:     tt.fields.UpdatedAt,
				UpdateVersion: tt.fields.UpdateVersion,
			}
			if got := s.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatasetSchema_SetID(t *testing.T) {
	type fields struct {
		ID            int64
		AppID         int32
		SpaceID       int64
		DatasetID     int64
		Fields        []*FieldSchema
		Immutable     bool
		CreatedBy     string
		CreatedAt     time.Time
		UpdatedBy     string
		UpdatedAt     time.Time
		UpdateVersion int64
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
			name: "Test set ID",
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
			s := &DatasetSchema{
				ID:            tt.fields.ID,
				AppID:         tt.fields.AppID,
				SpaceID:       tt.fields.SpaceID,
				DatasetID:     tt.fields.DatasetID,
				Fields:        tt.fields.Fields,
				Immutable:     tt.fields.Immutable,
				CreatedBy:     tt.fields.CreatedBy,
				CreatedAt:     tt.fields.CreatedAt,
				UpdatedBy:     tt.fields.UpdatedBy,
				UpdatedAt:     tt.fields.UpdatedAt,
				UpdateVersion: tt.fields.UpdateVersion,
			}
			s.SetID(tt.args.id)
			if s.ID != tt.args.id {
				t.Errorf("SetID() failed, expected ID: %v, got: %v", tt.args.id, s.ID)
			}
		})
	}
}

func TestFieldSchema_Available(t *testing.T) {
	type fields struct {
		Key            string
		Name           string
		Description    string
		ContentType    ContentType
		DefaultFormat  FieldDisplayFormat
		SchemaKey      SchemaKey
		TextSchema     *JSONSchema
		MultiModelSpec *MultiModalSpec
		Status         FieldStatus
		Hidden         bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Test available field",
			fields: fields{
				Status: FieldStatusAvailable,
			},
			want: true,
		},
		{
			name: "Test unavailable field",
			fields: fields{
				Status: FieldStatusDeleted,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FieldSchema{
				Key:            tt.fields.Key,
				Name:           tt.fields.Name,
				Description:    tt.fields.Description,
				ContentType:    tt.fields.ContentType,
				DefaultFormat:  tt.fields.DefaultFormat,
				SchemaKey:      tt.fields.SchemaKey,
				TextSchema:     tt.fields.TextSchema,
				MultiModelSpec: tt.fields.MultiModelSpec,
				Status:         tt.fields.Status,
				Hidden:         tt.fields.Hidden,
			}
			if got := s.Available(); got != tt.want {
				t.Errorf("Available() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldSchema_CompatibleWith(t *testing.T) {
	type fields struct {
		Key            string
		Name           string
		Description    string
		ContentType    ContentType
		DefaultFormat  FieldDisplayFormat
		SchemaKey      SchemaKey
		TextSchema     *JSONSchema
		MultiModelSpec *MultiModalSpec
		Status         FieldStatus
		Hidden         bool
	}
	type args struct {
		other *FieldSchema
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Test compatible fields",
			fields: fields{
				ContentType: ContentTypeText,
				SchemaKey:   "text",
			},
			args: args{
				other: &FieldSchema{
					ContentType: ContentTypeText,
					SchemaKey:   "text",
				},
			},
			want: true,
		},
		{
			name: "Test incompatible fields",
			fields: fields{
				ContentType: ContentTypeText,
			},
			args: args{
				other: &FieldSchema{
					ContentType: ContentTypeImage,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FieldSchema{
				Key:            tt.fields.Key,
				Name:           tt.fields.Name,
				Description:    tt.fields.Description,
				ContentType:    tt.fields.ContentType,
				DefaultFormat:  tt.fields.DefaultFormat,
				SchemaKey:      tt.fields.SchemaKey,
				TextSchema:     tt.fields.TextSchema,
				MultiModelSpec: tt.fields.MultiModelSpec,
				Status:         tt.fields.Status,
				Hidden:         tt.fields.Hidden,
			}
			if got := s.CompatibleWith(tt.args.other); got != tt.want {
				t.Errorf("CompatibleWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldSchema_ValidateData(t *testing.T) {
	type fields struct {
		Key            string
		Name           string
		Description    string
		ContentType    ContentType
		DefaultFormat  FieldDisplayFormat
		SchemaKey      SchemaKey
		TextSchema     *JSONSchema
		MultiModelSpec *MultiModalSpec
		Status         FieldStatus
		Hidden         bool
	}
	type args struct {
		d *FieldData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test text data",
			fields: fields{
				Key:         "test",
				ContentType: ContentTypeText,
				TextSchema: &JSONSchema{
					Raw: []byte(`{"type": "string"}`),
				},
			},
			args: args{
				d: &FieldData{
					Key:     "test",
					Content: `123`,
				},
			},
			wantErr: false,
		},
		{
			name: "Test image data",
			fields: fields{
				Key:            "test",
				ContentType:    ContentTypeImage,
				MultiModelSpec: &MultiModalSpec{},
			},
			args: args{
				d: &FieldData{
					Key:     "test",
					Content: `123`,
				},
			},
			wantErr: false,
		},
		{
			name: "Test image data, space is nil",
			fields: fields{
				Key:         "test",
				ContentType: ContentTypeImage,
			},
			args: args{
				d: &FieldData{
					Key:     "test",
					Content: `123`,
				},
			},
			wantErr: false,
		},
		{
			name: "Test image data, count err",
			fields: fields{
				Key:         "test",
				ContentType: ContentTypeImage,
				MultiModelSpec: &MultiModalSpec{
					MaxFileCount: 1,
				},
			},
			args: args{
				d: &FieldData{
					Key:         "test",
					Content:     `123`,
					Attachments: []*ObjectStorage{{}, {}},
				},
			},
			wantErr: true,
		},
		{
			name: "Test invalid parts data",
			fields: fields{
				ContentType: ContentTypeMultiPart,
				Key:         "test",
			},
			args: args{
				d: &FieldData{
					Key:     "test",
					Content: `123`,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FieldSchema{
				Key:            tt.fields.Key,
				Name:           tt.fields.Name,
				Description:    tt.fields.Description,
				ContentType:    tt.fields.ContentType,
				DefaultFormat:  tt.fields.DefaultFormat,
				SchemaKey:      tt.fields.SchemaKey,
				TextSchema:     tt.fields.TextSchema,
				MultiModelSpec: tt.fields.MultiModelSpec,
				Status:         tt.fields.Status,
				Hidden:         tt.fields.Hidden,
			}
			if err := s.ValidateData(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("ValidateData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFieldSchema_validateTextData(t *testing.T) {
	type fields struct {
		Key            string
		Name           string
		Description    string
		ContentType    ContentType
		DefaultFormat  FieldDisplayFormat
		SchemaKey      SchemaKey
		TextSchema     *JSONSchema
		MultiModelSpec *MultiModalSpec
		Status         FieldStatus
		Hidden         bool
	}
	type args struct {
		d *FieldData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test valid text data",
			fields: fields{
				TextSchema: &JSONSchema{
					Raw: []byte(`{"type": "string"}`),
				},
			},
			args: args{
				d: &FieldData{
					Content: `"test"`,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FieldSchema{
				Key:            tt.fields.Key,
				Name:           tt.fields.Name,
				Description:    tt.fields.Description,
				ContentType:    tt.fields.ContentType,
				DefaultFormat:  tt.fields.DefaultFormat,
				SchemaKey:      tt.fields.SchemaKey,
				TextSchema:     tt.fields.TextSchema,
				MultiModelSpec: tt.fields.MultiModelSpec,
				Status:         tt.fields.Status,
				Hidden:         tt.fields.Hidden,
			}
			if err := s.validateTextData(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("validateTextData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJSONSchema_CompatibleWith(t *testing.T) {
	type fields struct {
		Raw    json.RawMessage
		Schema *gojsonschema.Schema
	}
	type args struct {
		other *JSONSchema
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Test compatible schemas",
			fields: fields{
				Raw: []byte(`{"type": "string"}`),
			},
			args: args{
				other: &JSONSchema{
					Raw: []byte(`{"type": "string"}`),
				},
			},
			want: true,
		},
		{
			name: "Test incompatible schemas",
			fields: fields{
				Raw: []byte(`{"type": "string"}`),
			},
			args: args{
				other: &JSONSchema{
					Raw: []byte(`{"type": "number"}`),
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &JSONSchema{
				Raw:    tt.fields.Raw,
				Schema: tt.fields.Schema,
			}
			if got := s.CompatibleWith(tt.args.other); got != tt.want {
				t.Errorf("CompatibleWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONSchema_GetSingleType(t *testing.T) {
	type fields struct {
		Raw    json.RawMessage
		Schema *gojsonschema.Schema
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test single type",
			fields: fields{
				Raw: []byte(`{"type": "string"}`),
			},
			want: "string",
		},
		{
			name: "Test multiple types",
			fields: fields{
				Raw: []byte(`{"type": ["string", "number"]}`),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &JSONSchema{
				Raw:    tt.fields.Raw,
				Schema: tt.fields.Schema,
			}
			if got := s.GetSingleType(); got != tt.want {
				t.Errorf("GetSingleType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONSchema_MarshalJSON(t *testing.T) {
	type fields struct {
		Raw    json.RawMessage
		Schema *gojsonschema.Schema
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Test marshal JSON",
			fields: fields{
				Raw: []byte(`{"type": "string"}`),
			},
			want:    []byte(`{"type": "string"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &JSONSchema{
				Raw:    tt.fields.Raw,
				Schema: tt.fields.Schema,
			}
			got, err := s.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONSchema_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Raw    json.RawMessage
		Schema *gojsonschema.Schema
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Test unmarshal JSON",
			fields: fields{},
			args: args{
				data: []byte(`{"type": "string"}`),
			},
			wantErr: false,
		},
		{
			name:   "Test invalid JSON",
			fields: fields{},
			args: args{
				data: []byte(`invalid`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &JSONSchema{
				Raw:    tt.fields.Raw,
				Schema: tt.fields.Schema,
			}
			if err := s.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFieldSchemaValidateData(t *testing.T) {
	objSchema := newJSONObjectSchema(t)
	for _, tc := range []struct {
		name    string
		schema  *FieldSchema
		data    *FieldData
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "nil data",
			schema:  &FieldSchema{ContentType: ContentTypeText},
			data:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "key mismatch",
			schema:  &FieldSchema{Key: "key1", ContentType: ContentTypeText},
			data:    &FieldData{Key: "key2", ContentType: ContentTypeText},
			wantErr: assert.Error,
		},
		{
			name:    "text schema key",
			schema:  &FieldSchema{SchemaKey: SchemaKeyInteger, ContentType: ContentTypeText},
			data:    &FieldData{Content: "123"},
			wantErr: assert.NoError,
		},
		{
			name:    "text json schema",
			schema:  &FieldSchema{ContentType: ContentTypeText, TextSchema: objSchema},
			data:    &FieldData{Content: `{"first_name": "John", "last_name": "Doe", "birthday": "1990-01-01", "address": {"street_address": "123 Main St", "city": "Anytown", "state": "CA", "country": "USA"}}`},
			wantErr: assert.NoError,
		},
		{
			name:    "text json schema error",
			schema:  &FieldSchema{ContentType: ContentTypeText, TextSchema: objSchema},
			data:    &FieldData{Content: `{"first_name": "John", "last_name": "Doe", "birthday": "bad date", "address": {"street_address": "123 Main St", "city": "Anytown", "state": "CA", "country": "USA"}}`},
			wantErr: assert.Error,
		},
		{
			name:    "image",
			schema:  &FieldSchema{ContentType: ContentTypeImage, MultiModelSpec: &MultiModalSpec{MaxFileCount: 1}},
			data:    &FieldData{Attachments: []*ObjectStorage{{Provider: "tos", URI: "test.jpg"}}},
			wantErr: assert.NoError,
		},
		{
			name:   "image error",
			schema: &FieldSchema{ContentType: ContentTypeImage, MultiModelSpec: &MultiModalSpec{MaxFileCount: 1}},
			data: &FieldData{Attachments: []*ObjectStorage{
				{Provider: "tos", URI: "test1.jpg"},
				{Provider: "tos", URI: "test2.jpg"},
			}},
			wantErr: assert.Error,
		},
		{
			name:    "invalid json", // only valid json allowed to be checked by jsonSchema, otherwise the input content may be truncated
			schema:  &FieldSchema{ContentType: ContentTypeText, SchemaKey: SchemaKeyInteger},
			data:    &FieldData{Content: `2024-01-01`},
			wantErr: assert.Error,
		},
		{
			name:    "content with quote",
			schema:  &FieldSchema{ContentType: ContentTypeText, SchemaKey: SchemaKeyString},
			data:    &FieldData{Content: `"name": "alice"`},
			wantErr: assert.NoError,
		},
		{
			name:    "content with control character",
			schema:  &FieldSchema{ContentType: ContentTypeText, SchemaKey: SchemaKeyString},
			data:    &FieldData{Content: `Let \x03$*$ be an operation`},
			wantErr: assert.NoError,
		},
		{
			name:    "boolean content",
			schema:  &FieldSchema{ContentType: ContentTypeText, SchemaKey: SchemaKeyBool},
			data:    &FieldData{Content: `FaLSe`},
			wantErr: assert.NoError,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.schema.ValidateData(tc.data)
			tc.wantErr(t, err)
		})
	}
}

func newJSONObjectSchema(t *testing.T) *JSONSchema {
	schema, err := NewJSONSchema(`{
  "type": "object",
  "properties": {
    "first_name": { "type": "string" },
    "last_name": { "type": "string" },
    "birthday": { "type": "string", "format": "date" },
    "address": {
       "type": "object",
       "properties": {
         "street_address": { "type": "string" },
         "city": { "type": "string" },
         "state": { "type": "string" },
         "country": { "type" : "string" }
       }
    }
  }
}
`)
	require.NoError(t, err)
	return schema
}

func TestJSONSchema_getTypes(t *testing.T) {
	type fields struct {
		Raw    json.RawMessage
		Schema *gojsonschema.Schema
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Test single type",
			fields: fields{
				Raw: []byte(`{"type": "string"}`),
			},
			want: []string{"string"},
		},
		{
			name: "Test multiple types",
			fields: fields{
				Raw: []byte(`{"type": ["string", "number"]}`),
			},
			want: []string{"string", "number"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &JSONSchema{
				Raw:    tt.fields.Raw,
				Schema: tt.fields.Schema,
			}
			if got := s.getTypes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJSONSchema(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name    string
		args    args
		want    *JSONSchema
		wantErr bool
	}{
		{
			name: "Test valid schema",
			args: args{
				raw: `{"type": "string"}`,
			},
			want: &JSONSchema{
				Raw: []byte(`{"type": "string"}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJSONSchema(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJSONSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Raw, tt.want.Raw) {
				t.Errorf("NewJSONSchema() got = %v, want %v", got, tt.want)
			}
		})
	}
}
