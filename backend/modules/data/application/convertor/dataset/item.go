// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gslice"

	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/data/domain/dataset"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/pkg/time"
)

func ItemDO2DTO(i *entity.Item) *dataset.DatasetItem {
	return &dataset.DatasetItem{
		ID:           gptr.Of(i.ID),
		AppID:        gptr.Of(i.AppID),
		SpaceID:      gptr.Of(i.SpaceID),
		DatasetID:    gptr.Of(i.DatasetID),
		SchemaID:     gptr.Of(i.SchemaID),
		ItemID:       gptr.Of(i.ItemID),
		ItemKey:      gptr.Of(i.ItemKey),
		Data:         gslice.Map(i.Data, func(f *entity.FieldData) *dataset.FieldData { return FieldDataDO2DTO(f) }),
		RepeatedData: gslice.Map(i.RepeatedData, func(d *entity.ItemData) *dataset.ItemData { return ItemDataDO2DTO(d) }),
		CreatedBy:    gptr.Of(i.CreatedBy),
		CreatedAt:    gptr.Of(i.CreatedAt.UnixMilli()),
		UpdatedBy:    gptr.Of(i.UpdatedBy),
		UpdatedAt:    gptr.Of(i.UpdatedAt.UnixMilli()),
		DataOmitted:  gptr.Of(false), // Notice: hard code to false
	}
}

func ItemDataDO2DTO(i *entity.ItemData) *dataset.ItemData {
	return &dataset.ItemData{
		ID:   gptr.Of(i.ID),
		Data: gslice.Map(i.Data, func(f *entity.FieldData) *dataset.FieldData { return FieldDataDO2DTO(f) }),
	}
}

func FieldDataDO2DTO(f *entity.FieldData) *dataset.FieldData {
	return &dataset.FieldData{
		Key:         gptr.Of(f.Key),
		Name:        gptr.Of(f.Name),
		ContentType: gptr.Of(ContentTypeDO2DTO(f.ContentType)),
		Format:      gptr.Of(FieldDisplayFormatDO2DTO(f.Format)),
		Content:     gptr.Of(f.Content),
		Attachments: gslice.Map(f.Attachments, func(a *entity.ObjectStorage) *dataset.ObjectStorage { return ObjectStorageDO2DTO(a) }),
		Parts:       gslice.Map(f.Parts, func(p *entity.FieldData) *dataset.FieldData { return FieldDataDO2DTO(p) }),
	}
}

func ObjectStorageDO2DTO(o *entity.ObjectStorage) *dataset.ObjectStorage {
	return &dataset.ObjectStorage{
		Provider: gptr.Of(ProviderDO2DTO(o.Provider)),
		Name:     gptr.Of(o.Name),
		URI:      gptr.Of(o.URI),
		URL:      gptr.Of(o.URL),
		ThumbURL: gptr.Of(o.ThumbURL),
	}
}

func ItemDTO2DO(s *dataset.DatasetItem) *entity.Item {
	if s == nil {
		return nil
	}
	t := &entity.Item{
		ID:           s.GetID(),
		AppID:        s.GetAppID(),
		SpaceID:      s.GetSpaceID(),
		DatasetID:    s.GetDatasetID(),
		SchemaID:     s.GetSchemaID(),
		ItemID:       s.GetItemID(),
		ItemKey:      s.GetItemKey(),
		Data:         gslice.Map(s.Data, FieldDataDTO2DO),
		RepeatedData: gslice.Map(s.RepeatedData, ItemDataDTO2DO),
		CreatedBy:    s.GetCreatedBy(),
		CreatedAt:    time.UnixMilliToTime(s.GetCreatedAt()),
		UpdatedBy:    s.GetUpdatedBy(),
		UpdatedAt:    time.UnixMilliToTime(s.GetUpdatedAt()),
	}
	return t
}

func ItemDataDTO2DO(s *dataset.ItemData) *entity.ItemData {
	return &entity.ItemData{
		ID:   s.GetID(),
		Data: gslice.Map(s.Data, FieldDataDTO2DO),
	}
}

func FieldDataDTO2DO(s *dataset.FieldData) *entity.FieldData {
	return &entity.FieldData{
		Key:         s.GetKey(),
		Name:        s.GetName(),
		ContentType: ContentTypeDTO2DO(s.GetContentType()),
		Format:      FieldDisplayFormatDTO2DO(s.GetFormat()),
		Content:     s.GetContent(),
		Attachments: gslice.Map(s.GetAttachments(), ObjectStorageDTO2DO),
		Parts:       gslice.Map(s.Parts, FieldDataDTO2DO),
	}
}

func ObjectStorageDTO2DO(s *dataset.ObjectStorage) *entity.ObjectStorage {
	return &entity.ObjectStorage{
		Provider: StorageProviderDTO2DO(s.GetProvider()),
		Name:     s.GetName(),
		URI:      s.GetURI(),
		URL:      s.GetURL(),
		ThumbURL: s.GetThumbURL(),
	}
}

func ItemErrorGroupDO2DTO(s *entity.ItemErrorGroup) *dataset.ItemErrorGroup {
	return &dataset.ItemErrorGroup{
		Type:       gptr.Of(dataset.ItemErrorType(gptr.Indirect(s.Type))),
		Summary:    s.Summary,
		ErrorCount: s.ErrorCount,
		Details:    gslice.Map(s.Details, ItemErrorDetailDO2DTO),
	}
}

func ItemErrorDetailDO2DTO(s *entity.ItemErrorDetail) *dataset.ItemErrorDetail {
	return &dataset.ItemErrorDetail{
		Index:      s.Index,
		StartIndex: s.StartIndex,
		EndIndex:   s.EndIndex,
		Message:    s.Message,
	}
}

func ItemErrorGroupDTO2DO(s *dataset.ItemErrorGroup) *entity.ItemErrorGroup {
	return &entity.ItemErrorGroup{
		Type:       gptr.Of(entity.ItemErrorType(s.GetType())),
		Summary:    s.Summary,
		ErrorCount: s.ErrorCount,
		Details:    gslice.Map(s.GetDetails(), ItemErrorDetailDTO2DO),
	}
}

func ItemErrorDetailDTO2DO(s *dataset.ItemErrorDetail) *entity.ItemErrorDetail {
	return &entity.ItemErrorDetail{
		Index:      s.Index,
		StartIndex: s.StartIndex,
		EndIndex:   s.EndIndex,
		Message:    s.Message,
	}
}
