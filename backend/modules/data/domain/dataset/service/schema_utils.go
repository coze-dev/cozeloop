// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"regexp"
	"strconv"

	"github.com/bytedance/gg/collection/set"
	"github.com/bytedance/gg/gslice"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
)

var identityRegexp = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]{0,63}$")

// genFieldKeys 为 fields 生成 key， 确保 key 不冲突。
func genFieldKeys(fields []*entity.FieldSchema) {
	seen := set.New(gslice.Map(fields, func(f *entity.FieldSchema) string { return f.Key })...)
	seen.Add(`key`) // reserved

	genKey := func(prefix string) string {
		k, sfx := prefix, 1
		for seen.Contains(k) {
			k = prefix + "_" + strconv.Itoa(sfx)
			sfx += 1
		}
		seen.Add(k)
		return k
	}

	for _, field := range fields {
		if field.Key != "" {
			continue
		}

		prefix := "key"
		if identityRegexp.MatchString(field.Name) {
			prefix = field.Name
		}
		field.Key = genKey(prefix)
	}
}

func validateSchema(dataset *entity.Dataset, fields []*entity.FieldSchema) error {
	if dataset.Spec == nil || dataset.Features == nil {
		return errno.InternalErrorf("dataset spec or features is nil, dataset_id=%d", dataset.ID)
	}

	availFields := gslice.Filter(fields, func(f *entity.FieldSchema) bool { return f.Available() })
	if len(availFields) == 0 {
		return errors.New("no available fields")
	}

	if empty := gslice.Filter(fields, func(f *entity.FieldSchema) bool { return f.Name == "" || f.Key == "" }); len(empty) > 0 {
		return errors.Errorf("field name or key is empty")
	}
	if badKeys := gslice.Filter(fields, func(f *entity.FieldSchema) bool { return !identityRegexp.MatchString(f.Key) }); len(badKeys) > 0 {
		return errors.Errorf("field key '%v' is not valid", badKeys)
	}

	names := gslice.Map(availFields, func(f *entity.FieldSchema) string { return f.Name })
	if dup := gslice.Dup(names); len(dup) > 0 {
		return errors.Errorf("field name '%v' duplicated", dup)
	}
	keys := gslice.Map(fields, func(f *entity.FieldSchema) string { return f.Key })
	if dup := gslice.Dup(keys); len(dup) > 0 {
		return errors.Errorf("field key '%v' duplicated", dup)
	}

	if dataset.Spec.MaxFieldCount > 0 && len(availFields) > int(dataset.Spec.MaxFieldCount) {
		return errors.Errorf("field_count %d exceed column_limit %d", len(fields), dataset.Spec.MaxFieldCount)
	}

	if !dataset.Features.MultiModal {
		s := gslice.Filter(availFields, func(f *entity.FieldSchema) bool { return f.ContentType.IsMultiModal() })
		if len(s) > 0 {
			return errors.Errorf("multi_modal is not enabled, fields=%v", gslice.Map(s, func(f *entity.FieldSchema) string { return f.Name }))
		}
	}
	return nil
}

func schemaCompatible(preFields, curFields []*entity.FieldSchema) bool {
	curFields = gslice.Filter(curFields, func(f *entity.FieldSchema) bool { return f.Available() })
	preMap := gslice.ToMap(preFields, func(f *entity.FieldSchema) (string, *entity.FieldSchema) { return f.Key, f })
	for _, cur := range curFields {
		pre, ok := preMap[cur.Key]
		if ok && !pre.CompatibleWith(cur) {
			return false
		}
	}

	return true
}

// mergeFields 合并上一版本和当前版本，返回所有字段. 注意: curFields 中 key 为空的字段视作新增.
func mergeSchema(ds *entity.Dataset, preSchema *entity.DatasetSchema, fields []*entity.FieldSchema) ([]*entity.FieldSchema, error) {
	fields = gslice.Filter(fields, func(f *entity.FieldSchema) bool { return f.Available() })
	preMap := gslice.ToMap(preSchema.Fields, func(f *entity.FieldSchema) (string, *entity.FieldSchema) { return f.Key, f })

	for _, cur := range fields {
		cur.Status = entity.FieldStatusAvailable
		delete(preMap, cur.Key)
	}

	for _, pre := range preMap {
		field := &entity.FieldSchema{}
		if err := copier.Copy(field, pre); err != nil { // make a copy to avoid mutate preFields
			return nil, errors.WithMessagef(err, "copy a field, key=%s", pre.Key)
		}
		field.Status = entity.FieldStatusDeleted
		fields = append(fields, field)
	}

	genFieldKeys(fields)
	if err := validateSchema(ds, fields); err != nil {
		return nil, errno.InvalidParamErr(err)
	}
	return fields, nil
}

func newSchemaOfDataset(dataset *entity.Dataset, fields []*entity.FieldSchema) *entity.DatasetSchema {
	immutable := false
	if dataset.Features != nil {
		immutable = !dataset.Features.EditSchema
	}

	return &entity.DatasetSchema{
		ID:        dataset.SchemaID,
		AppID:     dataset.AppID,
		SpaceID:   dataset.SpaceID,
		DatasetID: dataset.ID,
		Fields:    fields,
		Immutable: immutable,
		CreatedBy: dataset.CreatedBy,
		UpdatedBy: dataset.CreatedBy,
	}
}
