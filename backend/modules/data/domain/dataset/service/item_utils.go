// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"fmt"
	"math"
	"strings"

	"github.com/bytedance/gg/collection/set"
	"github.com/bytedance/gg/gcond"
	"github.com/bytedance/gg/gmap"
	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gslice"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/consts"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
)

// SanitizeInputItem 根据 schema 等修剪、填补 item 数据内容，使之与类型相匹配。用于处理用户传入的数据。
func SanitizeInputItem(ds *DatasetWithSchema, items ...*entity.Item) {
	var (
		repeatedData = ds.Features.RepeatedData
		fields       = gslice.ToMap(ds.Schema.Fields, func(f *entity.FieldSchema) (string, *entity.FieldSchema) { return f.Key, f })
		nameToKey    = gslice.ToMap(ds.Schema.AvailableFields(), func(f *entity.FieldSchema) (string, string) { return f.Name, f.Key })
	)

	for _, item := range items {
		switch {
		case repeatedData:
			item.Data = nil
			keep := 0
			for _, data := range item.RepeatedData {
				pruned := sanitizeItemData(data.Data, fields, nameToKey)
				if len(pruned) == 0 {
					continue
				}
				data.Data = pruned
				item.RepeatedData[keep] = data
				keep++
			}
			item.RepeatedData = item.RepeatedData[:keep]

		default:
			item.RepeatedData = nil
			item.Data = sanitizeItemData(item.Data, fields, nameToKey)
		}
	}
}

func sanitizeItemData(data []*entity.FieldData, fields map[string]*entity.FieldSchema, nameToKey map[string]string) []*entity.FieldData {
	keep := 0
	for _, fd := range data {
		key := fd.Key
		if key == "" { // 请求中的 FieldData 中允许仅指定 PromptName 不指定 Key, 此处做填充。
			key = nameToKey[fd.Name]
			fd.Key = key
		}
		schema, ok := fields[key]
		if !ok {
			continue
		}

		fd.ContentType = schema.ContentType
		sanitizeFieldData(fd, 1)
		if fd.DataBytes() == 0 {
			continue
		}
		castFieldData(schema, fd)

		data[keep] = fd
		keep++
	}
	return data[:keep]
}

func sanitizeFieldData(fd *entity.FieldData, walkLevel int) {
	switch fd.ContentType {
	case entity.ContentTypeText:
		fd.Parts = nil
		fd.Content = strings.TrimSpace(fd.Content)

	case entity.ContentTypeImage, entity.ContentTypeAudio, entity.ContentTypeVideo:
		fd.Content = ""
		fd.Parts = nil

	case entity.ContentTypeMultiPart:
		fd.Content = ""
		fd.Attachments = nil
		if walkLevel == 0 {
			fd.Parts = nil
		}
		for _, part := range fd.Parts {
			sanitizeFieldData(part, walkLevel-1)
		}
		fd.Parts = gslice.Filter(fd.Parts, func(part *entity.FieldData) bool { return part.DataBytes() > 0 })

	default:
		fd.Content = ""
		fd.Parts = nil
		fd.Attachments = nil
	}
}

// castFieldData 将 FieldData 的内容转换为符合 schema 的类型。
func castFieldData(s *entity.FieldSchema, d *entity.FieldData) {
	if s.ContentType != entity.ContentTypeText || s.TextSchema == nil || s.TextSchema.Schema == nil || d.Content == "" {
		return
	}

	switch s.TextSchema.GetSingleType() {
	case consts.TypeBoolean:
		d.Content = strings.ToLower(d.Content)
	default:
	}
}

// ValidateItems 校验 items 是否符合 schema 等约束，返回信息中包含入参的 index 信息。
func ValidateItems(ds *DatasetWithSchema, items []*entity.Item) (good []*IndexedItem, bad []*entity.ItemErrorGroup) {
	iitems := make([]*IndexedItem, 0, len(items))
	for i, item := range items {
		iitems = append(iitems, &IndexedItem{Index: i, Item: item})
	}

	return ValidateIndexedItems(ds, iitems)
}

func ValidateIndexedItems(ds *DatasetWithSchema, items []*IndexedItem) (good []*IndexedItem, bad []*entity.ItemErrorGroup) {
	var (
		schemaByKey = gslice.ToMap(ds.Schema.AvailableFields(), func(f *entity.FieldSchema) (string, *entity.FieldSchema) { return f.Key, f })
		maxItemSize = ds.Spec.MaxItemSize
		errMap      = make(map[entity.ItemErrorType]*entity.ItemErrorGroup)
		keep        = 0
	)

	if maxItemSize == 0 {
		maxItemSize = math.MaxInt64
	}

	addErrItem := func(index int, errType entity.ItemErrorType, message string) {
		pre, ok := errMap[errType]
		if !ok {
			pre = &entity.ItemErrorGroup{Type: gptr.Of(errType), Details: make([]*entity.ItemErrorDetail, 0)}
			errMap[errType] = pre
		}
		pre.ErrorCount = gptr.Of(gptr.Indirect(pre.ErrorCount) + 1)
		pre.Details = append(pre.Details, &entity.ItemErrorDetail{Message: gptr.Of(message), Index: gptr.Of(int32(index))})
	}

	for _, item := range items {
		props := item.GetOrBuildProperties()
		if size := props.Bytes; size > maxItemSize {
			addErrItem(item.Index, entity.ItemErrorType_ExceedMaxItemSize, fmt.Sprintf("size of item %d exceeds max %d", size, maxItemSize))
			continue
		}

		hasInvalidData := false
		for _, data := range item.AllData() {
			dm := gslice.ToMap(data, func(fd *entity.FieldData) (string, *entity.FieldData) { return fd.Key, fd })
			for key, schema := range schemaByKey {
				field, ok := dm[key]
				if !ok {
					continue
				}
				if err := schema.ValidateData(field); err != nil {
					hasInvalidData = true
					addErrItem(item.Index, entity.ItemErrorType_MismatchSchema, fmt.Sprintf("field_name=%s, msg=%s", schema.Name, err.Error()))
					continue
				}
			}
		}
		if !hasInvalidData {
			items[keep] = item
			keep++
		}
	}

	return items[:keep], gmap.Values(errMap)
}

func ValidateItem(ds *DatasetWithSchema, item *entity.Item) error {
	_, bad := ValidateItems(ds, []*entity.Item{item})
	if len(bad) > 0 {
		b := bad[0]
		var errorCode int32
		switch gptr.Indirect(b.Type) {
		case entity.ItemErrorType_ExceedMaxItemSize:
			errorCode = errno.ItemDataSizeExceededCode
		case entity.ItemErrorType_MismatchSchema:
			errorCode = errno.SchemaMismatchCode
		case entity.ItemErrorType_ExceedDatasetCapacity:
			errorCode = errno.DatasetCapacityFullCode
		default:
			errorCode = errno.CommonBadRequestCode
		}
		msg := fmt.Sprintf("reason=%v", b.Type)
		if len(b.Details) > 0 {
			msg = fmt.Sprintf("reason=%v, message=%v", b.Type, b.Details[0].Message)
		}
		return errno.Errorf(errorCode, "invalid item, %s", msg)
	}
	return nil
}

// SanitizeOutputItem 根据 schema 修剪 items 内容。用于处理传给用户的数据。
func SanitizeOutputItem(schema *entity.DatasetSchema, items []*entity.Item) {
	fields := schema.AvailableFields()
	keys := gslice.Map(fields, func(f *entity.FieldSchema) string { return f.Key })
	key2Schema := gslice.ToMap(fields, func(f *entity.FieldSchema) (string, *entity.FieldSchema) { return f.Key, f })
	keySet := set.New(keys...)
	for _, item := range items {
		// 1. 隐藏已删除的字段
		if len(item.Data) > 0 {
			item.Data = gslice.Filter(item.Data, func(f *entity.FieldData) bool { return keySet.Contains(f.Key) })
		}
		for i, data := range item.RepeatedData {
			item.RepeatedData[i].Data = gslice.Filter(data.Data, func(f *entity.FieldData) bool { return keySet.Contains(f.Key) })
		}
		// 2. 字段回填
		for _, data := range item.AllData() {
			for _, field := range data {
				schema := key2Schema[field.Key]
				field.Name = schema.Name
				field.ContentType = gcond.If(field.ContentType == "", schema.ContentType, field.ContentType)
				field.Format = gcond.If(field.Format == "", schema.DefaultFormat, field.Format)
			}
		}
	}
}

var DatasetItemDataKey = `dataset:%d:item:%d:vn:%d` // dataset:dataset_id:item:item_id:vn:add_vn, string, item 内容

func FormatDatasetItemDataKey(datasetID, itemID, vn int64) string {
	return fmt.Sprintf(DatasetItemDataKey, datasetID, itemID, vn)
}
